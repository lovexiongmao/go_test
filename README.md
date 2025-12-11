# Go Web 后端服务

基于 Gin + GORM + Dig + Audit + Logrus 的 Web 后端服务框架。

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **依赖注入**: Dig (Uber)
- **认证**: JWT (JSON Web Token)
- **权限管理**: RBAC (基于角色的访问控制)
- **审计**: HTTP 审计中间件 + GORM 数据库审计插件
- **日志**: Logrus
- **数据库**: MySQL
- **API文档**: Swagger/OpenAPI
- **密码加密**: bcrypt

## 项目结构

```
.
├── cmd/
│   └── server/          # 服务入口
│       └── main.go
├── docs/
│   └── swagger/         # Swagger API 文档
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接和审计插件
│   ├── handler/         # HTTP处理器（用户、角色、权限、认证）
│   ├── logger/          # 日志模块
│   ├── middleware/      # 中间件（日志、审计、认证、权限校验）
│   ├── model/           # 数据模型（用户、角色、权限）
│   ├── repository/      # 数据访问层
│   ├── router/          # 路由配置
│   ├── service/         # 业务逻辑层
│   └── util/            # 工具函数（统一响应格式、JWT等）
├── sql/
│   └── permission_related_init_data.sql  # 权限系统初始数据
├── pkg/
│   └── dig/             # 依赖注入容器
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置数据库连接信息等。

### 3. 创建数据库

在MySQL中创建数据库：

```sql
CREATE DATABASE testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 初始化权限数据（可选）

如果需要使用预定义的权限、角色和用户数据，可以执行初始化 SQL：

```bash
mysql -u root -p testdb < sql/permission_related_init_data.sql
```

该 SQL 文件包含：
- 47 个预定义权限（用户、角色、权限、服务器、部署、监控等）
- 7 个预定义角色（超级管理员、管理员、运维工程师、开发人员等）
- 9 个示例用户（密码均为 `123456`，请在生产环境中修改）

### 5. 运行服务

```bash
go run cmd/server/main.go
```

或使用 Makefile：

```bash
make run
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### Swagger API 文档

项目集成了 Swagger UI，启动服务后访问：

```
http://localhost:8080/swagger/index.html
```

Swagger 文档提供了完整的 API 接口说明和在线测试功能。

### 认证接口

#### 用户登录

```bash
POST /api/v1/login
```

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "123456"
  }'
```

**响应示例**：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "超级管理员",
      "email": "admin@example.com"
    }
  }
}
```

### 用户管理

所有用户管理接口都需要 JWT 认证和相应权限。

- `POST /api/v1/users` - 创建用户（需要 `user:create` 权限）
- `GET /api/v1/users` - 获取用户列表（需要 `user:read` 权限）
- `GET /api/v1/users/:id` - 获取用户详情（需要 `user:read` 权限）
- `PUT /api/v1/users/:id` - 更新用户（需要 `user:update` 权限）
- `DELETE /api/v1/users/:id` - 删除用户（需要 `user:delete` 权限）

**请求示例**（需要先登录获取 token）：
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "password": "123456"
  }'
```

### 角色管理

所有角色管理接口都需要 JWT 认证和相应权限。

- `POST /api/v1/roles` - 创建角色（需要 `role:create` 权限）
- `GET /api/v1/roles` - 获取角色列表（需要 `role:read` 权限）
- `GET /api/v1/roles/:id` - 获取角色详情（需要 `role:read` 权限）
- `PUT /api/v1/roles/:id` - 更新角色（需要 `role:update` 权限）
- `DELETE /api/v1/roles/:id` - 删除角色（需要 `role:delete` 权限）
- `POST /api/v1/roles/:id/permissions` - 为角色分配权限（需要 `role:update` 权限）
- `DELETE /api/v1/roles/:id/permissions` - 移除角色权限（需要 `role:update` 权限）
- `GET /api/v1/roles/:id/permissions` - 获取角色权限列表（需要 `role:read` 权限）
- `POST /api/v1/roles/:id/users` - 为角色分配用户（需要 `role:update` 权限）
- `DELETE /api/v1/roles/:id/users` - 移除角色用户（需要 `role:update` 权限）
- `GET /api/v1/roles/:id/users` - 获取角色用户列表（需要 `role:read` 权限）

### 权限管理

所有权限管理接口都需要 JWT 认证和相应权限。

- `POST /api/v1/permissions` - 创建权限（需要 `permission:create` 权限）
- `GET /api/v1/permissions` - 获取权限列表（需要 `permission:read` 权限）
- `GET /api/v1/permissions/:id` - 获取权限详情（需要 `permission:read` 权限）
- `PUT /api/v1/permissions/:id` - 更新权限（需要 `permission:update` 权限）
- `DELETE /api/v1/permissions/:id` - 删除权限（需要 `permission:delete` 权限）

### 健康检查

```
GET /health
```

## 功能特性

### 核心功能

- ✅ RESTful API 设计
- ✅ 依赖注入（Dig）
- ✅ **JWT 认证** - 基于 JWT 的用户认证机制
- ✅ **RBAC 权限管理** - 基于角色的访问控制（用户-角色-权限）
- ✅ **权限校验中间件** - 自动校验用户是否有权限访问资源
- ✅ **数据库连接池管理** - 自动配置连接池参数，优化数据库性能
- ✅ 请求日志记录（Logrus）
- ✅ 优雅关闭服务
- ✅ 自动数据库迁移
- ✅ **环境变量配置** - 支持通过 `.env` 文件配置所有参数
- ✅ **密码加密** - 使用 bcrypt 加密用户密码

### 审计功能

项目实现了**双层审计机制**：

1. **HTTP 层面审计**（中间件）
   - 记录所有 API 请求信息
   - 包括：请求方法、路径、客户端 IP、User-Agent、用户 ID
   - 记录到日志文件

2. **数据库层面审计**（GORM 插件）
   - 自动记录所有表的增删改操作
   - 记录操作前后的数据变化（JSON 格式）
   - 记录操作者信息（用户 ID、IP）
   - 存储在 `audit_logs` 表中
   - **适用于所有模型和表**，无需额外配置

### 认证与授权

#### JWT 认证

- ✅ 用户登录后返回 JWT Token
- ✅ Token 包含用户 ID 和邮箱信息
- ✅ Token 过期时间可配置（默认 24 小时）
- ✅ 所有需要认证的接口都需要在 Header 中携带 Token

**Token 使用方式**：
```
Authorization: Bearer <your-token>
```

#### RBAC 权限管理

项目实现了完整的基于角色的访问控制（RBAC）系统：

- **用户（User）** - 系统使用者
- **角色（Role）** - 权限的集合，如"管理员"、"开发人员"等
- **权限（Permission）** - 具体的操作权限，如"user:create"、"role:read"等

**权限格式**：`资源:操作`
- 资源：如 `user`、`role`、`permission`、`server` 等
- 操作：如 `create`、`read`、`update`、`delete`、`execute` 等

**权限校验流程**：
1. 用户登录后获取 JWT Token
2. 请求需要认证的接口时，携带 Token
3. JWT 中间件验证 Token 并提取用户信息
4. 权限校验中间件检查用户是否有对应资源的操作权限
5. 权限检查通过后，请求继续处理

**权限检查逻辑**：
- 系统会查找用户关联的所有角色
- 检查这些角色是否拥有请求的资源操作权限
- 只要有一个角色拥有权限，即允许访问

### API 文档

- ✅ Swagger/OpenAPI 集成
- ✅ 自动生成 API 文档
- ✅ 在线 API 测试界面
- ✅ **Swagger UI 支持 Token 认证** - 在 Swagger UI 中可以设置 Token 进行接口测试

### 统一响应格式

- ✅ 统一的 HTTP 响应结构
- ✅ 标准化的错误处理
- ✅ 支持自定义响应消息

## 开发说明

### 添加新的API

1. 在 `internal/model/` 中定义数据模型
2. 在 `internal/repository/` 中实现数据访问层
3. 在 `internal/service/` 中实现业务逻辑
4. 在 `internal/handler/` 中实现HTTP处理器（添加 Swagger 注释）
5. 在 `internal/router/router.go` 中注册路由
6. 在 `pkg/dig/container.go` 中注册依赖
7. 运行 `swag init` 重新生成 Swagger 文档

### Swagger 文档

项目使用 [swaggo/swag](https://github.com/swaggo/swag) 生成 API 文档。

**生成文档**：
```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go
```

**访问文档**：
启动服务后访问 `http://localhost:8080/swagger/index.html`

**添加 API 注释示例**：
```go
// CreateUser 创建用户
// @Summary      创建用户
// @Description  创建一个新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "用户信息"
// @Success      201   {object}  util.Response{data=UserResponse}
// @Failure      400   {object}  util.Response
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ...
}
```

### 数据库连接池管理

项目在 `internal/database/database.go` 中实现了数据库连接池的配置：

```go:30:38:internal/database/database.go
// 配置连接池
sqlDB, err := db.DB()
if err != nil {
    return nil, err
}

sqlDB.SetMaxIdleConns(10)        // 最大空闲连接数：10
sqlDB.SetMaxOpenConns(100)       // 最大打开连接数：100
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间：1小时
```

**配置说明**：
- `SetMaxIdleConns(10)`: 设置连接池中空闲连接的最大数量，保持一定数量的连接可以快速复用
- `SetMaxOpenConns(100)`: 设置数据库打开连接的最大数量，防止连接数过多导致数据库压力过大
- `SetConnMaxLifetime(time.Hour)`: 设置连接的最大生存时间，超过时间的连接会被关闭并重新创建，避免长时间连接导致的网络问题

这些配置可以有效优化数据库连接的性能和稳定性。

### 环境变量配置

项目使用 `godotenv` 库支持通过 `.env` 文件配置所有参数，配置加载逻辑在 `internal/config/config.go` 中：

```go:39:64:internal/config/config.go
func LoadConfig() (*Config, error) {
    // 加载.env文件（如果存在）
    _ = godotenv.Load()

    config := &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
            Host: getEnv("SERVER_HOST", "0.0.0.0"),
            Mode: getEnv("GIN_MODE", "debug"),
        },
        Database: DatabaseConfig{
            Host:        getEnv("DB_HOST", "localhost"),
            Port:        getEnv("DB_PORT", "3306"),
            User:        getEnv("DB_USER", "root"),
            Password:    getEnv("DB_PASSWORD", ""),
            DBName:      getEnv("DB_NAME", "testdb"),
            AutoMigrate: getEnv("DB_AUTO_MIGRATE", "true") == "true",
        },
        Log: LogConfig{
            Level:     getEnv("LOG_LEVEL", "info"),
            Format:    getEnv("LOG_FORMAT", "text"),
            Output:    getEnv("LOG_OUTPUT", "stdout"),
            LogFile:   getEnv("APP_LOG_FILE", "logs/app.log"),
            AuditFile: getEnv("AUDIT_LOG_FILE", "logs/audit.log"),
        },
    }
    // ...
}
```

**支持的环境变量**：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `SERVER_PORT` | 服务端口 | `8080` |
| `SERVER_HOST` | 服务地址 | `0.0.0.0` |
| `GIN_MODE` | Gin 模式（debug/release/test） | `debug` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASSWORD` | 数据库密码 | 空 |
| `DB_NAME` | 数据库名称 | `testdb` |
| `DB_AUTO_MIGRATE` | 是否自动迁移 | `true` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `LOG_FORMAT` | 日志格式（json/text） | `text` |
| `LOG_OUTPUT` | 日志输出（stdout/file/both） | `stdout` |
| `APP_LOG_FILE` | 应用日志文件路径 | `logs/app.log` |
| `AUDIT_LOG_FILE` | 审计日志文件路径 | `logs/audit.log` |
| `JWT_SECRET` | JWT 密钥（生产环境必须修改） | `your-secret-key-change-in-production` |
| `JWT_EXPIRE_TIME` | Token 过期时间（分钟） | `1440`（24小时） |

**使用方式**：
1. 创建 `.env` 文件（项目根目录）
2. 设置需要的环境变量，例如：
   ```bash
   SERVER_PORT=8080
   DB_HOST=localhost
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=testdb
   ```
3. 项目启动时会自动加载 `.env` 文件中的配置
4. 如果环境变量未设置，会使用代码中的默认值

**注意**：
- `.env` 文件通常不应提交到版本控制系统，建议添加到 `.gitignore` 中
- 生产环境必须修改 `JWT_SECRET`，使用强随机字符串

### 日志级别

- `debug`: 详细调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 审计日志

#### HTTP 层面审计（中间件）

所有 API 请求都会自动记录审计日志到日志文件，包括：
- 请求方法
- 请求路径
- 客户端 IP
- User-Agent
- 用户 ID（如果已认证）

#### 数据库层面审计（GORM 插件）

所有数据库的增删改操作都会自动记录到 `audit_logs` 表，包括：
- 表名（`table_name`）
- 记录 ID（`record_id`）
- 操作类型（`action`: create/update/delete）
- 操作前的数据（`old_values`，JSON 格式）
- 操作后的数据（`new_values`，JSON 格式）
- 操作者用户 ID（`user_id`）
- 操作者 IP（`ip`）
- 操作时间（`created_at`）

**特点**：
- 自动适用于所有模型和表，无需额外配置
- 使用独立的数据库连接，不影响原事务
- 自动跳过 `audit_logs` 表自身的操作，避免递归
- 通过 Context 传递用户信息，支持在 Handler/Service 层设置

**使用示例**：

```go
// 在 Handler 中，用户信息已通过中间件自动设置到 Context
// 数据库操作会自动记录审计日志

// 创建用户
db.Create(&user)  // ✅ 自动记录创建审计日志

// 更新用户
db.Save(&user)    // ✅ 自动记录更新审计日志（包含旧值和新值）

// 删除用户
db.Delete(&user)  // ✅ 自动记录删除审计日志（包含旧值）
```

### 统一响应格式

项目提供了统一的 HTTP 响应格式（`internal/util/response.go`）：

```go
// 成功响应
util.Success(c, data)
util.SuccessWithMessage(c, "自定义消息", data)

// 创建成功响应
util.Created(c, data)
util.CreatedWithMessage(c, "创建成功", data)

// 错误响应
util.BadRequest(c, "错误消息")
util.NotFound(c, "资源不存在")
util.InternalServerError(c, "服务器错误")
```

响应格式：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": { ... }
}
```

### 权限系统使用说明

#### 权限模型

项目使用标准的 RBAC 模型：

```
用户 (User) ←→ 用户角色关联 (UserRole) ←→ 角色 (Role) ←→ 角色权限关联 (RolePermission) ←→ 权限 (Permission)
```

#### 权限初始化

项目提供了 `sql/permission_related_init_data.sql` 文件，包含：

1. **47 个预定义权限**，涵盖：
   - 用户管理（4个）：create、read、update、delete
   - 角色管理（4个）：create、read、update、delete
   - 权限管理（4个）：create、read、update、delete
   - 服务器管理（5个）：create、read、update、delete、execute
   - 应用部署（5个）：create、read、update、delete、execute
   - 监控管理（2个）：read、alert
   - 日志管理（2个）：read、download
   - 配置管理（4个）：create、read、update、delete
   - 审计日志（1个）：read
   - 数据库管理（6个）：create、read、update、delete、backup、restore
   - 容器管理（7个）：create、read、update、delete、start、stop、restart

2. **7 个预定义角色**：
   - `super_admin` - 超级管理员（所有权限）
   - `admin` - 管理员（大部分权限，排除危险操作）
   - `ops_engineer` - 运维工程师（服务器、部署、监控等）
   - `developer` - 开发人员（部署、查看等）
   - `viewer` - 只读用户（仅查看权限）
   - `dba` - 数据库管理员（数据库相关权限）
   - `sre` - SRE工程师（监控、日志、可靠性相关）

3. **9 个示例用户**（密码均为 `123456`，bcrypt 加密）

#### 添加新权限

1. 在数据库中创建权限记录：
   ```sql
   INSERT INTO permissions (name, display_name, description, resource, action, status, created_at, updated_at) 
   VALUES ('resource:action', '显示名称', '权限描述', 'resource', 'action', 1, NOW(), NOW());
   ```

2. 将权限分配给角色：
   ```sql
   INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
   VALUES (角色ID, 权限ID, NOW(), NOW());
   ```

3. 在路由中使用权限校验中间件：
   ```go
   router.POST("/resource", middleware.RequirePermission(userService, "resource", "action"), handler.CreateResource)
   ```

#### 权限校验示例

```go
// 在 router.go 中
users := auth.Group("/users")
{
    // 创建用户需要 user:create 权限
    users.POST("", middleware.RequirePermission(userService, "user", "create"), userHandler.CreateUser)
    
    // 查看用户需要 user:read 权限
    users.GET("", middleware.RequirePermission(userService, "user", "read"), userHandler.ListUsers)
}
```

## Makefile 命令

项目提供了 Makefile 方便开发：

```bash
make run        # 运行服务
make build      # 构建二进制文件
make test       # 运行测试
make clean      # 清理构建文件
make deps       # 下载依赖
make fmt        # 格式化代码
make vet        # 代码检查
make check      # 运行所有检查（fmt + vet + test）
make swagger    # 生成 Swagger 文档
```

## 许可证

MIT

