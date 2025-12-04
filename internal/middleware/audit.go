package middleware

import (
	"go_test/internal/logger"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware 审计中间件
func AuditMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求信息到审计日志
		auditInfo := map[string]interface{}{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"client_ip": c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// 如果有用户信息，可以在这里添加
		// if userID, exists := c.Get("user_id"); exists {
		// 	auditInfo["user_id"] = userID
		// }

		log.WithFields(auditInfo).Info("审计日志")

		c.Next()
	}
}

