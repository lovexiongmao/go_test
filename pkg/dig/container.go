package dig

import (
	"go_web/internal/config"
	"go_web/internal/database"
	"go_web/internal/handler"
	"go_web/internal/logger"
	"go_web/internal/middleware"
	"go_web/internal/model"
	"go_web/internal/repository"
	"go_web/internal/router"
	"go_web/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// Container 依赖注入容器
type Container struct {
	*dig.Container
}

// NewContainer 创建新的依赖注入容器
func NewContainer() *Container {
	c := dig.New()

	// 提供配置
	c.Provide(config.LoadConfig)

	// 提供请求日志Logger
	c.Provide(logger.NewLogger)

	// 提供审计日志Logger（使用命名参数区分）
	c.Provide(func(cfg *config.Config) *logger.Logger {
		return logger.NewAuditLogger(cfg)
	}, dig.Name("auditLogger"))

	// 提供数据库
	c.Provide(database.NewDatabase)

	// 提供Repository
	c.Provide(repository.NewUserRepository)
	c.Provide(repository.NewRoleRepository)
	c.Provide(repository.NewPermissionRepository)

	// 提供Service
	c.Provide(service.NewUserService)
	c.Provide(service.NewRoleService)
	c.Provide(service.NewPermissionService)

	// 提供Handler
	c.Provide(handler.NewUserHandler)
	c.Provide(handler.NewRoleHandler)
	c.Provide(handler.NewPermissionHandler)
	c.Provide(handler.NewAuthHandler)

	// 提供中间件（使用命名参数区分）
	c.Provide(func(log *logger.Logger) gin.HandlerFunc {
		return middleware.LoggerMiddleware(log)
	}, dig.Name("logger"))

	// 审计中间件参数结构体
	type AuditMiddlewareParams struct {
		dig.In
		AuditLogger *logger.Logger `name:"auditLogger"`
	}
	c.Provide(func(params AuditMiddlewareParams) gin.HandlerFunc {
		return middleware.AuditMiddleware(params.AuditLogger)
	}, dig.Name("audit"))

	// JWT 认证中间件
	c.Provide(func(cfg *config.Config) gin.HandlerFunc {
		return middleware.JWTAuthMiddleware(cfg)
	}, dig.Name("jwt"))

	// 提供路由
	c.Provide(router.SetupRouter)

	return &Container{Container: c}
}

// InitializeDatabase 初始化数据库表
// 注意：生产环境应禁用自动迁移，使用专门的迁移工具（如 golang-migrate）
func InitializeDatabase(db *gorm.DB, cfg *config.Config) error {
	// 只在配置允许时执行迁移
	if !cfg.Database.AutoMigrate {
		return nil
	}

	return db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&database.AuditLog{},
	)
}
