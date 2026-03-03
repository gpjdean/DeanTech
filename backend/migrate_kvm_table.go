package main

import (
	"log"
	"deantech/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 使用正确的数据库连接
	dsn := "usernmame:password@tcp(127.0.0.1:3306)/deantech?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")

	// 1. 检查旧表是否存在
	var oldTableExists int64
	db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'deantech' AND table_name = 'k_vm_hosts'").Scan(&oldTableExists)
	
	// 2. 检查新表是否存在
	var newTableExists int64
	db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'deantech' AND table_name = 'kvm_hosts'").Scan(&newTableExists)

	log.Printf("Old table (k_vm_hosts) exists: %v", oldTableExists > 0)
	log.Printf("New table (kvm_hosts) exists: %v", newTableExists > 0)

	// 3. 迁移KVMHost模型（这将创建新表kvm_hosts）
	log.Println("Starting migration for KVMHost model with new table name...")
	if err := db.AutoMigrate(&models.KVMHost{}); err != nil {
		log.Fatalf("Failed to migrate KVMHost: %v", err)
	}

	log.Println("KVMHost migration completed successfully")

	// 4. 如果旧表存在且新表不存在旧数据，尝试迁移数据
	if oldTableExists > 0 && newTableExists == 0 {
		log.Println("Migrating data from old table to new table...")
		// 注意：这只是一个简单的迁移示例，实际情况可能需要更复杂的处理
		result := db.Exec(`INSERT INTO kvm_hosts SELECT * FROM k_vm_hosts`)
		if result.Error != nil {
			log.Printf("Warning: Failed to migrate data from old table: %v", result.Error)
		} else {
			log.Printf("Successfully migrated %d rows from old table to new table", result.RowsAffected)
			
			// 可以选择删除旧表
			// log.Println("Dropping old table...")
			// if err := db.Exec(`DROP TABLE k_vm_hosts`).Error; err != nil {
			// 	log.Printf("Warning: Failed to drop old table: %v", err)
			// } else {
			// 	log.Println("Successfully dropped old table")
			// }
		}
	}

	// 5. 验证新表是否成功创建
	db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'deantech' AND table_name = 'kvm_hosts'").Scan(&newTableExists)
	if newTableExists > 0 {
		log.Println("\nkvm_hosts table created successfully!")
		
		// 查看新表结构
		log.Println("New table structure:")
		var columns []struct {
			Field   string
			Type    string
			Null    string
			Key     string
			Default interface{}
			Extra   string
		}
		db.Raw("SHOW COLUMNS FROM kvm_hosts").Scan(&columns)
		for _, col := range columns {
			log.Printf("  %s: %s %s %s", col.Field, col.Type, col.Null, col.Key)
		}
	} else {
		log.Fatalf("Failed to create kvm_hosts table")
	}

	log.Println("\nMigration completed successfully!")
}
