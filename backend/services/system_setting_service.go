package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"deantech/models"
)

// SystemSettingService 系统设置服务结构体
type SystemSettingService struct {
	db *gorm.DB
}

// NewSystemSettingService 创建系统设置服务实例
func NewSystemSettingService(db *gorm.DB) *SystemSettingService {
	return &SystemSettingService{db: db}
}

// RegisterRoutes 注册系统设置相关路由
func (s *SystemSettingService) RegisterRoutes(router *gin.RouterGroup) {
	settingGroup := router.Group("/settings")
	{
		settingGroup.GET("", s.ListSettingsHandler)          // 列出所有系统设置
		settingGroup.GET("/:key", s.GetSettingHandler)       // 获取单个系统设置
		settingGroup.POST("", s.CreateSettingHandler)        // 创建系统设置
		settingGroup.PUT("/:key", s.UpdateSettingHandler)    // 更新系统设置
		settingGroup.DELETE("/:key", s.DeleteSettingHandler) // 删除系统设置
	}
}

// ListSettingsHandler 列出所有系统设置
func (s *SystemSettingService) ListSettingsHandler(c *gin.Context) {
	var settings []models.SystemSetting
	result := s.db.Find(&settings)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list settings: " + result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// GetSettingHandler 获取单个系统设置
func (s *SystemSettingService) GetSettingHandler(c *gin.Context) {
	key := c.Param("key")
	var setting models.SystemSetting
	result := s.db.Where("`key` = ?", key).First(&setting)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Setting not found",
		})
		return
	}
	c.JSON(http.StatusOK, setting)
}

// CreateSettingHandler 创建系统设置
func (s *SystemSettingService) CreateSettingHandler(c *gin.Context) {
	var setting models.SystemSetting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 检查key是否已存在
	var existingSetting models.SystemSetting
	result := s.db.Where("`key` = ?", setting.Key).First(&existingSetting)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Setting key already exists",
		})
		return
	}

	result = s.db.Create(&setting)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create setting: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, setting)
}

// UpdateSettingHandler 更新系统设置
func (s *SystemSettingService) UpdateSettingHandler(c *gin.Context) {
	key := c.Param("key")
	var updatedSetting models.SystemSetting
	if err := c.ShouldBindJSON(&updatedSetting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 查找现有设置
	var setting models.SystemSetting
	result := s.db.Where("`key` = ?", key).First(&setting)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Setting not found",
		})
		return
	}

	// 更新字段
	setting.Value = updatedSetting.Value
	setting.Description = updatedSetting.Description

	result = s.db.Save(&setting)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update setting: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// DeleteSettingHandler 删除系统设置
func (s *SystemSettingService) DeleteSettingHandler(c *gin.Context) {
	key := c.Param("key")

	result := s.db.Where("`key` = ?", key).Delete(&models.SystemSetting{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete setting: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Setting deleted successfully",
	})
}
