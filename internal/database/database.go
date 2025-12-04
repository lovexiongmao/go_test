package database

import (
	"time"

	"go_test/internal/config"
	"go_test/internal/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.Config, log *logger.Logger) (*gorm.DB, error) {
	// 配置GORM日志
	var gormLog gormLogger.Interface
	if cfg.Log.Level == "debug" {
		gormLog = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		gormLog = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 注册自定义audit插件（可选，HTTP层面的审计已在中间件中实现）
	// auditPlugin := NewAuditPlugin(db)
	// if err := db.Use(auditPlugin); err != nil {
	// 	return nil, err
	// }

	log.Info("数据库连接成功")

	return db, nil
}

