package middleware

import (
	"context"

	"go_web/internal/database"
	"go_web/internal/logger"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware 审计中间件
func AuditMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求信息到审计日志
		auditInfo := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// 获取用户ID（如果有认证中间件设置）
		var userID uint
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint); ok {
				userID = id
				auditInfo["user_id"] = id
			} else if id, ok := uid.(uint64); ok {
				userID = uint(id)
				auditInfo["user_id"] = id
			} else if id, ok := uid.(int); ok && id > 0 {
				userID = uint(id)
				auditInfo["user_id"] = id
			}
		}

		// 将审计信息设置到 request context，供数据库审计插件使用
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, database.AuditUserIDKey, userID)
		ctx = context.WithValue(ctx, database.AuditIPKey, c.ClientIP())
		c.Request = c.Request.WithContext(ctx)

		log.WithFields(auditInfo).Info("审计日志")

		c.Next()
	}
}
