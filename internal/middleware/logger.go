package middleware

import (
	"time"

	"go_test/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// 收集所有类型的错误信息
		var errorMessages []string
		if errs := c.Errors.ByType(gin.ErrorTypePrivate); len(errs) > 0 {
			for _, err := range errs {
				errorMessages = append(errorMessages, err.Error())
			}
		}
		if errs := c.Errors.ByType(gin.ErrorTypePublic); len(errs) > 0 {
			for _, err := range errs {
				errorMessages = append(errorMessages, err.Error())
			}
		}
		if errs := c.Errors.ByType(gin.ErrorTypeBind); len(errs) > 0 {
			for _, err := range errs {
				errorMessages = append(errorMessages, err.Error())
			}
		}
		if errs := c.Errors.ByType(gin.ErrorTypeAny); len(errs) > 0 {
			for _, err := range errs {
				errorMessages = append(errorMessages, err.Error())
			}
		}

		if raw != "" {
			path = path + "?" + raw
		}

		entry := log.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"user_agent": c.Request.UserAgent(),
		})

		// 如果有错误信息，记录到日志中
		if len(errorMessages) > 0 {
			errorMsg := ""
			for i, msg := range errorMessages {
				if i > 0 {
					errorMsg += "; "
				}
				errorMsg += msg
			}
			entry = entry.WithField("error", errorMsg)
		}

		if statusCode >= 500 {
			entry.Error("HTTP请求错误")
		} else if statusCode >= 400 {
			entry.Warn("HTTP请求警告")
		} else {
			entry.Info("HTTP请求")
		}
	}
}
