package services

import (
	"deantech/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MD5Service MD5服务结构体
type MD5Service struct {
	db *gorm.DB
}

// NewMD5Service 创建MD5服务实例
func NewMD5Service(db *gorm.DB) *MD5Service {
	return &MD5Service{db: db}
}

// CreateMD5History 创建MD5加密历史记录
func (s *MD5Service) CreateMD5History(history *models.MD5History) error {
	return s.db.Create(history).Error
}

// GetMD5Histories 获取用户的MD5加密历史记录
func (s *MD5Service) GetMD5Histories(userID uint, limit, offset int) ([]models.MD5History, int64, error) {
	var histories []models.MD5History
	var total int64

	// 查询总记录数
	if err := s.db.Model(&models.MD5History{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询历史记录
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// DeleteMD5History 删除MD5加密历史记录（硬删除）
func (s *MD5Service) DeleteMD5History(userID, historyID uint) error {
	return s.db.Unscoped().Where("user_id = ? AND id = ?", userID, historyID).Delete(&models.MD5History{}).Error
}

// ClearMD5Histories 清空用户的MD5加密历史记录（硬删除）
func (s *MD5Service) ClearMD5Histories(userID uint) error {
	return s.db.Unscoped().Where("user_id = ?", userID).Delete(&models.MD5History{}).Error
}

// CreateMD5HistoryHandler 创建MD5加密历史记录的API处理方法
func (s *MD5Service) CreateMD5HistoryHandler(c *gin.Context) {
	var history models.MD5History

	// 绑定请求数据
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 创建历史记录
	if err := s.CreateMD5History(&history); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "MD5加密历史记录创建成功", "data": history})
}

// GetMD5HistoriesHandler 获取用户的MD5加密历史记录的API处理方法
func (s *MD5Service) GetMD5HistoriesHandler(c *gin.Context) {
	// 从请求头获取用户名
	username := c.GetHeader("X-Username")
	if username == "" {
		// 如果没有用户名，返回空的历史记录
		c.JSON(200, gin.H{
			"message": "获取MD5加密历史记录成功",
			"data": gin.H{
				"histories": []models.MD5History{},
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
			"message": "获取MD5加密历史记录成功",
			"data": gin.H{
				"histories": []models.MD5History{},
				"total":    0,
			},
		})
		return
	}

	// 获取分页参数
	limit := 20
	offset := 0

	// 查询历史记录
	histories, total, err := s.GetMD5Histories(user.ID, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "获取MD5加密历史记录成功",
		"data": gin.H{
			"histories": histories,
			"total":    total,
		},
	})
}

// DeleteMD5HistoryHandler 删除MD5加密历史记录的API处理方法
func (s *MD5Service) DeleteMD5HistoryHandler(c *gin.Context) {
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
	if err := s.DeleteMD5History(user.ID, id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "MD5加密历史记录删除成功"})
}

// ClearMD5HistoriesHandler 清空用户的MD5加密历史记录的API处理方法
func (s *MD5Service) ClearMD5HistoriesHandler(c *gin.Context) {
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
	if err := s.ClearMD5Histories(user.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "MD5加密历史记录清空成功"})
}