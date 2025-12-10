# 数据库迁移最佳实践

## 问题分析

### 当前实现的问题

在 `main.go` 中每次服务启动都执行 `AutoMigrate()`：

```go
// 初始化数据库表
if err := dig.InitializeDatabase(db); err != nil {
    return fmt.Errorf("初始化数据库失败: %v", err)
}
```

**生产环境的风险：**

1. **自动修改表结构** - AutoMigrate 会自动添加/删除字段，可能破坏现有数据
2. **无版本控制** - 无法追踪数据库变更历史
3. **无回滚机制** - 迁移失败无法回滚
4. **无审核流程** - 直接在生产环境修改表结构
5. **可能影响服务** - 迁移过程中可能锁表，影响正在运行的服务
6. **无法控制时机** - 每次启动都执行，无法选择迁移时机

---

## 最佳实践

### 1. **开发环境 vs 生产环境**

| 环境 | 策略 | 工具 |
|------|------|------|
| **开发环境** | 可以使用 AutoMigrate | GORM AutoMigrate |
| **测试环境** | 使用迁移工具 | golang-migrate / gormigrate |
| **生产环境** | **必须使用迁移工具** | golang-migrate / gormigrate |

### 2. **推荐的迁移工具**

#### 方案 A：golang-migrate（推荐）

- ✅ 版本化的 SQL 迁移文件
- ✅ 支持 up/down 迁移（可回滚）
- ✅ 独立的迁移命令
- ✅ 支持多种数据库
- ✅ 迁移历史追踪

#### 方案 B：gormigrate

- ✅ 基于 GORM
- ✅ Go 代码编写迁移
- ✅ 支持回滚
- ✅ 迁移历史追踪

---

## 实现方案

### 方案 1：环境变量控制（快速改进）

**优点：**
- 实现简单
- 不需要额外工具
- 可以快速应用到现有项目

**缺点：**
- 仍然使用 AutoMigrate（不够安全）
- 没有版本控制

### 方案 2：使用迁移工具（推荐）

**优点：**
- 版本控制
- 可回滚
- 生产环境安全
- 可审核

**缺点：**
- 需要额外工具
- 需要学习成本

---

## 实现代码

### 方案 1：环境变量控制

#### 1. 更新配置

```go
// internal/config/config.go
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    DSN      string
    AutoMigrate bool // 新增：是否自动迁移
}
```

#### 2. 更新初始化逻辑

```go
// pkg/dig/container.go
func InitializeDatabase(db *gorm.DB, cfg *config.Config) error {
    // 只在开发环境或明确配置时执行迁移
    if !cfg.Database.AutoMigrate {
        return nil
    }
    
    return db.AutoMigrate(
        &model.User{},
        &database.AuditLog{},
    )
}
```

#### 3. 更新 main.go

```go
// cmd/server/main.go
func startServer(
    cfg *config.Config,
    log *logger.Logger,
    db *gorm.DB,
    r *gin.Engine,
) error {
    // 只在配置允许时初始化数据库表
    if cfg.Database.AutoMigrate {
        if err := dig.InitializeDatabase(db, cfg); err != nil {
            return fmt.Errorf("初始化数据库失败: %v", err)
        }
        log.Info("数据库表初始化成功！")
    } else {
        log.Info("跳过数据库自动迁移（生产环境模式）")
    }
    
    // ... 其他代码
}
```

---

## 完整实现（推荐方案）

我已经为你实现了**方案 1：环境变量控制**，这是最快速且安全的改进方式。

### 已完成的修改

1. ✅ 在 `config.go` 中添加了 `AutoMigrate` 配置项
2. ✅ 在 `container.go` 中更新了 `InitializeDatabase` 函数，支持配置控制
3. ✅ 在 `main.go` 中添加了条件判断，只在配置允许时执行迁移

### 使用方法

#### 开发环境（启用自动迁移）

在 `.env` 文件中：
```bash
DB_AUTO_MIGRATE=true
```

或者不设置（默认为 `true`）

#### 生产环境（禁用自动迁移）

在 `.env` 文件中：
```bash
DB_AUTO_MIGRATE=false
```

启动服务时会跳过自动迁移，日志会显示：
```
跳过数据库自动迁移（生产环境模式，请使用专门的迁移工具）
```

---

## 生产环境迁移工具推荐

### 方案 A：golang-migrate（强烈推荐）

#### 安装

```bash
# macOS
brew install golang-migrate

# 或使用 Go
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### 创建迁移文件

```bash
# 创建迁移目录
mkdir -p migrations

# 创建迁移文件
migrate create -ext sql -dir migrations -seq create_users_table
```

这会创建两个文件：
- `migrations/000001_create_users_table.up.sql` - 升级迁移
- `migrations/000001_create_users_table.down.sql` - 回滚迁移

#### 编写迁移 SQL

**000001_create_users_table.up.sql:**
```sql
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    status INT DEFAULT 1,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**000001_create_users_table.down.sql:**
```sql
DROP TABLE IF EXISTS users;
```

#### 执行迁移

```bash
# 升级到最新版本
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/dbname" up

# 回滚一个版本
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/dbname" down 1

# 查看当前版本
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/dbname" version
```

### 方案 B：gormigrate（基于 GORM）

#### 安装

```bash
go get -u github.com/go-gormigrate/gormigrate/v2
```

#### 使用示例

```go
// migrations/migrations.go
package migrations

import (
    "go_test/internal/model"
    "go_test/internal/database"
    "gorm.io/gorm"
    "github.com/go-gormigrate/gormigrate/v2"
)

func GetMigrations() []*gormigrate.Migration {
    return []*gormigrate.Migration{
        {
            ID: "20240101000001",
            Migrate: func(tx *gorm.DB) error {
                return tx.AutoMigrate(&model.User{})
            },
            Rollback: func(tx *gorm.DB) error {
                return tx.Migrator().DropTable(&model.User{})
            },
        },
        {
            ID: "20240101000002",
            Migrate: func(tx *gorm.DB) error {
                return tx.AutoMigrate(&database.AuditLog{})
            },
            Rollback: func(tx *gorm.DB) error {
                return tx.Migrator().DropTable(&database.AuditLog{})
            },
        },
    }
}
```

#### 执行迁移

```go
// cmd/migrate/main.go
package main

import (
    "go_test/internal/config"
    "go_test/internal/database"
    "go_test/migrations"
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    cfg, _ := config.LoadConfig()
    db, _ := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
    
    m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())
    
    if err := m.Migrate(); err != nil {
        panic(err)
    }
}
```

---

## 总结

### ✅ 已实现的改进

1. **环境变量控制** - 通过 `DB_AUTO_MIGRATE` 控制是否执行自动迁移
2. **生产环境安全** - 生产环境默认禁用自动迁移
3. **清晰的日志** - 明确提示是否执行了迁移

### 📋 最佳实践清单

#### 开发环境
- ✅ 可以使用 `DB_AUTO_MIGRATE=true`
- ✅ 快速迭代，方便开发

#### 测试环境
- ⚠️ 建议使用迁移工具
- ⚠️ 模拟生产环境流程

#### 生产环境
- ❌ **必须禁用** `DB_AUTO_MIGRATE=false`
- ✅ 使用专门的迁移工具（golang-migrate 或 gormigrate）
- ✅ 迁移前备份数据库
- ✅ 在维护窗口执行迁移
- ✅ 测试回滚方案
- ✅ 监控迁移过程

### 🚀 下一步建议

1. **短期**：使用当前的环境变量控制方案
2. **中期**：引入 golang-migrate 工具
3. **长期**：建立完整的 CI/CD 迁移流程

---

## 常见问题

### Q: 为什么生产环境不能使用 AutoMigrate？

**A:** AutoMigrate 有以下风险：
- 自动修改表结构，可能破坏数据
- 无法回滚
- 无法控制迁移时机
- 无版本控制
- 可能锁表影响服务

### Q: 如何知道数据库是否需要迁移？

**A:** 使用迁移工具可以：
- 查看当前版本：`migrate version`
- 查看待迁移版本：`migrate -path migrations -database "..." up -dryrun`

### Q: 迁移失败怎么办？

**A:** 
1. 使用迁移工具的回滚功能
2. 从备份恢复数据库
3. 修复迁移脚本后重新执行

---

## 参考资源

- [golang-migrate 文档](https://github.com/golang-migrate/migrate)
- [gormigrate 文档](https://github.com/go-gormigrate/gormigrate)
- [GORM 迁移文档](https://gorm.io/docs/migration.html)
