package handler

import (
	"go_web/internal/config"
	"go_web/internal/service"
	"go_web/internal/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService service.UserService
	config      *config.Config
}

func NewAuthHandler(userService service.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:      cfg,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@example.com"` // 用户邮箱
	Password string `json:"password" binding:"required" example:"123456"`               // 用户密码
}

type LoginResponse struct {
	Token string      `json:"token"` // JWT Token
	User  interface{} `json:"user"`  // 用户信息
}

// Login 用户登录
// @Summary      用户登录
// @Description  用户登录接口，验证用户名密码后返回 JWT Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        login  body      LoginRequest  true  "登录信息"
// @Success      200    {object}  util.Response{data=LoginResponse}
// @Failure      400    {object}  util.Response
// @Failure      401    {object}  util.Response
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	// 1. 根据邮箱查找用户
	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		util.Unauthorized(c, "邮箱或密码错误")
		return
	}

	// 2. 验证密码（使用 bcrypt）
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		util.Unauthorized(c, "邮箱或密码错误")
		return
	}

	// 3. 检查用户状态
	if user.Status != 1 {
		util.Unauthorized(c, "用户已被禁用")
		return
	}

	// 4. 生成 JWT Token
	token, err := util.GenerateToken(h.config, user.ID, user.Email)
	if err != nil {
		util.InternalServerError(c, "生成 token 失败")
		return
	}

	// 5. 返回 token 和用户信息
	util.Success(c, LoginResponse{
		Token: token,
		User: gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
