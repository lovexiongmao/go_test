package router

import (
	"go_test/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	loggerMiddleware gin.HandlerFunc,
	auditMiddleware gin.HandlerFunc,
	userHandler *handler.UserHandler,
) *gin.Engine {
	// 注意：Gin模式应该在main.go中通过gin.SetMode()设置
	// 但gin.SetMode是全局设置，在路由创建后设置也能生效
	r := gin.Default()

	// 全局中间件
	r.Use(loggerMiddleware)
	r.Use(auditMiddleware)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "服务运行正常",
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return r
}

