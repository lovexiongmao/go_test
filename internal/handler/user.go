package handler

import (
	"strconv"

	"go_test/internal/service"
	"go_test/internal/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name   *string `json:"name"`   // 使用指针，nil表示不更新
	Status *int    `json:"status"` // 使用指针，nil表示不更新
}

// CreateUser 创建用户
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
