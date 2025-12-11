package router

import (
	"go_web/docs/swagger" // Swagger 文档
	"go_web/internal/config"
	"go_web/internal/handler"
	"go_web/internal/middleware"
	"go_web/internal/service"
	"go_web/internal/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/dig"
)

// RouterParams 路由参数结构体，用于接收命名参数
type RouterParams struct {
	dig.In

	Config            *config.Config
	LoggerMiddleware  gin.HandlerFunc `name:"logger"`
	AuditMiddleware   gin.HandlerFunc `name:"audit"`
	JWTAuthMiddleware gin.HandlerFunc `name:"jwt"`
	UserHandler       *handler.UserHandler
	RoleHandler       *handler.RoleHandler
	PermissionHandler *handler.PermissionHandler
	AuthHandler       *handler.AuthHandler
	UserService       service.UserService
}

func SetupRouter(params RouterParams) *gin.Engine {
	cfg := params.Config
	loggerMiddleware := params.LoggerMiddleware
	auditMiddleware := params.AuditMiddleware
	jwtAuthMiddleware := params.JWTAuthMiddleware
	userHandler := params.UserHandler
	roleHandler := params.RoleHandler
	permissionHandler := params.PermissionHandler
	authHandler := params.AuthHandler
	userService := params.UserService
	// 在创建路由之前设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()

	// 全局中间件
	r.Use(loggerMiddleware)
	r.Use(auditMiddleware)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		util.SuccessWithMessage(c, "服务运行正常", gin.H{
			"status": "ok",
		})
	})

	// Swagger UI 文档
	// 确保 Swagger 文档被注册（导入 docs/swagger 包会自动执行 init 函数）
	_ = swagger.SwaggerInfo
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	api := r.Group("/api/v1")
	{
		// 登录接口（不需要认证）
		api.POST("/login", authHandler.Login)

		// 需要认证的路由组
		auth := api.Group("")
		auth.Use(jwtAuthMiddleware) // 添加 JWT 认证中间件
		{
			// 用户相关路由
			users := auth.Group("/users")
			{
				users.POST("", middleware.RequirePermission(userService, "user", "create"), userHandler.CreateUser)
				users.GET("", middleware.RequirePermission(userService, "user", "read"), userHandler.ListUsers)
				users.GET("/:id", middleware.RequirePermission(userService, "user", "read"), userHandler.GetUser)
				users.PUT("/:id", middleware.RequirePermission(userService, "user", "update"), userHandler.UpdateUser)
				users.DELETE("/:id", middleware.RequirePermission(userService, "user", "delete"), userHandler.DeleteUser)
			}

			// 角色相关路由
			roles := auth.Group("/roles")
			{
				roles.POST("", middleware.RequirePermission(userService, "role", "create"), roleHandler.CreateRole)
				roles.GET("", middleware.RequirePermission(userService, "role", "read"), roleHandler.ListRoles)
				roles.GET("/:id", middleware.RequirePermission(userService, "role", "read"), roleHandler.GetRole)
				roles.PUT("/:id", middleware.RequirePermission(userService, "role", "update"), roleHandler.UpdateRole)
				roles.DELETE("/:id", middleware.RequirePermission(userService, "role", "delete"), roleHandler.DeleteRole)
				// 角色权限管理
				roles.POST("/:id/permissions", middleware.RequirePermission(userService, "role", "update"), roleHandler.AssignPermissions)
				roles.DELETE("/:id/permissions", middleware.RequirePermission(userService, "role", "update"), roleHandler.RemovePermissions)
				roles.GET("/:id/permissions", middleware.RequirePermission(userService, "role", "read"), roleHandler.GetRolePermissions)
				// 角色用户管理
				roles.POST("/:id/users", middleware.RequirePermission(userService, "role", "update"), roleHandler.AssignUsers)
				roles.DELETE("/:id/users", middleware.RequirePermission(userService, "role", "update"), roleHandler.RemoveUsers)
				roles.GET("/:id/users", middleware.RequirePermission(userService, "role", "read"), roleHandler.GetRoleUsers)
			}

			// 权限相关路由
			permissions := auth.Group("/permissions")
			{
				permissions.POST("", middleware.RequirePermission(userService, "permission", "create"), permissionHandler.CreatePermission)
				permissions.GET("", middleware.RequirePermission(userService, "permission", "read"), permissionHandler.ListPermissions)
				permissions.GET("/:id", middleware.RequirePermission(userService, "permission", "read"), permissionHandler.GetPermission)
				permissions.PUT("/:id", middleware.RequirePermission(userService, "permission", "update"), permissionHandler.UpdatePermission)
				permissions.DELETE("/:id", middleware.RequirePermission(userService, "permission", "delete"), permissionHandler.DeletePermission)
			}
		}
	}

	return r
}
