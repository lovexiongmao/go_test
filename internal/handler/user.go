package handler

import (
	"strconv"
	"time"

	"go_web/internal/service"
	"go_web/internal/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required" example:"张三"`                          // 用户姓名
	Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"` // 用户邮箱
	Password string `json:"password" binding:"required,min=6" example:"123456"`            // 用户密码（最少6位）
}

type UpdateUserRequest struct {
	Name   *string `json:"name" example:"李四"`  // 用户姓名（可选）
	Status *int    `json:"status" example:"1"` // 用户状态：1-正常，0-禁用（可选）
}

// UserResponse 用户响应结构体（用于 Swagger 文档）
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`                            // 用户ID
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"` // 更新时间
	Name      string    `json:"name" example:"张三"`                         // 用户姓名
	Email     string    `json:"email" example:"zhangsan@example.com"`      // 用户邮箱
	Status    int       `json:"status" example:"1"`                        // 用户状态：1-正常，0-禁用
}

// CreateUser 创建用户
// @Summary      创建用户
// @Description  创建一个新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        Authorization header    string             true  "Bearer {token}"  default(Bearer )
// @Param        user          body      CreateUserRequest  true  "用户信息"
// @Success      201           {object}  util.Response{data=UserResponse}
// @Failure      400           {object}  util.Response
// @Failure      401           {object}  util.Response
// @Failure      500           {object}  util.Response
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		util.InternalServerErrorWithError(c, "创建用户失败", err)
		return
	}

	util.CreatedWithMessage(c, "用户创建成功", user)
}

// GetUser 获取用户详情
// @Summary      获取用户详情
// @Description  根据用户ID获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id            path      int     true  "用户ID"
// @Param        Authorization header    string  true  "Bearer {token}"  default(Bearer )
// @Success      200           {object}  util.Response{data=UserResponse}
// @Failure      400           {object}  util.Response
// @Failure      401           {object}  util.Response
// @Failure      404           {object}  util.Response
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		util.NotFound(c, "用户不存在")
		return
	}

	util.Success(c, user)
}

// UpdateUser 更新用户
// @Summary      更新用户
// @Description  根据用户ID更新用户信息（支持部分更新）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id            path      int                true  "用户ID"
// @Param        Authorization header    string             true  "Bearer {token}"  default(Bearer )
// @Param        user          body      UpdateUserRequest  true  "用户信息"
// @Success      200           {object}  util.Response{data=UserResponse}
// @Failure      400           {object}  util.Response
// @Failure      401           {object}  util.Response
// @Failure      500           {object}  util.Response
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的用户ID")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	// 处理指针参数
	var name string
	var status int = -1 // -1表示不更新
	if req.Name != nil {
		name = *req.Name
	}
	if req.Status != nil {
		status = *req.Status
	}

	user, err := h.userService.UpdateUser(uint(id), name, status)
	if err != nil {
		util.InternalServerErrorWithError(c, "更新用户失败", err)
		return
	}

	util.SuccessWithMessage(c, "用户更新成功", user)
}

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  根据用户ID删除用户（软删除）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id            path      int     true  "用户ID"
// @Param        Authorization header    string  true  "Bearer {token}"  default(Bearer )
// @Success      200           {object}  util.Response
// @Failure      400           {object}  util.Response
// @Failure      401           {object}  util.Response
// @Failure      500           {object}  util.Response
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的用户ID")
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		util.InternalServerErrorWithError(c, "删除用户失败", err)
		return
	}

	util.SuccessWithMessage(c, "用户删除成功", nil)
}

// ListUsers 用户列表
// @Summary      获取用户列表
// @Description  分页获取用户列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        page          query     int     false  "页码"      default(1)
// @Param        page_size     query     int     false  "每页数量"   default(10)
// @Param        Authorization header    string  false  "Bearer {token}"  default(Bearer )
// @Success      200           {object}  util.Response{data=object{list=[]UserResponse,total=int,page=int,page_size=int}}
// @Failure      401           {object}  util.Response
// @Failure      500           {object}  util.Response
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.userService.ListUsers(page, pageSize)
	if err != nil {
		util.InternalServerErrorWithError(c, "获取用户列表失败", err)
		return
	}

	util.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
