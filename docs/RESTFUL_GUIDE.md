# RESTful API 设计指南

## 一、项目中 RESTful 风格的体现

### 1. **HTTP 方法正确使用** ✅

项目中正确使用了 HTTP 方法来表示不同的操作：

```go
// router/router.go
users.POST("", userHandler.CreateUser)        // POST   - 创建资源
users.GET("", userHandler.ListUsers)          // GET    - 获取资源列表
users.GET("/:id", userHandler.GetUser)       // GET    - 获取单个资源
users.PUT("/:id", userHandler.UpdateUser)    // PUT    - 更新资源
users.DELETE("/:id", userHandler.DeleteUser) // DELETE - 删除资源
```

**对应关系：**
- `POST /api/v1/users` → 创建新用户
- `GET /api/v1/users` → 获取用户列表
- `GET /api/v1/users/:id` → 获取指定用户
- `PUT /api/v1/users/:id` → 更新指定用户
- `DELETE /api/v1/users/:id` → 删除指定用户

### 2. **资源路径设计** ✅

- 使用名词复数形式：`/users` 而不是 `/user`
- 使用层级结构：`/api/v1/users`（版本控制）
- 使用路径参数：`/:id` 表示资源标识符

### 3. **HTTP 状态码使用** ✅

```go
// handler/user.go
c.JSON(http.StatusCreated, ...)    // 201 - 创建成功
c.JSON(http.StatusOK, ...)          // 200 - 操作成功
c.JSON(http.StatusBadRequest, ...)  // 400 - 请求错误
c.JSON(http.StatusNotFound, ...)    // 404 - 资源不存在
c.JSON(http.StatusInternalServerError, ...) // 500 - 服务器错误
```

### 4. **JSON 格式响应** ✅

统一使用 JSON 格式进行数据交换。

---

## 二、RESTful 设计原则

### 核心原则

1. **资源（Resource）**：URL 应该表示资源，而不是动作
   - ✅ `/api/v1/users` - 用户资源
   - ❌ `/api/v1/createUser` - 动作

2. **HTTP 方法表示操作**：
   - `GET` - 查询（幂等）
   - `POST` - 创建（非幂等）
   - `PUT` - 完整更新（幂等）
   - `PATCH` - 部分更新（幂等）
   - `DELETE` - 删除（幂等）

3. **无状态（Stateless）**：每个请求都应该包含所有必要信息

4. **统一接口**：使用标准的 HTTP 方法和状态码

---

## 三、如何保持 RESTful 风格

### 1. **添加新资源时的模板**

#### 步骤 1：创建 Model
```go
// internal/model/article.go
type Article struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
    
    Title   string `gorm:"type:varchar(200);not null" json:"title"`
    Content string `gorm:"type:text" json:"content"`
    UserID  uint   `gorm:"not null;index" json:"user_id"`
}
```

#### 步骤 2：创建 Repository
```go
// internal/repository/article.go
type ArticleRepository interface {
    Create(article *model.Article) error
    GetByID(id uint) (*model.Article, error)
    List(offset, limit int) ([]*model.Article, int64, error)
    Update(article *model.Article) error
    Delete(id uint) error
}
```

#### 步骤 3：创建 Service
```go
// internal/service/article.go
type ArticleService interface {
    CreateArticle(title, content string, userID uint) (*model.Article, error)
    GetArticleByID(id uint) (*model.Article, error)
    ListArticles(page, pageSize int) ([]*model.Article, int64, error)
    UpdateArticle(id uint, title, content string) (*model.Article, error)
    DeleteArticle(id uint) error
}
```

#### 步骤 4：创建 Handler
```go
// internal/handler/article.go
type ArticleHandler struct {
    articleService service.ArticleService
}

// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var req CreateArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // ... 业务逻辑
    c.JSON(http.StatusCreated, gin.H{
        "message": "文章创建成功",
        "data":    article,
    })
}

// GetArticle 获取文章详情
func (h *ArticleHandler) GetArticle(c *gin.Context) {
    // ...
    c.JSON(http.StatusOK, gin.H{"data": article})
}

// ListArticles 获取文章列表
func (h *ArticleHandler) ListArticles(c *gin.Context) {
    // ...
    c.JSON(http.StatusOK, gin.H{"data": gin.H{
        "list": articles,
        "total": total,
    }})
}

// UpdateArticle 更新文章
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
    // ...
    c.JSON(http.StatusOK, gin.H{
        "message": "文章更新成功",
        "data":    article,
    })
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
    // ...
    c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
```

#### 步骤 5：注册路由
```go
// internal/router/router.go
api := r.Group("/api/v1")
{
    // 用户路由
    users := api.Group("/users")
    {
        users.POST("", userHandler.CreateUser)
        users.GET("", userHandler.ListUsers)
        users.GET("/:id", userHandler.GetUser)
        users.PUT("/:id", userHandler.UpdateUser)
        users.DELETE("/:id", userHandler.DeleteUser)
    }
    
    // 文章路由（新增）
    articles := api.Group("/articles")
    {
        articles.POST("", articleHandler.CreateArticle)
        articles.GET("", articleHandler.ListArticles)
        articles.GET("/:id", articleHandler.GetArticle)
        articles.PUT("/:id", articleHandler.UpdateArticle)
        articles.DELETE("/:id", articleHandler.DeleteArticle)
    }
}
```

### 2. **路由命名规范**

#### ✅ 好的实践：
```
POST   /api/v1/users              - 创建用户
GET    /api/v1/users              - 获取用户列表
GET    /api/v1/users/:id         - 获取用户详情
PUT    /api/v1/users/:id         - 更新用户
DELETE /api/v1/users/:id         - 删除用户

POST   /api/v1/articles          - 创建文章
GET    /api/v1/articles          - 获取文章列表
GET    /api/v1/articles/:id      - 获取文章详情
PUT    /api/v1/articles/:id      - 更新文章
DELETE /api/v1/articles/:id      - 删除文章
```

#### ❌ 避免的做法：
```
POST   /api/v1/createUser        - 动作在 URL 中
GET    /api/v1/getUser/:id       - 动作在 URL 中
POST   /api/v1/user/update/:id   - 动作在 URL 中
POST   /api/v1/user/delete/:id  - 应该用 DELETE 方法
```

### 3. **嵌套资源处理**

当资源有层级关系时：

```
GET    /api/v1/users/:user_id/articles        - 获取某用户的所有文章
POST   /api/v1/users/:user_id/articles        - 为用户创建文章
GET    /api/v1/users/:user_id/articles/:id    - 获取用户的某篇文章
PUT    /api/v1/users/:user_id/articles/:id    - 更新用户的文章
DELETE /api/v1/users/:user_id/articles/:id    - 删除用户的文章
```

**实现示例：**
```go
// router/router.go
users := api.Group("/users")
{
    users.GET("/:user_id/articles", articleHandler.ListUserArticles)
    users.POST("/:user_id/articles", articleHandler.CreateUserArticle)
    users.GET("/:user_id/articles/:id", articleHandler.GetUserArticle)
    users.PUT("/:user_id/articles/:id", articleHandler.UpdateUserArticle)
    users.DELETE("/:user_id/articles/:id", articleHandler.DeleteUserArticle)
}
```

### 4. **查询参数规范**

对于列表查询，使用查询参数：

```
GET /api/v1/users?page=1&page_size=10&status=1&name=张三
```

```go
// handler/user.go
func (h *UserHandler) ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    status := c.Query("status")
    name := c.Query("name")
    // ...
}
```

### 5. **HTTP 状态码使用规范**

| 状态码 | 含义 | 使用场景 |
|--------|------|----------|
| 200 OK | 成功 | GET、PUT、PATCH 成功 |
| 201 Created | 创建成功 | POST 创建资源成功 |
| 204 No Content | 成功无内容 | DELETE 成功 |
| 400 Bad Request | 请求错误 | 参数验证失败 |
| 401 Unauthorized | 未授权 | 需要登录 |
| 403 Forbidden | 禁止访问 | 权限不足 |
| 404 Not Found | 资源不存在 | 资源未找到 |
| 409 Conflict | 冲突 | 资源冲突（如邮箱已存在） |
| 422 Unprocessable Entity | 无法处理 | 业务逻辑错误 |
| 500 Internal Server Error | 服务器错误 | 服务器内部错误 |

### 6. **响应格式统一**

建议统一响应格式：

```go
// 成功响应
{
    "code": 200,
    "message": "操作成功",
    "data": { ... }
}

// 错误响应
{
    "code": 400,
    "message": "错误信息",
    "error": "详细错误描述"
}
```

可以创建一个统一的响应工具：

```go
// internal/utils/response.go
package utils

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, gin.H{
        "code": 200,
        "message": "操作成功",
        "data": data,
    })
}

func Created(c *gin.Context, data interface{}) {
    c.JSON(http.StatusCreated, gin.H{
        "code": 201,
        "message": "创建成功",
        "data": data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{
        "code": code,
        "message": message,
    })
}
```

### 7. **PUT vs PATCH**

- **PUT**：完整更新，需要提供所有字段
- **PATCH**：部分更新，只需提供要更新的字段

**建议：**
- 使用 `PATCH` 进行部分更新（更符合 RESTful 规范）
- 如果使用 `PUT`，应该要求提供完整资源

```go
// 使用 PATCH 进行部分更新
articles.PATCH("/:id", articleHandler.PatchArticle)

// 使用 PUT 进行完整更新
articles.PUT("/:id", articleHandler.PutArticle)
```

---

## 四、最佳实践总结

### ✅ DO（应该做）

1. **使用复数名词**：`/users` 而不是 `/user`
2. **使用 HTTP 方法**：GET、POST、PUT、PATCH、DELETE
3. **使用路径参数**：`/:id` 表示资源标识
4. **使用查询参数**：用于过滤、分页、排序
5. **正确使用状态码**：200、201、400、404、500 等
6. **版本控制**：`/api/v1/`、`/api/v2/`
7. **统一响应格式**：保持一致的 JSON 结构

### ❌ DON'T（不应该做）

1. **不要在 URL 中使用动词**：`/createUser` ❌
2. **不要使用单数名词**：`/user` ❌（除非是单例资源）
3. **不要混用 HTTP 方法**：用 POST 删除 ❌
4. **不要在 URL 中使用动作**：`/user/update/:id` ❌
5. **不要忽略状态码**：所有请求都返回 200 ❌

---

## 五、扩展功能示例

### 示例：添加文章（Article）资源

1. **创建 Model** → `internal/model/article.go`
2. **创建 Repository** → `internal/repository/article.go`
3. **创建 Service** → `internal/service/article.go`
4. **创建 Handler** → `internal/handler/article.go`
5. **注册路由** → `internal/router/router.go` 中添加：
   ```go
   articles := api.Group("/articles")
   {
       articles.POST("", articleHandler.CreateArticle)
       articles.GET("", articleHandler.ListArticles)
       articles.GET("/:id", articleHandler.GetArticle)
       articles.PUT("/:id", articleHandler.UpdateArticle)
       articles.DELETE("/:id", articleHandler.DeleteArticle)
   }
   ```

遵循这个模式，可以保持整个项目的 RESTful 风格一致性。
