package database

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// Context keys for audit information
const (
	AuditUserIDKey    = "audit:user_id"
	AuditIPKey        = "audit:ip"
	AuditOldValuesKey = "audit:old_values" // 用于存储更新前的旧值
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`

	ModelTableName string `gorm:"type:varchar(100);index;column:table_name"` // 表名，使用column标签避免与方法名冲突
	RecordID       uint   `gorm:"index"`
	Action         string `gorm:"type:varchar(20);index"` // create, update, delete
	OldValues      string `gorm:"type:text"`
	NewValues      string `gorm:"type:text"`
	UserID         uint   `gorm:"index"`
	IP             string `gorm:"type:varchar(50)"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditPlugin GORM审计插件
type AuditPlugin struct {
	db *gorm.DB
}

// NewAuditPlugin 创建审计插件
func NewAuditPlugin(db *gorm.DB) *AuditPlugin {
	return &AuditPlugin{db: db}
}

// Name 返回插件名称
func (p *AuditPlugin) Name() string {
	return "audit"
}

// Initialize 初始化插件
func (p *AuditPlugin) Initialize(db *gorm.DB) error {
	p.db = db

	// 注册回调
	callback := db.Callback()

	// Create回调
	callback.Create().After("gorm:create").Register("audit:create", p.auditCreate)

	// Update回调：在更新前获取旧值，在更新后记录审计日志
	callback.Update().Before("gorm:update").Register("audit:before_update", p.auditBeforeUpdate)
	callback.Update().After("gorm:update").Register("audit:update", p.auditUpdate)

	// Delete回调
	callback.Delete().After("gorm:delete").Register("audit:delete", p.auditDelete)

	return nil
}

// auditCreate 记录创建操作的审计日志
func (p *AuditPlugin) auditCreate(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 跳过审计日志表自身的操作，避免递归
	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" {
		return
	}

	// 获取表名
	tableName := p.getTableName(db)
	if tableName == "" {
		return
	}

	// 获取记录ID
	recordID := p.getRecordID(db)
	if recordID == 0 {
		return
	}

	// 获取新值
	newValues := p.serializeModel(db.Statement.Dest)

	// 获取用户ID和IP
	userID := p.getUserID(db)
	ip := p.getIP(db)

	// 创建审计日志
	auditLog := AuditLog{
		ModelTableName: tableName,
		RecordID:       recordID,
		Action:         "create",
		OldValues:      "", // 创建操作没有旧值
		NewValues:      newValues,
		UserID:         userID,
		IP:             ip,
	}

	// 使用新的数据库连接保存审计日志，避免影响原事务
	p.db.Session(&gorm.Session{NewDB: true}).Create(&auditLog)
}

// auditBeforeUpdate 在更新前获取旧值并存储到context中
func (p *AuditPlugin) auditBeforeUpdate(db *gorm.DB) {
	// 跳过审计日志表自身的操作
	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" {
		return
	}

	// 获取表名
	tableName := p.getTableName(db)
	if tableName == "" {
		return
	}

	// 获取记录ID
	recordID := p.getRecordID(db)
	if recordID == 0 {
		return
	}

	// 在更新前查询旧值
	if db.Statement.Schema != nil {
		oldModel := reflect.New(db.Statement.Schema.ModelType).Interface()
		// 使用新的数据库连接查询，避免影响原事务
		// 使用 NewDB: true 确保使用独立的连接，避免事务隔离问题
		// 但保留原 context，以便后续回调可以访问
		oldDB := db.Session(&gorm.Session{NewDB: true})
		if err := oldDB.First(oldModel, recordID).Error; err == nil {
			oldValues := p.serializeModel(oldModel)
			// 将旧值存储到context中
			if db.Statement.Context != nil {
				ctx := db.Statement.Context
				ctx = context.WithValue(ctx, AuditOldValuesKey, oldValues)
				db.Statement.Context = ctx
			} else {
				// 如果 context 为空，创建一个新的 context
				ctx := context.Background()
				ctx = context.WithValue(ctx, AuditOldValuesKey, oldValues)
				db.Statement.Context = ctx
			}
		}
	}
}

// auditUpdate 记录更新操作的审计日志
func (p *AuditPlugin) auditUpdate(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 跳过审计日志表自身的操作
	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" {
		return
	}

	// 获取表名
	tableName := p.getTableName(db)
	if tableName == "" {
		return
	}

	// 获取记录ID
	recordID := p.getRecordID(db)
	if recordID == 0 {
		return
	}

	// 获取新旧值
	// 旧值应该已经在 auditBeforeUpdate 中获取并存储到context中
	oldValues := ""
	if db.Statement.Context != nil {
		if oldVals, ok := db.Statement.Context.Value(AuditOldValuesKey).(string); ok {
			oldValues = oldVals
		}
	}

	// 如果context中没有旧值，尝试查询（兼容性处理）
	if oldValues == "" && db.Statement.Schema != nil {
		oldModel := reflect.New(db.Statement.Schema.ModelType).Interface()
		// 使用 Unscoped 查询，因为可能已经被软删除
		if err := p.db.Session(&gorm.Session{NewDB: true}).Unscoped().First(oldModel, recordID).Error; err == nil {
			oldValues = p.serializeModel(oldModel)
		}
	}

	newValues := p.serializeModel(db.Statement.Dest)

	// 获取用户ID和IP
	userID := p.getUserID(db)
	ip := p.getIP(db)

	// 创建审计日志
	auditLog := AuditLog{
		ModelTableName: tableName,
		RecordID:       recordID,
		Action:         "update",
		OldValues:      oldValues,
		NewValues:      newValues,
		UserID:         userID,
		IP:             ip,
	}

	// 使用新的数据库连接保存审计日志
	p.db.Session(&gorm.Session{NewDB: true}).Create(&auditLog)
}

// auditDelete 记录删除操作的审计日志
func (p *AuditPlugin) auditDelete(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 跳过审计日志表自身的操作
	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" {
		return
	}

	// 获取表名
	tableName := p.getTableName(db)
	if tableName == "" {
		return
	}

	// 获取记录ID
	recordID := p.getRecordID(db)
	if recordID == 0 {
		return
	}

	// 获取旧值（删除前的值）
	oldValues := ""
	if db.Statement.Schema != nil {
		oldModel := reflect.New(db.Statement.Schema.ModelType).Interface()
		// 使用 Unscoped 查询，因为可能已经被软删除
		if err := db.Session(&gorm.Session{}).Unscoped().First(oldModel, recordID).Error; err == nil {
			oldValues = p.serializeModel(oldModel)
		}
	}

	// 获取用户ID和IP
	userID := p.getUserID(db)
	ip := p.getIP(db)

	// 创建审计日志
	auditLog := AuditLog{
		ModelTableName: tableName,
		RecordID:       recordID,
		Action:         "delete",
		OldValues:      oldValues,
		NewValues:      "", // 删除操作没有新值
		UserID:         userID,
		IP:             ip,
	}

	// 使用新的数据库连接保存审计日志
	p.db.Session(&gorm.Session{NewDB: true}).Create(&auditLog)
}

// getTableName 获取表名
func (p *AuditPlugin) getTableName(db *gorm.DB) string {
	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}
	if db.Statement.Table != "" {
		return db.Statement.Table
	}
	return ""
}

// getRecordID 获取记录ID（主键值）
func (p *AuditPlugin) getRecordID(db *gorm.DB) uint {
	// 首先尝试从 Dest 反射获取ID
	if db.Statement.Dest != nil {
		destValue := reflect.ValueOf(db.Statement.Dest)
		if destValue.Kind() == reflect.Ptr {
			if destValue.IsNil() {
				return 0
			}
			destValue = destValue.Elem()
		}
		if destValue.Kind() == reflect.Struct {
			idField := destValue.FieldByName("ID")
			if idField.IsValid() {
				switch idField.Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					return uint(idField.Uint())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					id := idField.Int()
					if id > 0 {
						return uint(id)
					}
				}
			}
		}
	}

	// 尝试从 Schema 的主键字段获取
	if db.Statement.Schema != nil && db.Statement.Schema.PrioritizedPrimaryField != nil {
		if db.Statement.Dest != nil {
			destValue := reflect.ValueOf(db.Statement.Dest)
			if destValue.Kind() == reflect.Ptr {
				if destValue.IsNil() {
					return 0
				}
				destValue = destValue.Elem()
			}
			if destValue.Kind() == reflect.Struct {
				fieldValue := destValue.FieldByName(db.Statement.Schema.PrioritizedPrimaryField.Name)
				if fieldValue.IsValid() {
					switch fieldValue.Kind() {
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						return uint(fieldValue.Uint())
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						id := fieldValue.Int()
						if id > 0 {
							return uint(id)
						}
					}
				}
			}
		}
	}

	// 尝试从 Statement 的 Vars 中获取（用于 Delete 操作）
	if len(db.Statement.Vars) > 0 {
		for _, v := range db.Statement.Vars {
			if id, ok := v.(uint); ok && id > 0 {
				return id
			}
			if id, ok := v.(uint64); ok && id > 0 {
				return uint(id)
			}
			if id, ok := v.(int); ok && id > 0 {
				return uint(id)
			}
			if id, ok := v.(int64); ok && id > 0 {
				return uint(id)
			}
		}
	}

	return 0
}

// serializeModel 序列化模型为JSON字符串
func (p *AuditPlugin) serializeModel(model interface{}) string {
	if model == nil {
		return ""
	}

	// 使用反射获取实际值
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}

	// 如果是切片，处理每个元素
	if value.Kind() == reflect.Slice {
		var results []map[string]interface{}
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i).Interface()
			if data, err := p.modelToMap(elem); err == nil {
				results = append(results, data)
			}
		}
		if data, err := json.Marshal(results); err == nil {
			return string(data)
		}
		return ""
	}

	// 单个模型
	if data, err := p.modelToMap(model); err == nil {
		if jsonData, err := json.Marshal(data); err == nil {
			return string(jsonData)
		}
	}

	return ""
}

// modelToMap 将模型转换为map，过滤敏感字段
func (p *AuditPlugin) modelToMap(model interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return data, nil
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return data, fmt.Errorf("model is not a struct")
	}

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// 跳过未导出的字段
		if !fieldValue.CanInterface() {
			continue
		}

		// 跳过密码等敏感字段
		if field.Name == "Password" {
			continue
		}

		// 跳过 DeletedAt（软删除字段）
		if field.Name == "DeletedAt" {
			continue
		}

		// 获取字段值
		fieldInterface := fieldValue.Interface()

		// 处理指针类型
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				data[field.Name] = nil
				continue
			}
			fieldInterface = fieldValue.Elem().Interface()
		}

		data[field.Name] = fieldInterface
	}

	return data, nil
}

// getUserID 从context中获取用户ID
func (p *AuditPlugin) getUserID(db *gorm.DB) uint {
	if db.Statement.Context == nil {
		return 0
	}

	if userID, ok := db.Statement.Context.Value(AuditUserIDKey).(uint); ok {
		return userID
	}

	if userID, ok := db.Statement.Context.Value(AuditUserIDKey).(uint64); ok {
		return uint(userID)
	}

	if userID, ok := db.Statement.Context.Value(AuditUserIDKey).(int); ok && userID > 0 {
		return uint(userID)
	}

	return 0
}

// getIP 从context中获取IP地址
func (p *AuditPlugin) getIP(db *gorm.DB) string {
	if db.Statement.Context == nil {
		return ""
	}

	if ip, ok := db.Statement.Context.Value(AuditIPKey).(string); ok {
		return ip
	}

	return ""
}

// SetAuditContext 设置审计上下文信息（在handler或service中调用）
// 需要在调用数据库操作前设置context
// 使用示例：
//
//	import "context"
//	ctx := context.WithValue(c.Request.Context(), database.AuditUserIDKey, userID)
//	ctx = context.WithValue(ctx, database.AuditIPKey, clientIP)
//	db = db.WithContext(ctx)
func SetAuditContext(db *gorm.DB, userID uint, ip string) *gorm.DB {
	// 注意：这个函数需要配合 context 使用
	// 更好的方式是在 handler 中直接设置 context
	return db
}
