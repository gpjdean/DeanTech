package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"deantech/models"
)

// LogService 日志服务结构体
type LogService struct {
	db *gorm.DB
}

// NewLogService 创建日志服务实例
func NewLogService(db *gorm.DB) *LogService {
	return &LogService{db: db}
}

// RegisterRoutes 注册日志相关路由
func (s *LogService) RegisterRoutes(router *gin.RouterGroup) {
	logGroup := router.Group("/logs")
	{
		// 操作日志
		logGroup.GET("/operation", s.ListOperationLogsHandler)
		logGroup.GET("/operation/:id", s.GetOperationLogHandler)
		logGroup.DELETE("/operation/:id", s.DeleteOperationLogHandler)
		logGroup.DELETE("/operation", s.BatchDeleteOperationLogsHandler)

		// 登录日志
		logGroup.GET("/login", s.ListLoginLogsHandler)
		logGroup.GET("/login/:id", s.GetLoginLogHandler)
		logGroup.DELETE("/login/:id", s.DeleteLoginLogHandler)
		logGroup.DELETE("/login", s.BatchDeleteLoginLogsHandler)
	}
}

// ListOperationLogsHandler 获取操作日志列表
func (s *LogService) ListOperationLogsHandler(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 获取筛选参数
	status := c.Query("status")
	search := c.Query("search")

	// 查询操作日志列表
	var operationLogs []models.OperationLog
	var total int64

	// 构建查询，包含软删除的记录
	query := s.db.Unscoped().Model(&models.OperationLog{})

	// 添加筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("username LIKE ? OR menu LIKE ? OR operation LIKE ? OR api LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 统计总数
	query.Count(&total)

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&operationLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list operation logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": operationLogs,
		"page":  page,
		"size":  pageSize,
	})
}

// GetOperationLogHandler 获取单个操作日志
func (s *LogService) GetOperationLogHandler(c *gin.Context) {
	id := c.Param("id")

	var operationLog models.OperationLog
	if err := s.db.First(&operationLog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operation log not found"})
		return
	}

	c.JSON(http.StatusOK, operationLog)
}

// DeleteOperationLogHandler 删除单个操作日志
func (s *LogService) DeleteOperationLogHandler(c *gin.Context) {
	id := c.Param("id")

	// 使用硬删除，确保记录真正被删除
	if err := s.db.Unscoped().Delete(&models.OperationLog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete operation log: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Operation log deleted successfully"})
}

// BatchDeleteOperationLogsHandler 批量删除操作日志
func (s *LogService) BatchDeleteOperationLogsHandler(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 使用硬删除，确保记录真正被删除
	if err := s.db.Unscoped().Delete(&models.OperationLog{}, request.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to batch delete operation logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Operation logs deleted successfully"})
}

// ListLoginLogsHandler 获取登录日志列表
func (s *LogService) ListLoginLogsHandler(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查询登录日志列表
	var loginLogs []models.LoginLog
	var total int64

	// 构建查询，包含软删除的记录
	query := s.db.Unscoped().Model(&models.LoginLog{})

	// 统计总数
	query.Count(&total)

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&loginLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list login logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": loginLogs,
		"page":  page,
		"size":  pageSize,
	})
}

// GetLoginLogHandler 获取单个登录日志
func (s *LogService) GetLoginLogHandler(c *gin.Context) {
	id := c.Param("id")

	var loginLog models.LoginLog
	if err := s.db.First(&loginLog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Login log not found"})
		return
	}

	c.JSON(http.StatusOK, loginLog)
}

// DeleteLoginLogHandler 删除单个登录日志
func (s *LogService) DeleteLoginLogHandler(c *gin.Context) {
	id := c.Param("id")

	// 使用硬删除，确保记录真正被删除
	if err := s.db.Unscoped().Delete(&models.LoginLog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete login log: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login log deleted successfully"})
}

// BatchDeleteLoginLogsHandler 批量删除登录日志
func (s *LogService) BatchDeleteLoginLogsHandler(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 使用硬删除，确保记录真正被删除
	if err := s.db.Unscoped().Delete(&models.LoginLog{}, request.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to batch delete login logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login logs deleted successfully"})
}
