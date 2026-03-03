package main

import (
	"log"
	"deantech/config"
	"deantech/internal/database"
	"deantech/models"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 专门迁移SLSConfig模型
	log.Println("Creating sls_configs table...")
	if err := db.AutoMigrate(&models.SLSConfig{}); err != nil {
		log.Fatalf("Failed to migrate SLSConfig: %v", err)
	}

	log.Println("sls_configs table created successfully!")
}