package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// 数据库配置
	oldDB := "promalert"
	newDB := "deantech"
	dbUser := "root"
	dbPassword := "password" // 请替换为实际的root密码
	dbHost := "43.143.194.230"
	dbPort := "3306"
	
	// 1. 创建新数据库
	log.Println("1. 创建新数据库")
	createDB := exec.Command(
		"mysql",
		"-u", dbUser,
		"-p"+dbPassword,
		"-h", dbHost,
		"-P", dbPort,
		"-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", newDB),
	)
	createDB.Stdout = os.Stdout
	createDB.Stderr = os.Stderr
	if err := createDB.Run(); err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}
	
	// 2. 创建新用户并授权
	log.Println("2. 创建新用户并授权")
	createUser := exec.Command(
		"mysql",
		"-u", dbUser,
		"-p"+dbPassword,
		"-h", dbHost,
		"-P", dbPort,
		"-e", fmt.Sprintf(`
		CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s';
		GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';
		GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost';
		FLUSH PRIVILEGES;
		`, newDB, newDB, newDB, newDB, newDB, newDB),
	)
	createUser.Stdout = os.Stdout
	createUser.Stderr = os.Stderr
	if err := createUser.Run(); err != nil {
		log.Fatalf("创建用户并授权失败: %v", err)
	}
	
	// 3. 导出旧数据库数据
	log.Println("3. 导出旧数据库数据")
	dumpCmd := exec.Command(
		"mysqldump",
		"-u", dbUser,
		"-p"+dbPassword,
		"-h", dbHost,
		"-P", dbPort,
		oldDB,
	)
	
	// 4. 导入数据到新数据库
	log.Println("4. 导入数据到新数据库")
	importCmd := exec.Command(
		"mysql",
		"-u", dbUser,
		"-p"+dbPassword,
		"-h", dbHost,
		"-P", dbPort,
		newDB,
	)
	
	// 连接管道
	importCmd.Stdin, _ = dumpCmd.StdoutPipe()
	importCmd.Stdout = os.Stdout
	importCmd.Stderr = os.Stderr
	
	// 启动导入命令
	if err := importCmd.Start(); err != nil {
		log.Fatalf("启动导入命令失败: %v", err)
	}
	
	// 启动导出命令
	if err := dumpCmd.Run(); err != nil {
		log.Fatalf("导出数据失败: %v", err)
	}
	
	// 等待导入命令完成
	if err := importCmd.Wait(); err != nil {
		log.Fatalf("导入数据失败: %v", err)
	}
	
	log.Println("数据库迁移完成！")
	log.Printf("旧数据库: %s", oldDB)
	log.Printf("新数据库: %s", newDB)
	log.Printf("新用户: %s", newDB)
	log.Printf("新密码: %s", newDB)
}
