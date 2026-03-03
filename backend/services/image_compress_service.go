package services

import (
	"deantech/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImageCompressService 图片压缩服务结构体
type ImageCompressService struct {
	db *gorm.DB
}

// NewImageCompressService 创建图片压缩服务实例
func NewImageCompressService(db *gorm.DB) *ImageCompressService {
	return &ImageCompressService{db: db}
}

// CreateImageCompressHistory 创建图片压缩历史记录
func (s *ImageCompressService) CreateImageCompressHistory(history *models.ImageCompressHistory) error {
	return s.db.Create(history).Error
}

// GetImageCompressHistories 获取用户的图片压缩历史记录
func (s *ImageCompressService) GetImageCompressHistories(userID uint, limit, offset int) ([]models.ImageCompressHistory, int64, error) {
	var histories []models.ImageCompressHistory
	var total int64

	// 查询总记录数
	if err := s.db.Model(&models.ImageCompressHistory{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询历史记录
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// DeleteImageCompressHistory 删除图片压缩历史记录（硬删除）
func (s *ImageCompressService) DeleteImageCompressHistory(userID, historyID uint) error {
	return s.db.Unscoped().Where("user_id = ? AND id = ?", userID, historyID).Delete(&models.ImageCompressHistory{}).Error
}

// ClearImageCompressHistories 清空用户的图片压缩历史记录（硬删除）
func (s *ImageCompressService) ClearImageCompressHistories(userID uint) error {
	return s.db.Unscoped().Where("user_id = ?", userID).Delete(&models.ImageCompressHistory{}).Error
}

// CreateImageCompressHistoryHandler 创建图片压缩历史记录的API处理方法
func (s *ImageCompressService) CreateImageCompressHistoryHandler(c *gin.Context) {
	var history models.ImageCompressHistory

	// 绑定请求数据
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 创建历史记录
	if err := s.CreateImageCompressHistory(&history); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "图片压缩历史记录创建成功", "data": history})
}

// GetImageCompressHistoriesHandler 获取用户的图片压缩历史记录的API处理方法
func (s *ImageCompressService) GetImageCompressHistoriesHandler(c *gin.Context) {
	// 从请求头获取用户名
	username := c.GetHeader("X-Username")
	if username == "" {
		// 如果没有用户名，返回空的历史记录
		c.JSON(200, gin.H{
			"message": "获取图片压缩历史记录成功",
			"data": gin.H{
				"histories": []models.ImageCompressHistory{},
				"total":    0,
			},
		})
		return
	}

	// 根据用户名查询用户ID
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		// 如果查询不到用户，返回空的历史记录
		c.JSON(200, gin.H{
			"message": "获取图片压缩历史记录成功",
			"data": gin.H{
				"histories": []models.ImageCompressHistory{},
				"total":    0,
			},
		})
		return
	}

	// 获取分页参数
	limit := 20
	offset := 0

	// 查询历史记录
	histories, total, err := s.GetImageCompressHistories(user.ID, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "获取图片压缩历史记录成功",
		"data": gin.H{
			"histories": histories,
			"total":    total,
		},
	})
}

// DeleteImageCompressHistoryHandler 删除图片压缩历史记录的API处理方法
func (s *ImageCompressService) DeleteImageCompressHistoryHandler(c *gin.Context) {
	// 从请求头获取用户名
	username := c.GetHeader("X-Username")
	if username == "" {
		c.JSON(401, gin.H{"error": "未授权"})
		return
	}

	// 根据用户名查询用户ID
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "未授权"})
		return
	}

	// 获取历史记录ID
	historyID := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(historyID, "%d", &id); err != nil {
		c.JSON(400, gin.H{"error": "无效的历史记录ID"})
		return
	}

	// 删除历史记录
	if err := s.DeleteImageCompressHistory(user.ID, id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "图片压缩历史记录删除成功"})
}

// ClearImageCompressHistoriesHandler 清空用户的图片压缩历史记录的API处理方法
func (s *ImageCompressService) ClearImageCompressHistoriesHandler(c *gin.Context) {
	// 从请求头获取用户名
	username := c.GetHeader("X-Username")
	if username == "" {
		c.JSON(401, gin.H{"error": "未授权"})
		return
	}

	// 根据用户名查询用户ID
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "未授权"})
		return
	}

	// 清空历史记录
	if err := s.ClearImageCompressHistories(user.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "图片压缩历史记录清空成功"})
}