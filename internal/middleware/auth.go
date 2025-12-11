package middleware

import (
	"strings"

	"go_web/internal/config"
	"go_web/internal/util"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Authorization header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next() // 不设置 user_id，让后续中间件处理（可能返回未登录）
			return
		}

		// 2. 解析 Bearer token
		// 格式：Authorization: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]

		// 3. 解析 JWT token
		claims, err := util.ParseToken(cfg, token)
		if err != nil {
			c.Next() // token 无效，让后续中间件处理
			return
		}

		// 4. 将用户ID设置到 context 中
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
