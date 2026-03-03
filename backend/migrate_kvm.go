package main

import (
	"log"
	"deantech/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 使用与原项目相同的数据库连接
	dsn := "deantech:deantech@tcp(43.143.194.230:3306)/deantech?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")

	// 只迁移KVMHost模型
	log.Println("Starting migration for KVMHost model...")
	if err := db.AutoMigrate(&models.KVMHost{}); err != nil {
		log.Fatalf("Failed to migrate KVMHost: %v", err)
	}

	log.Println("KVMHost migration completed successfully")

	// 检查表是否存在
	var count int64
	db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'deantech' AND table_name = 'k_vm_hosts'").Scan(&count)
	if count > 0 {
		log.Println("k_vm_hosts table exists")
		
		// 查看表结构
		log.Println("Table structure:")
		var columns []struct {
			ColumnName string
			DataType   string
		}
		db.Raw("SHOW COLUMNS FROM k_vm_hosts").Scan(&columns)
		for _, col := range columns {
			log.Printf("  %s: %s", col.ColumnName, col.DataType)
		}
	} else {
		log.Println("k_vm_hosts table does NOT exist")
	}
}
