package middleware

import (
	"deantech/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperationLogger 操作日志中间件
func OperationLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()
		
		// 执行请求处理
		c.Next()
		
		// 记录请求结束时间
		endTime := time.Now()
		
		// 计算请求耗时
		timeCost := endTime.Sub(startTime).Milliseconds()
		
		// 跳过静态文件请求和健康检查
		if strings.HasPrefix(c.Request.URL.Path, "/static") || c.Request.URL.Path == "/health" {
			return
		}
		
		// 跳过日志相关的API，避免递归记录
		if strings.HasPrefix(c.Request.URL.Path, "/api/logs/") {
			return
		}
		
		// 获取用户名
		username := c.GetHeader("X-Username")
		if username == "" {
			username = "unknown"
		}
		
		// 从上下文中获取用户ID
		userID, exists := c.Get("userID")
		var userIDUint uint
		if !exists {
			userIDUint = 0
		} else {
			// 处理类型转换，将int转换为uint
			if uidInt, ok := userID.(int); ok {
				userIDUint = uint(uidInt)
			} else if uidUint, ok := userID.(uint); ok {
				userIDUint = uidUint
			} else {
				userIDUint = 0
			}
		}
		
		// 获取请求状态
		status := "success"
		if c.Writer.Status() >= 400 {
			status = "failed"
		}
		
		// 构建API路径（去掉/api前缀）
		apiPath := c.Request.URL.Path
		if strings.HasPrefix(apiPath, "/api") {
			apiPath = strings.TrimPrefix(apiPath, "/api")
		}
		
		// 根据API路径推断菜单
		menu := getMenuFromAPIPath(apiPath)
		
		// 创建操作日志
		operationLog := models.OperationLog{
			UserID:      userIDUint,
			Username:    username,
			Menu:        menu,
			Operation:   c.Request.Method + " " + apiPath,
			API:         apiPath,
			Method:      c.Request.Method,
			Params:      getRequestParams(c),
			Result:      getResponseStatus(c),
			Status:      status,
			IP:          c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			TimeCost:    timeCost,
		}
		
		// 保存到数据库
		db.Create(&operationLog)
	}
}

// getMenuFromAPIPath 根据API路径推断菜单
func getMenuFromAPIPath(apiPath string) string {
	// 分割API路径
	parts := strings.Split(strings.Trim(apiPath, "/"), "/")
	if len(parts) == 0 {
		return "首页"
	}
	
	// 菜单映射
	menuMap := map[string]string{
		"alerts":        "告警管理",
		"media":         "告警介质",
		"templates":     "告警模板",
		"silences":      "静默规则",
		"inhibitions":   "抑制规则",
		"prometheus":    "Prometheus管理",
		"hosts":         "主机管理",
		"kvm":           "KVM管理",
		"clusters":      "集群管理",
		"users":         "用户管理",
		"settings":      "系统设置",
		"sls":           "SLS配置",
		"dashboard":     "平台概览",
		"email":         "邮件配置",
		"email-settings": "邮件配置",
		"logs":          "日志管理",
		"nodes":         "节点管理",
		"pods":          "Pod管理",
		"deployments":   "Deployment管理",
		"statefulsets":  "StatefulSet管理",
		"daemonsets":    "DaemonSet管理",
		"jobs":          "Job管理",
		"cronjobs":      "CronJob管理",
		"services":      "Service管理",
		"configmaps":    "ConfigMap管理",
		"secrets":       "Secret管理",
		"ingresses":     "Ingress管理",
		"namespaces":    "命名空间管理",
		"replicasets":   "ReplicaSet管理",
		"pvcs":          "PVC管理",
		"pvs":           "PV管理",
		"storageclasses": "存储类管理",
		"events":        "事件管理",
		"ssh":           "SSH管理",
	}
	
	// 获取第一级路径作为菜单
	if menu, exists := menuMap[parts[0]]; exists {
		return menu
	}
	
	return "其他"
}

// getRequestParams 获取请求参数
func getRequestParams(c *gin.Context) string {
	// 简单处理，实际项目中可以根据需要格式化
	if c.Request.Method == "GET" {
		return c.Request.URL.RawQuery
	}
	
	// 对于POST、PUT等方法，可以从c.Request.Body中读取
	// 但需要注意，这里c.Request.Body已经被读取过了
	// 所以我们简单返回"[body]"表示有请求体
	return "[body]"
}

// getResponseStatus 获取响应状态
func getResponseStatus(c *gin.Context) string {
	return http.StatusText(c.Writer.Status())
}
