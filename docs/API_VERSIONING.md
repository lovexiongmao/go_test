# RESTful API 版本控制详解

## 一、为什么需要版本控制？

### 1. **业务需求变化**

在实际开发中，API 会不断演进：

**场景示例：**

**v1 版本（初始版本）：**
```json
// GET /api/v1/users/:id
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

**v2 版本（业务需求变化）：**
```json
// GET /api/v2/users/:id
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "avatar": "https://example.com/avatar.jpg",  // 新增字段
  "profile": {                                  // 新增嵌套对象
    "bio": "软件工程师",
    "location": "北京"
  }
}
```

**问题：** 如果直接修改 v1，会导致：
- 旧版客户端（移动 App、前端页面）可能因为缺少新字段而崩溃
- 无法向后兼容

### 2. **破坏性变更（Breaking Changes）**

以下变更会破坏现有客户端：

| 变更类型 | 示例 | 影响 |
|---------|------|------|
| **删除字段** | 移除 `email` 字段 | 客户端读取 `email` 会失败 |
| **修改字段类型** | `age` 从 `int` 改为 `string` | 类型不匹配导致解析错误 |
| **修改字段名** | `name` 改为 `username` | 客户端找不到字段 |
| **修改响应结构** | 扁平结构改为嵌套结构 | 客户端解析逻辑失效 |
| **修改必需字段** | 新增必填字段 | 创建请求会失败 |
| **修改 HTTP 方法** | POST 改为 PUT | 客户端调用方式错误 |

### 3. **多客户端共存**

一个 API 通常服务于多个客户端：

```
API 服务
├── iOS App (v1.0)      → 使用 /api/v1
├── Android App (v1.5)  → 使用 /api/v1
├── Web 前端 (v2.0)     → 使用 /api/v2
└── 第三方集成          → 使用 /api/v1
```

**没有版本控制的问题：**
- 更新 API 后，所有客户端必须同时更新
- 无法逐步迁移
- 旧版客户端可能无法使用

---

## 二、版本控制的好处

### 1. **向后兼容性** ✅

**允许旧版客户端继续工作：**

```
时间线：
2024-01-01: 发布 v1 API
2024-06-01: 发布 v2 API（新功能）
2024-12-01: 计划废弃 v1

在此期间：
- 旧客户端继续使用 /api/v1（不受影响）
- 新客户端使用 /api/v2（新功能）
- 逐步迁移，平滑过渡
```

### 2. **平滑升级** ✅

**给客户端充足的迁移时间：**

```
阶段 1: 发布 v2，v1 和 v2 共存
  ↓
阶段 2: 通知客户端 v1 将在 6 个月后废弃
  ↓
阶段 3: 提供迁移指南和工具
  ↓
阶段 4: 设置废弃日期，继续维护 v1
  ↓
阶段 5: 废弃 v1，只保留 v2
```

### 3. **降低风险** ✅

**避免"一刀切"的更新：**

- ✅ 可以逐步测试新版本
- ✅ 可以回滚到旧版本
- ✅ 减少生产环境故障

### 4. **清晰的变更管理** ✅

**版本号明确标识变更：**

```
/api/v1/users  → 稳定版本，不轻易改动
/api/v2/users  → 新版本，包含破坏性变更
/api/v3/users  → 未来版本
```

### 5. **支持 A/B 测试** ✅

**可以同时运行多个版本：**

```go
// 新功能在 v2 中测试
apiV2 := r.Group("/api/v2")
{
    users := apiV2.Group("/users")
    {
        users.GET("", userHandlerV2.ListUsers) // 新实现
    }
}

// v1 保持稳定
apiV1 := r.Group("/api/v1")
{
    users := apiV1.Group("/users")
    {
        users.GET("", userHandlerV1.ListUsers) // 稳定实现
    }
}
```

---

## 三、版本控制的方式

### 1. **URL 路径版本控制**（推荐，你项目中使用的方式）

```go
/api/v1/users
/api/v2/users
```

**优点：**
- ✅ 清晰直观
- ✅ 易于实现
- ✅ 缓存友好
- ✅ 符合 RESTful 原则

**实现方式：**
```go
// router/router.go
apiV1 := r.Group("/api/v1")
{
    users := apiV1.Group("/users")
    {
        users.GET("", userHandlerV1.ListUsers)
    }
}

apiV2 := r.Group("/api/v2")
{
    users := apiV2.Group("/users")
    {
        users.GET("", userHandlerV2.ListUsers)
    }
}
```

### 2. **请求头版本控制**

```
GET /api/users
Headers:
  Accept: application/vnd.api.v1+json
  或
  API-Version: v1
```

**优点：**
- ✅ URL 更简洁
- ✅ 版本信息集中管理

**缺点：**
- ❌ 不够直观
- ❌ 浏览器测试不方便
- ❌ 缓存可能有问题

### 3. **查询参数版本控制**

```
GET /api/users?version=v1
GET /api/users?v=2
```

**缺点：**
- ❌ 不符合 RESTful 原则
- ❌ 容易与业务查询参数混淆
- ❌ 不推荐使用

---

## 四、实际应用场景

### 场景 1：字段变更

**v1 版本：**
```json
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

**v2 版本（添加新字段）：**
```json
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "phone": "13800138000",  // 新增
  "avatar": "https://..."  // 新增
}
```

**实现：**
```go
// handler/user_v1.go
func (h *UserHandlerV1) GetUser(c *gin.Context) {
    user := // ... 查询用户
    c.JSON(200, gin.H{
        "id": user.ID,
        "name": user.Name,
        "email": user.Email,
        // v1 不返回 phone 和 avatar
    })
}

// handler/user_v2.go
func (h *UserHandlerV2) GetUser(c *gin.Context) {
    user := // ... 查询用户
    c.JSON(200, gin.H{
        "id": user.ID,
        "name": user.Name,
        "email": user.Email,
        "phone": user.Phone,    // v2 新增
        "avatar": user.Avatar,  // v2 新增
    })
}
```

### 场景 2：响应结构变更

**v1 版本（扁平结构）：**
```json
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "status": 1
}
```

**v2 版本（嵌套结构）：**
```json
{
  "id": 1,
  "name": "张三",
  "contact": {
    "email": "zhangsan@example.com"
  },
  "meta": {
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 场景 3：业务逻辑变更

**v1：简单分页**
```
GET /api/v1/users?page=1&page_size=10
```

**v2：高级查询**
```
GET /api/v2/users?page=1&page_size=10&status=1&role=admin&sort=created_at&order=desc
```

---

## 五、版本控制最佳实践

### 1. **版本命名规范**

```
主版本号.次版本号.修订号
v1.0.0  → 主版本（破坏性变更）
v1.1.0  → 次版本（新功能，向后兼容）
v1.1.1  → 修订版本（bug 修复）
```

**RESTful API 通常只使用主版本号：**
```
/api/v1
/api/v2
/api/v3
```

### 2. **何时创建新版本？**

**应该创建新版本（v2）的情况：**
- ✅ 删除字段
- ✅ 修改字段类型
- ✅ 修改字段名
- ✅ 修改响应结构
- ✅ 修改必需字段
- ✅ 修改业务逻辑（影响客户端）

**不需要创建新版本的情况：**
- ✅ 添加可选字段（向后兼容）
- ✅ 添加新的端点
- ✅ Bug 修复
- ✅ 性能优化

### 3. **版本生命周期管理**

```go
// 示例：版本管理中间件
func VersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.Param("version") // 从路径提取版本
        
        // 检查版本是否支持
        supportedVersions := []string{"v1", "v2"}
        if !contains(supportedVersions, version) {
            c.JSON(400, gin.H{
                "error": "不支持的 API 版本",
                "supported_versions": supportedVersions,
            })
            c.Abort()
            return
        }
        
        // 检查版本是否已废弃
        deprecatedVersions := map[string]string{
            "v1": "2025-12-31", // v1 将在 2025-12-31 废弃
        }
        if deprecationDate, ok := deprecatedVersions[version]; ok {
            c.Header("Deprecation", "true")
            c.Header("Sunset", deprecationDate)
            c.Header("Link", "</api/v2>; rel=\"successor-version\"")
        }
        
        c.Set("api_version", version)
        c.Next()
    }
}
```

### 4. **版本废弃策略**

**步骤 1：标记为废弃**
```http
GET /api/v1/users
Headers:
  Deprecation: true
  Sunset: 2025-12-31
  Link: </api/v2>; rel="successor-version"
```

**步骤 2：通知客户端**
- 在文档中明确标注
- 发送邮件通知
- 在响应头中添加警告

**步骤 3：设置废弃日期**
- 给客户端 6-12 个月的迁移时间
- 继续维护旧版本（bug 修复）

**步骤 4：正式废弃**
- 停止维护
- 返回 410 Gone 状态码

---

## 六、项目中的实现建议

### 当前实现（推荐保持）

```go
// router/router.go
api := r.Group("/api/v1")
{
    users := api.Group("/users")
    {
        users.POST("", userHandler.CreateUser)
        users.GET("", userHandler.ListUsers)
        users.GET("/:id", userHandler.GetUser)
        users.PUT("/:id", userHandler.UpdateUser)
        users.DELETE("/:id", userHandler.DeleteUser)
    }
}
```

### 未来扩展（多版本支持）

```go
// router/router.go
func SetupRouter(params RouterParams) *gin.Engine {
    // ...
    
    // v1 API（稳定版本）
    apiV1 := r.Group("/api/v1")
    {
        usersV1 := apiV1.Group("/users")
        {
            usersV1.POST("", userHandlerV1.CreateUser)
            usersV1.GET("", userHandlerV1.ListUsers)
            usersV1.GET("/:id", userHandlerV1.GetUser)
            usersV1.PUT("/:id", userHandlerV1.UpdateUser)
            usersV1.DELETE("/:id", userHandlerV1.DeleteUser)
        }
    }
    
    // v2 API（新版本，包含新功能）
    apiV2 := r.Group("/api/v2")
    {
        usersV2 := apiV2.Group("/users")
        {
            usersV2.POST("", userHandlerV2.CreateUser)
            usersV2.GET("", userHandlerV2.ListUsers)
            usersV2.GET("/:id", userHandlerV2.GetUser)
            usersV2.PATCH("/:id", userHandlerV2.PatchUser) // v2 使用 PATCH
            usersV2.DELETE("/:id", userHandlerV2.DeleteUser)
        }
    }
    
    return r
}
```

### 共享代码策略

```go
// 共享 Service 层逻辑
type UserService interface {
    CreateUser(name, email, password string) (*model.User, error)
    // ...
}

// v1 Handler 使用共享 Service
type UserHandlerV1 struct {
    userService service.UserService
}

// v2 Handler 也使用共享 Service，但响应格式不同
type UserHandlerV2 struct {
    userService service.UserService
}

func (h *UserHandlerV2) GetUser(c *gin.Context) {
    user, err := h.userService.GetUserByID(id)
    if err != nil {
        // ...
    }
    
    // v2 返回更丰富的响应
    c.JSON(200, gin.H{
        "data": gin.H{
            "user": user,
            "meta": gin.H{
                "version": "v2",
                "includes": []string{"profile", "settings"},
            },
        },
    })
}
```

---

## 七、总结

### 版本控制的核心价值

1. **向后兼容** - 保护现有客户端
2. **平滑升级** - 给客户端迁移时间
3. **降低风险** - 避免大规模故障
4. **灵活演进** - 支持业务需求变化
5. **清晰管理** - 明确标识变更

### 关键原则

- ✅ **向后兼容的变更** → 不需要新版本
- ✅ **破坏性变更** → 必须创建新版本
- ✅ **给客户端迁移时间** → 至少 6 个月
- ✅ **明确废弃策略** → 使用 HTTP 头通知
- ✅ **保持旧版本稳定** → 只修复 bug，不添加功能

### 你的项目

当前使用 `/api/v1` 是**最佳实践**，为未来的版本演进打下了良好基础。当需要做破坏性变更时，只需创建 `/api/v2` 即可。
