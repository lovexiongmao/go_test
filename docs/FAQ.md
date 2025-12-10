# 常见问题解答（FAQ）

本文档整理了项目开发中的常见问题和快速解答。

---

## 数据库相关

### Q: 如何在生产环境禁用自动数据库迁移？

**A:** 在 `.env` 文件中设置：
```bash
DB_AUTO_MIGRATE=false
```

详细说明见：`docs/DATABASE_MIGRATION_BEST_PRACTICES.md`

---

### Q: 审计表没有被创建，为什么？

**A:** 检查以下几点：
1. `DB_AUTO_MIGRATE` 是否设置为 `true`
2. `pkg/dig/container.go` 中 `AuditLog` 是否取消注释
3. 查看启动日志是否有错误信息

---

## API 设计相关

### Q: 如何添加新的 RESTful API？

**A:** 遵循以下步骤：
1. 创建 Model → `internal/model/`
2. 创建 Repository → `internal/repository/`
3. 创建 Service → `internal/service/`
4. 创建 Handler → `internal/handler/`
5. 注册路由 → `internal/router/router.go`

详细指南见：`docs/RESTFUL_GUIDE.md`

---

### Q: 什么时候需要创建新的 API 版本（v2）？

**A:** 当有以下破坏性变更时：
- 删除字段
- 修改字段类型
- 修改字段名
- 修改响应结构
- 修改必需字段

详细说明见：`docs/API_VERSIONING.md`

---

### Q: 如何保持 RESTful 风格？

**A:** 遵循以下原则：
- 使用 HTTP 方法表示操作（GET、POST、PUT、DELETE）
- 使用资源路径（复数名词、版本控制）
- 使用正确的 HTTP 状态码
- 统一响应格式

详细指南见：`docs/RESTFUL_GUIDE.md`

---

## 代码实现相关

### Q: omitempty 标签的作用是什么？

**A:** 
- `omitempty` 是 JSON 序列化标签，不是 GORM 标签
- 当字段为零值时，序列化为 JSON 时忽略该字段
- 例如：`json:"deleted_at,omitempty"` 表示如果 `deleted_at` 为 `nil`，则不包含在 JSON 中

---

### Q: 如何让密码字段不返回给前端？

**A:** 在模型中使用 `json:"-"`：
```go
Password string `gorm:"type:varchar(255);not null" json:"-"`
```
前端不会收到该字段，而不是用 `*` 代替。

---

### Q: 普通索引和唯一索引有什么区别？

**A:**
- **普通索引（index）**：允许重复值，主要用于提高查询速度
- **唯一索引（uniqueIndex）**：不允许重复值，既提高查询速度又保证唯一性

---

## 项目配置相关

### Q: 如何查看所有环境变量配置？

**A:** 查看 `internal/config/config.go` 文件，所有配置项都有说明。

---

### Q: 如何修改日志级别？

**A:** 在 `.env` 文件中设置：
```bash
LOG_LEVEL=debug  # debug, info, warn, error
```

---

## 开发流程相关

### Q: 如何添加新的功能模块？

**A:** 参考用户模块的实现：
1. Model → Repository → Service → Handler → Router
2. 在 `pkg/dig/container.go` 中注册依赖
3. 遵循 RESTful 设计规范

---

### Q: 如何查找之前的问题记录？

**A:** 
1. 查看 `docs/QUESTIONS.md` - 详细的问题记录
2. 查看 `docs/FAQ.md` - 常见问题快速解答
3. 使用文档中的关键词索引

---

## 部署相关

### Q: 生产环境需要注意什么？

**A:**
1. **数据库迁移**：设置 `DB_AUTO_MIGRATE=false`，使用专门的迁移工具
2. **日志级别**：设置为 `info` 或 `warn`，不要使用 `debug`
3. **Gin 模式**：设置为 `release`
4. **数据库连接**：配置正确的生产数据库连接信息

---

## 更多资源

- `docs/RESTFUL_GUIDE.md` - RESTful API 设计指南
- `docs/API_VERSIONING.md` - API 版本控制详解
- `docs/DATABASE_MIGRATION_BEST_PRACTICES.md` - 数据库迁移最佳实践
- `docs/QUESTIONS.md` - 详细的问题记录

---

## 贡献

如果发现新的常见问题，请添加到本文档中。
