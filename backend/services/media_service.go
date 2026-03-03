package services

import (
	"net/http"
	"deantech/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MediaService struct {
	db *gorm.DB
}

func NewMediaService(db *gorm.DB) *MediaService {
	return &MediaService{
		db: db,
	}
}

// ListMedia 列出所有告警介质
func (s *MediaService) ListMedia(c *gin.Context) {
	var mediaList []models.AlertMedia
	if err := s.db.Preload("Template").Find(&mediaList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mediaList)
}

// CreateMedia 创建告警介质
func (s *MediaService) CreateMedia(c *gin.Context) {
	var media models.AlertMedia
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Create(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, media)
}

// GetMedia 获取单个告警介质
func (s *MediaService) GetMedia(c *gin.Context) {
	id := c.Param("id")
	var media models.AlertMedia
	if err := s.db.Preload("Template").First(&media, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}
	c.JSON(http.StatusOK, media)
}

// UpdateMedia 更新告警介质
func (s *MediaService) UpdateMedia(c *gin.Context) {
	id := c.Param("id")
	var media models.AlertMedia
	if err := s.db.First(&media, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Save(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, media)
}

// DeleteMedia 删除告警介质
func (s *MediaService) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	if err := s.db.Delete(&models.AlertMedia{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

// ListTemplates 列出所有告警模板
func (s *MediaService) ListTemplates(c *gin.Context) {
	var templates []models.AlertTemplate
	if err := s.db.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// CreateTemplate 创建告警模板
func (s *MediaService) CreateTemplate(c *gin.Context) {
	var template models.AlertTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果设置为默认模板，将其他模板的默认状态取消
	if template.IsDefault {
		s.db.Model(&models.AlertTemplate{}).Update("is_default", false)
	}

	if err := s.db.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetTemplate 获取单个告警模板
func (s *MediaService) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.AlertTemplate
	if err := s.db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// UpdateTemplate 更新告警模板
func (s *MediaService) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.AlertTemplate
	if err := s.db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果设置为默认模板，将其他模板的默认状态取消
	if template.IsDefault {
		s.db.Model(&models.AlertTemplate{}).Update("is_default", false)
	}

	if err := s.db.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate 删除告警模板
func (s *MediaService) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := s.db.Delete(&models.AlertTemplate{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}