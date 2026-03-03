package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"deantech/models"
)

// SLSConfigService SLS配置服务结构体
type SLSConfigService struct {
	db *gorm.DB
}

// NewSLSConfigService 创建SLS配置服务实例
func NewSLSConfigService(db *gorm.DB) *SLSConfigService {
	return &SLSConfigService{db: db}
}

// RegisterRoutes 注册SLS配置相关路由
func (s *SLSConfigService) RegisterRoutes(router *gin.RouterGroup) {
	slsGroup := router.Group("/sls")
	{
		slsGroup.GET("/configs", s.ListSLSConfigsHandler)
		slsGroup.GET("/configs/:id", s.GetSLSConfigByIDHandler)
		slsGroup.POST("/configs", s.CreateSLSConfigHandler)
		slsGroup.PUT("/configs/:id", s.UpdateSLSConfigHandler)
		slsGroup.DELETE("/configs/:id", s.DeleteSLSConfigHandler)
		slsGroup.PUT("/configs/:id/active", s.SetActiveSLSConfigHandler)
	}
}

// ListSLSConfigsHandler 获取SLS配置列表
func (s *SLSConfigService) ListSLSConfigsHandler(c *gin.Context) {
	var slsConfigs []models.SLSConfig
	result := s.db.Find(&slsConfigs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, slsConfigs)
}

// GetSLSConfigByIDHandler 根据ID获取SLS配置
func (s *SLSConfigService) GetSLSConfigByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var slsConfig models.SLSConfig
	result := s.db.First(&slsConfig, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "SLS配置不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, slsConfig)
}

// CreateSLSConfigHandler 创建SLS配置
func (s *SLSConfigService) CreateSLSConfigHandler(c *gin.Context) {
	var slsConfig models.SLSConfig
	if err := c.ShouldBindJSON(&slsConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果设置为活跃状态，将其他配置设置为非活跃
	if slsConfig.IsActive {
		s.db.Model(&models.SLSConfig{}).Update("is_active", false)
	}

	result := s.db.Create(&slsConfig)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SLS配置创建成功",
		"data":    slsConfig,
	})
}

// UpdateSLSConfigHandler 更新SLS配置
func (s *SLSConfigService) UpdateSLSConfigHandler(c *gin.Context) {
	id := c.Param("id")
	var slsConfig models.SLSConfig
	if err := c.ShouldBindJSON(&slsConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查配置是否存在
	var existingConfig models.SLSConfig
	result := s.db.First(&existingConfig, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "SLS配置不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 如果设置为活跃状态，将其他配置设置为非活跃
	if slsConfig.IsActive {
		s.db.Model(&models.SLSConfig{}).Update("is_active", false)
	}

	// 更新配置
	slsConfig.ID = existingConfig.ID
	result = s.db.Save(&slsConfig)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SLS配置更新成功",
		"data":    slsConfig,
	})
}

// DeleteSLSConfigHandler 删除SLS配置
func (s *SLSConfigService) DeleteSLSConfigHandler(c *gin.Context) {
	id := c.Param("id")

	// 检查配置是否存在
	var slsConfig models.SLSConfig
	result := s.db.First(&slsConfig, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "SLS配置不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 删除配置
	result = s.db.Delete(&slsConfig)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SLS配置删除成功"})
}

// SetActiveSLSConfigHandler 设置活跃SLS配置
func (s *SLSConfigService) SetActiveSLSConfigHandler(c *gin.Context) {
	id := c.Param("id")

	// 检查配置是否存在
	var slsConfig models.SLSConfig
	result := s.db.First(&slsConfig, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "SLS配置不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 将所有配置设置为非活跃，然后将当前配置设置为活跃
	s.db.Model(&models.SLSConfig{}).Update("is_active", false)
	s.db.Model(&slsConfig).Update("is_active", true)

	c.JSON(http.StatusOK, gin.H{"message": "活跃配置设置成功"})
}