package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"deantech/models"
)

// EmailSettingService 邮箱配置服务结构体
type EmailSettingService struct {
	db *gorm.DB
}

// NewEmailSettingService 创建邮箱配置服务实例
func NewEmailSettingService(db *gorm.DB) *EmailSettingService {
	return &EmailSettingService{db: db}
}

// RegisterRoutes 注册邮箱配置相关路由
func (s *EmailSettingService) RegisterRoutes(router *gin.RouterGroup) {
	emailGroup := router.Group("/email-settings")
	{
		emailGroup.GET("", s.GetEmailSettingHandler)      // 获取邮箱配置
		emailGroup.PUT("", s.UpdateEmailSettingHandler)   // 更新邮箱配置
	}
}

// GetEmailSettingHandler 获取邮箱配置
func (s *EmailSettingService) GetEmailSettingHandler(c *gin.Context) {
	var emailSetting models.EmailSetting
	result := s.db.First(&emailSetting)
	if result.Error != nil {
		// 如果不存在，返回空的默认配置
		c.JSON(http.StatusOK, models.EmailSetting{
			SMTPPort:       465,
			SMTPEncryption: "ssl",
		})
		return
	}
	c.JSON(http.StatusOK, emailSetting)
}

// UpdateEmailSettingHandler 更新邮箱配置
func (s *EmailSettingService) UpdateEmailSettingHandler(c *gin.Context) {
	var emailSetting models.EmailSetting
	if err := c.ShouldBindJSON(&emailSetting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 查找现有配置
	var existingSetting models.EmailSetting
	result := s.db.First(&existingSetting)
	if result.Error != nil {
		// 如果不存在，创建新配置
		result = s.db.Create(&emailSetting)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create email setting: " + result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, emailSetting)
		return
	}

	// 更新现有配置
	existingSetting.SMTPServer = emailSetting.SMTPServer
	existingSetting.SMTPPort = emailSetting.SMTPPort
	existingSetting.FromEmail = emailSetting.FromEmail
	existingSetting.SMTPUsername = emailSetting.SMTPUsername
	existingSetting.SMTPPassword = emailSetting.SMTPPassword
	existingSetting.SMTPEncryption = emailSetting.SMTPEncryption

	result = s.db.Save(&existingSetting)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update email setting: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, existingSetting)
}
