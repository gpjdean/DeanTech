package database

import (
	"deantech/config"
	"deantech/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// 配置日志级别
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// 获取底层的sql.DB对象，配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// 迁移所有需要的模型
	if err := db.AutoMigrate(
		&models.User{},
		&models.Alert{},
		&models.AlertMedia{},
		&models.AlertTemplate{},
		&models.Silence{},
		&models.InhibitionRule{},
		&models.Host{},
		&models.KVMHost{},
		&models.Cluster{},
		&models.SLSConfig{},
		&models.SystemSetting{},
		&models.EmailSetting{},
		&models.OperationLog{},
		&models.LoginLog{},
		&models.MD5History{},
		&models.ImageCompressHistory{},
	); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
		return nil, err
	}

	// 手动添加缺失的列（如果不存在）
	// Name列
	if err := db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(50) DEFAULT ''").Error; err != nil {
		log.Printf("Warning: Failed to add name column: %v", err)
	}
	// Phone列
	if err := db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20) DEFAULT ''").Error; err != nil {
		log.Printf("Warning: Failed to add phone column: %v", err)
	}
	// Department列
	if err := db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS department VARCHAR(50) DEFAULT ''").Error; err != nil {
		log.Printf("Warning: Failed to add department column: %v", err)
	}
	// Position列
	if err := db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS position VARCHAR(50) DEFAULT ''").Error; err != nil {
		log.Printf("Warning: Failed to add position column: %v", err)
	}
	// Avatar列
	if err := db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar TEXT DEFAULT NULL").Error; err != nil {
		log.Printf("Warning: Failed to add avatar column: %v", err)
	}
	// Remark列
	if err := db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS remark TEXT DEFAULT NULL").Error; err != nil {
		log.Printf("Warning: Failed to add remark column: %v", err)
	}

	// 初始化默认用户
	initDefaultUser(db)

	return db, nil
}

// initDefaultUser 初始化默认用户
func initDefaultUser(db *gorm.DB) {
	// 检查是否已有用户
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Default user already exists, skipping initialization")
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// 创建默认管理员用户
	defaultUser := models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@example.com",
		Role:     "admin",
	}

	if err := db.Create(&defaultUser).Error; err != nil {
		log.Fatalf("Failed to create default user: %v", err)
	}

	log.Println("Default user created successfully")
}
