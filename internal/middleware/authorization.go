package middleware

import (
	"go_web/internal/service"
	"go_web/internal/util"

	"github.com/gin-gonic/gin"
)

// RequirePermission 权限校验中间件
// resource: 资源类型，如 "user", "role", "permission"
// action: 操作类型，如 "create", "read", "update", "delete"
func RequirePermission(userService service.UserService, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID（从认证中间件获取）
		userID, exists := c.Get("user_id")
		if !exists {
			util.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 转换用户ID类型
		var uid uint
		switch v := userID.(type) {
		case uint:
			uid = v
		case uint64:
			uid = uint(v)
		case int:
			if v > 0 {
				uid = uint(v)
			}
		default:
			util.Unauthorized(c, "用户ID格式错误")
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := userService.HasPermission(uid, resource, action)
		if err != nil {
			util.InternalServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			util.Forbidden(c, "权限不足，禁止访问")
			c.Abort()
			return
		}

		c.Next()
	}
}
