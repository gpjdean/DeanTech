package services

import (
	"net/http"
	"deantech/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AlertService struct {
	db                 *gorm.DB
	prometheusService *PrometheusService
}

func NewAlertService(db *gorm.DB, prometheusService *PrometheusService) *AlertService {
	return &AlertService{
		db:                 db,
		prometheusService: prometheusService,
	}
}

// ListAlerts 列出所有告警
func (s *AlertService) ListAlerts(c *gin.Context) {
	var alerts []models.Alert
	if err := s.db.Preload("AlertMedia").Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

// GetAlert 获取单个告警
func (s *AlertService) GetAlert(c *gin.Context) {
	id := c.Param("id")
	var alert models.Alert
	if err := s.db.Preload("AlertMedia").First(&alert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}
	c.JSON(http.StatusOK, alert)
}

// ResolveAlert 解决告警
func (s *AlertService) ResolveAlert(c *gin.Context) {
	id := c.Param("id")
	var alert models.Alert
	if err := s.db.First(&alert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	alert.Status = "resolved"
	alert.ResolvedAt = time.Now()
	// 这里可以添加解决人信息，从JWT token中获取
	// alert.ResolvedBy = c.GetHeader("X-User")

	if err := s.db.Save(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// DeleteAlert 删除告警
func (s *AlertService) DeleteAlert(c *gin.Context) {
	id := c.Param("id")
	if err := s.db.Delete(&models.Alert{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})
}

// ListSilences 列出所有静默规则
func (s *AlertService) ListSilences(c *gin.Context) {
	var silences []models.Silence
	if err := s.db.Find(&silences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, silences)
}

// CreateSilence 创建静默规则
func (s *AlertService) CreateSilence(c *gin.Context) {
	var silence models.Silence
	if err := c.ShouldBindJSON(&silence); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	silence.Status = "active"
	if err := s.db.Create(&silence).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, silence)
}

// GetSilence 获取单个静默规则
func (s *AlertService) GetSilence(c *gin.Context) {
	id := c.Param("id")
	var silence models.Silence
	if err := s.db.First(&silence, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Silence not found"})
		return
	}
	c.JSON(http.StatusOK, silence)
}

// DeleteSilence 删除静默规则
func (s *AlertService) DeleteSilence(c *gin.Context) {
	id := c.Param("id")
	if err := s.db.Delete(&models.Silence{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Silence deleted successfully"})
}

// ListInhibitions 列出所有抑制规则
func (s *AlertService) ListInhibitions(c *gin.Context) {
	var inhibitions []models.InhibitionRule
	if err := s.db.Find(&inhibitions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inhibitions)
}

// CreateInhibition 创建抑制规则
func (s *AlertService) CreateInhibition(c *gin.Context) {
	var inhibition models.InhibitionRule
	if err := c.ShouldBindJSON(&inhibition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Create(&inhibition).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inhibition)
}

// GetInhibition 获取单个抑制规则
func (s *AlertService) GetInhibition(c *gin.Context) {
	id := c.Param("id")
	var inhibition models.InhibitionRule
	if err := s.db.First(&inhibition, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inhibition rule not found"})
		return
	}
	c.JSON(http.StatusOK, inhibition)
}

// UpdateInhibition 更新抑制规则
func (s *AlertService) UpdateInhibition(c *gin.Context) {
	id := c.Param("id")
	var inhibition models.InhibitionRule
	if err := s.db.First(&inhibition, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inhibition rule not found"})
		return
	}

	if err := c.ShouldBindJSON(&inhibition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Save(&inhibition).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inhibition)
}

// DeleteInhibition 删除抑制规则
func (s *AlertService) DeleteInhibition(c *gin.Context) {
	id := c.Param("id")
	if err := s.db.Delete(&models.InhibitionRule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inhibition rule deleted successfully"})
}