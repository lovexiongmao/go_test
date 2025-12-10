# 问题记录文档

本文档用于记录开发过程中遇到的问题和解决方案，方便后续查找和参考。

## 使用说明

- 按时间倒序记录（最新的问题在顶部）
- 使用关键词标签方便搜索
- 每个问题包含：问题描述、解决方案、相关文件、关键词

---

## 问题记录

### 2024-01-XX - 数据库迁移最佳实践

**问题：** 在 main.go 中每次服务启动都会执行 AutoMigrate，这在生产环境是否合理？

**解决方案：**
- 添加了环境变量 `DB_AUTO_MIGRATE` 控制是否执行自动迁移
- 开发环境：`DB_AUTO_MIGRATE=true`（默认）
- 生产环境：`DB_AUTO_MIGRATE=false`，使用专门的迁移工具（如 golang-migrate）

**相关文件：**
- `internal/config/config.go` - 添加 AutoMigrate 配置
- `pkg/dig/container.go` - 更新 InitializeDatabase 函数
- `cmd/server/main.go` - 添加条件判断
- `docs/DATABASE_MIGRATION_BEST_PRACTICES.md` - 详细文档

**关键词：** `数据库迁移` `AutoMigrate` `生产环境` `golang-migrate` `最佳实践`

---

### 2024-01-XX - 审计功能实现

**问题：** audit.go 文件在执行新增操作时好像没有执行，没有在数据库里创建审计表，这个文件的主要作用是什么？

**解决方案：**
- audit.go 中的回调函数原本是空的，需要实现
- 已实现完整的审计功能（auditCreate、auditUpdate、auditDelete）
- 审计表会在启动时自动创建（如果 DB_AUTO_MIGRATE=true）

**相关文件：**
- `internal/database/audit.go` - 审计插件实现
- `internal/middleware/audit.go` - HTTP 层面审计
- `pkg/dig/container.go` - 审计表自动迁移

**关键词：** `审计` `audit` `数据库审计` `GORM插件` `回调函数`

---

### 2024-01-XX - API 版本控制

**问题：** 在 RESTful 中，资源路径为什么要使用版本控制，这样的好处是什么？

**解决方案：**
- 版本控制允许向后兼容，保护现有客户端
- 支持平滑升级，给客户端迁移时间
- 降低风险，避免大规模故障
- 支持业务需求变化

**相关文件：**
- `docs/API_VERSIONING.md` - 详细文档
- `internal/router/router.go` - 路由版本控制示例

**关键词：** `版本控制` `API版本` `向后兼容` `RESTful` `v1` `v2`

---

### 2024-01-XX - 版本兼容性问题

**问题：** 如果是删除字段或修改字段类型这种破坏性的改变，即使使用了 v2 版本，如果底层还是查询的同一张表，对于 v1 的 API 也是错误的，这个要怎么解决？

**解决方案：**
- 使用 DTO（Data Transfer Object）进行版本转换
- 数据库表结构保持向后兼容（不删除字段，只添加新字段）
- 在应用层使用不同的 DTO 来适配不同版本的 API
- 从同一个 Model 转换到不同版本的响应结构

**相关文件：**
- `docs/VERSION_COMPATIBILITY.md` - 详细文档（如果存在）
- `docs/API_VERSIONING.md` - 版本控制文档

**关键词：** `版本兼容` `DTO` `字段删除` `字段类型变更` `向后兼容`

---

### 2024-01-XX - RESTful 风格

**问题：** 项目中的 RESTful 风格是如何体现的，以及以后如果不断完善功能，要怎么做才能保持 RESTful 风格？

**解决方案：**
- 使用 HTTP 方法正确表示操作（GET、POST、PUT、DELETE）
- 使用资源路径设计（复数名词、版本控制）
- 使用正确的 HTTP 状态码
- 统一响应格式

**相关文件：**
- `docs/RESTFUL_GUIDE.md` - 完整的 RESTful 设计指南
- `internal/router/router.go` - 路由实现
- `internal/handler/user.go` - Handler 实现

**关键词：** `RESTful` `HTTP方法` `资源路径` `状态码` `API设计`

---

### 2024-01-XX - omitempty 标签作用

**问题：** omitempty 在 gorm 中的标签是什么作用？对于密码不返回给前端，那前端得到的是什么？普通索引和唯一索引有什么区别？

**解决方案：**
- `omitempty` 是 JSON 序列化标签，不是 GORM 标签
- 当字段为零值时，序列化为 JSON 时忽略该字段
- 密码字段使用 `json:"-"` 完全忽略，前端不会收到该字段
- 普通索引允许重复值，唯一索引不允许重复值

**相关文件：**
- `internal/model/user.go` - 模型定义示例

**关键词：** `omitempty` `json标签` `索引` `唯一索引` `密码字段`

---

## 关键词索引

### 数据库相关
- `数据库迁移` - 见 "数据库迁移最佳实践"
- `AutoMigrate` - 见 "数据库迁移最佳实践"
- `审计` - 见 "审计功能实现"
- `索引` - 见 "omitempty 标签作用"

### API 设计相关
- `RESTful` - 见 "RESTful 风格"
- `版本控制` - 见 "API 版本控制"
- `版本兼容` - 见 "版本兼容性问题"
- `HTTP方法` - 见 "RESTful 风格"

### 代码实现相关
- `DTO` - 见 "版本兼容性问题"
- `GORM插件` - 见 "审计功能实现"
- `回调函数` - 见 "审计功能实现"

### 最佳实践
- `生产环境` - 见 "数据库迁移最佳实践"
- `最佳实践` - 见多个问题

---

## 添加新问题的模板

```markdown
### YYYY-MM-DD - 问题标题

**问题：** 问题描述

**解决方案：**
- 解决方案要点1
- 解决方案要点2

**相关文件：**
- `文件路径1` - 说明
- `文件路径2` - 说明

**关键词：** `关键词1` `关键词2` `关键词3`
```

---

## 搜索技巧

1. **使用 Cmd+F（Mac）或 Ctrl+F（Windows）** 在当前文档中搜索关键词
2. **查看关键词索引** 快速定位相关问题
3. **按时间查找** 如果记得大概时间，可以按日期查找
4. **按文件查找** 如果记得相关文件，可以在相关文件部分查找

---

## 更新日志

- 2024-01-XX - 创建文档
- 2024-01-XX - 添加数据库迁移问题
- 2024-01-XX - 添加审计功能问题
- 2024-01-XX - 添加 API 版本控制问题
