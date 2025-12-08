package router

import (
	"go_test/internal/config"
	"go_test/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	loggerMiddleware gin.HandlerFunc,
	auditMiddleware gin.HandlerFunc,
	userHandler *handler.UserHandler,
) *gin.Engine {
	// 在创建路由之前设置Gin模式
	gin.SetMode(cfg.Server.Mode)
	
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

