package main

import (
	"deantech/api"
	"deantech/config"
	"deantech/internal/database"
	"deantech/services"
	"fmt"
	"log"
)

func main() {
	// 加载配置，指定配置文件路径
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化服务
	telnetService := services.NewTelnetService()
	pingService := services.NewPingService()
	traceRouteService := services.NewTraceRouteService()
	prometheusService := services.NewPrometheusService(&cfg.Prometheus)
	userService := services.NewUserService(db)
	systemSettingService := services.NewSystemSettingService(db)
	hostService := services.NewHostService(db)
	sshService := services.NewSSHService(db)
	emailSettingService := services.NewEmailSettingService(db)
	emailService := services.NewEmailService(db)
	logService := services.NewLogService(db)
	clusterService := services.NewClusterService(db)
	kvmService := services.NewKVMService(db)
	dashboardService := services.NewDashboardService(db, cfg)
	slsConfigService := services.NewSLSConfigService(db)

	// 初始化API路由，使用所有可用服务
	router := api.SetupRouter(nil, nil, prometheusService, userService, hostService, sshService, clusterService, slsConfigService, systemSettingService, dashboardService, kvmService, emailService, emailSettingService, logService, nil, nil, nil, telnetService, pingService, traceRouteService, db)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
