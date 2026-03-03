package main

import (
	"deantech/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	// 使用正确的数据库连接字符串
	dsn := "deantech:deantech@tcp(43.143.194.230:3306)/deantech?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移模型
	db.AutoMigrate(&models.User{})

	// 查找admin用户
	var admin models.User
	result := db.Where("username = ?", "admin").First(&admin)
	if result.Error != nil {
		log.Fatalf("Failed to find admin user: %v", result.Error)
	}

	// 设置新密码
	newPassword := "handsome"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// 更新密码
	admin.Password = string(hashedPassword)
	if err := db.Save(&admin).Error; err != nil {
		log.Fatalf("Failed to update admin password: %v", err)
	}

	log.Printf("Admin password reset successfully. New password: %s", newPassword)
}
