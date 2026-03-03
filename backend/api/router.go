package api

import (
	"deantech/internal/middleware"
	"deantech/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(alertService *services.AlertService, mediaService *services.MediaService, prometheusService *services.PrometheusService, userService *services.UserService, hostService *services.HostService, sshService *services.SSHService, clusterService *services.ClusterService, slsConfigService *services.SLSConfigService, systemSettingService *services.SystemSettingService, dashboardService *services.DashboardService, kvmService *services.KVMService, emailService *services.EmailService, emailSettingService *services.EmailSettingService, logService *services.LogService, aiService *services.AiService, md5Service *services.MD5Service, imageCompressService *services.ImageCompressService, telnetService *services.TelnetService, pingService *services.PingService, traceRouteService *services.TraceRouteService, db *gorm.DB) *gin.Engine {
	r := gin.Default()
	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Username")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API分组
	api := r.Group("/api")
	// 应用操作日志中间件，只有当db非nil时才应用
	if db != nil {
		api.Use(middleware.OperationLogger(db))
	}
	{
		// Prometheus相关 - 总是注册，因为prometheusService在main.go中初始化了
		prometheus := api.Group("/prometheus")
		{
			prometheus.GET("/alerts", prometheusService.GetPrometheusAlerts)
			prometheus.GET("/rules", prometheusService.GetPrometheusRules)
			prometheus.GET("/targets", prometheusService.GetPrometheusTargets)
		}

		// 只有非nil的服务才注册路由
		if alertService != nil {
			// 告警相关
			alerts := api.Group("/alerts")
			{
				alerts.GET("", alertService.ListAlerts)
				alerts.GET("/:id", alertService.GetAlert)
				alerts.PUT("/:id/resolve", alertService.ResolveAlert)
				alerts.DELETE("/:id", alertService.DeleteAlert)
			}

			// 静默规则相关
			silences := api.Group("/silences")
			{
				silences.GET("", alertService.ListSilences)
				silences.POST("", alertService.CreateSilence)
				silences.GET("/:id", alertService.GetSilence)
				silences.DELETE("/:id", alertService.DeleteSilence)
			}

			// 抑制规则相关
			inhibitions := api.Group("/inhibitions")
			{
				inhibitions.GET("", alertService.ListInhibitions)
				inhibitions.POST("", alertService.CreateInhibition)
				inhibitions.GET("/:id", alertService.GetInhibition)
				inhibitions.PUT("/:id", alertService.UpdateInhibition)
				inhibitions.DELETE("/:id", alertService.DeleteInhibition)
			}
		}

		if mediaService != nil {
			// 告警介质相关
			media := api.Group("/media")
			{
				media.GET("", mediaService.ListMedia)
				media.POST("", mediaService.CreateMedia)
				media.GET("/:id", mediaService.GetMedia)
				media.PUT("/:id", mediaService.UpdateMedia)
				media.DELETE("/:id", mediaService.DeleteMedia)
			}

			// 告警模板相关
			templates := api.Group("/templates")
			{
				templates.GET("", mediaService.ListTemplates)
				templates.POST("", mediaService.CreateTemplate)
				templates.GET("/:id", mediaService.GetTemplate)
				templates.PUT("/:id", mediaService.UpdateTemplate)
				templates.DELETE("/:id", mediaService.DeleteTemplate)
			}
		}

		if userService != nil {
			// 用户相关
			users := api.Group("/users")
			{
				users.GET("", userService.ListUsersHandler)
				users.POST("", userService.CreateUserHandler)
				users.GET("/current", userService.GetCurrentUserHandler)
				users.POST("/register", userService.RegisterHandler)
				users.POST("/login", userService.LoginHandler)
				users.GET("/:id", userService.GetUserHandler)
				users.PUT("/:id", userService.UpdateUserHandler)
				users.PUT("/profile", userService.UpdateProfileHandler)
				users.POST("/change-password", userService.ChangePasswordHandler)
				users.DELETE("/:id", userService.DeleteUserHandler)

				// 忘记密码相关
				forgotPassword := users.Group("/forgot-password")
				{
					forgotPassword.POST("/verify-email", userService.VerifyEmailHandler)
					forgotPassword.POST("/reset", userService.ResetPasswordHandler)
				}
			}
		}

		if hostService != nil && sshService != nil {
			// 主机相关
			hosts := api.Group("/hosts")
			{
				hosts.GET("", hostService.ListHostsHandler)
				hosts.POST("", hostService.CreateHostHandler)
				hosts.GET("/:id", hostService.GetHostHandler)
				hosts.PUT("/:id", hostService.UpdateHostHandler)
				hosts.DELETE("/:id", hostService.DeleteHostHandler)
				hosts.POST("/:id/test-connection", hostService.TestConnectionHandler)

				// SSH相关
				hosts.POST("/:id/ssh/command", sshService.ExecuteSSHCommandHandler)
				hosts.POST("/:id/ssh/upload", sshService.UploadFileHandler)
				hosts.GET("/:id/ssh/download", sshService.DownloadFileHandler)
				hosts.POST("/:id/ssh/restart", sshService.RestartHostHandler)
				hosts.POST("/:id/ssh/shutdown", sshService.ShutdownHostHandler)
				// 交互式SSH终端WebSocket
				hosts.GET("/:id/ssh/ws", sshService.SSHWebSocketHandler)
				// 获取服务器资源统计信息
				hosts.GET("/:id/stats", sshService.GetServerStatsHandler)
			}
		}

		// 只有非nil的服务才注册路由
		if clusterService != nil {
			clusterService.RegisterRoutes(api)
		}

		if kvmService != nil {
			// KVM相关
			kvmService.RegisterRoutes(api)
		}

		if slsConfigService != nil {
			// SLS配置相关
			slsConfigService.RegisterRoutes(api)
		}

		if systemSettingService != nil {
			// 系统设置相关
			systemSettingService.RegisterRoutes(api)
		}

		if emailService != nil {
			// 邮件相关
			emailService.RegisterRoutes(api)
		}

		if emailSettingService != nil {
			// 邮箱配置相关
			emailSettingService.RegisterRoutes(api)
		}

		if dashboardService != nil {
			// 仪表盘相关
			dashboard := api.Group("/dashboard")
			{
				dashboard.GET("/stats", dashboardService.GetDashboardStats)
				dashboard.GET("/resources", dashboardService.GetResourceStats)
				dashboard.GET("/system-info", dashboardService.GetSystemInfo)
				dashboard.GET("/health-status", dashboardService.GetHealthStatus)
			}
		}

		if logService != nil {
			// 日志相关
			logService.RegisterRoutes(api)
		}

		if aiService != nil {
			// AI相关
			aiService.RegisterRoutes(api)
		}

		if md5Service != nil {
			// MD5相关
			md5Histories := api.Group("/md5-histories")
			{
				md5Histories.POST("", md5Service.CreateMD5HistoryHandler)
				md5Histories.GET("", md5Service.GetMD5HistoriesHandler)
				md5Histories.DELETE("/:id", md5Service.DeleteMD5HistoryHandler)
				md5Histories.DELETE("", md5Service.ClearMD5HistoriesHandler)
			}
		}

		if imageCompressService != nil {
			// 图片压缩相关
			imageCompressHistories := api.Group("/image-compress-histories")
			{
				imageCompressHistories.POST("", imageCompressService.CreateImageCompressHistoryHandler)
				imageCompressHistories.GET("", imageCompressService.GetImageCompressHistoriesHandler)
				imageCompressHistories.DELETE("/:id", imageCompressService.DeleteImageCompressHistoryHandler)
				imageCompressHistories.DELETE("", imageCompressService.ClearImageCompressHistoriesHandler)
			}
		}

		if telnetService != nil {
			// Telnet相关
			telnetService.RegisterRoutes(api)
		}

		if pingService != nil {
			// Ping相关
			pingService.RegisterRoutes(api)
		}

		if traceRouteService != nil {
			// 路由跟踪相关
			traceRouteService.RegisterRoutes(api)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return r
}
