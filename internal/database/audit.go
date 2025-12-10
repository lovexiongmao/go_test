package database

import (
	"time"

	"gorm.io/gorm"
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

	// Update回调
	callback.Update().After("gorm:update").Register("audit:update", p.auditUpdate)

	// Delete回调
	callback.Delete().After("gorm:delete").Register("audit:delete", p.auditDelete)

	return nil
}

func (p *AuditPlugin) auditCreate(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 这里可以记录创建操作的审计日志
	// 实际实现可以根据需求扩展
}

func (p *AuditPlugin) auditUpdate(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 这里可以记录更新操作的审计日志
	// 实际实现可以根据需求扩展
}

func (p *AuditPlugin) auditDelete(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 这里可以记录删除操作的审计日志
	// 实际实现可以根据需求扩展
}
