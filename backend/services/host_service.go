package services

import (
	"net/http"
	"deantech/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HostService 主机服务
type HostService struct {
	db *gorm.DB
}

// NewHostService 创建主机服务实例
func NewHostService(db *gorm.DB) *HostService {
	return &HostService{db: db}
}

// HostListItem 主机列表项，不包含敏感信息
type HostListItem struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	IP        string    `json:"ip"`
	Remark    string    `json:"remark"`
	IsActive  bool      `json:"isActive"`
	Status    string    `json:"status"`
}

// ListHostsHandler 获取主机列表
func (s *HostService) ListHostsHandler(c *gin.Context) {
	var hosts []models.Host
	result := s.db.Find(&hosts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get hosts",
		})
		return
	}

	// 测试每个主机的连接状态
	for i := range hosts {
		// 直接通过SSH连接测试，只要能建立连接就说明连接成功
		client, err := NewSSHClient(hosts[i])
		if err != nil {
			hosts[i].Status = "error"
		} else {
			hosts[i].Status = "running"
			// 关闭SSH连接
			client.Close()
		}
	}

	// 转换为不含敏感信息的列表项
	var hostListItems []HostListItem
	for _, host := range hosts {
		hostListItems = append(hostListItems, HostListItem{
			ID:        host.ID,
			CreatedAt: host.CreatedAt,
			UpdatedAt: host.UpdatedAt,
			Name:      host.Name,
			IP:        host.IP,
			Remark:    host.Remark,
			IsActive:  host.IsActive,
			Status:    host.Status,
		})
	}

	c.JSON(http.StatusOK, hostListItems)
}

// CreateHostHandler 创建主机
func (s *HostService) CreateHostHandler(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 创建主机
	result := s.db.Create(&host)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create host",
		})
		return
	}

	// 测试连接并更新状态
	client, err := NewSSHClient(host)
	if err != nil {
		host.Status = "error"
	} else {
		host.Status = "running"
		client.Close()
	}

	// 更新主机状态
	s.db.Model(&host).Update("status", host.Status)

	// 返回不含敏感信息的主机信息
	hostItem := HostListItem{
		ID:        host.ID,
		CreatedAt: host.CreatedAt,
		UpdatedAt: host.UpdatedAt,
		Name:      host.Name,
		IP:        host.IP,
		Remark:    host.Remark,
		IsActive:  host.IsActive,
		Status:    host.Status,
	}

	c.JSON(http.StatusCreated, hostItem)
}

// GetHostHandler 获取主机详情
func (s *HostService) GetHostHandler(c *gin.Context) {
	id := c.Param("id")
	var host models.Host
	result := s.db.First(&host, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	// 返回完整的主机信息，包括敏感字段，因为编辑、SSH等功能需要
	c.JSON(http.StatusOK, host)
}

// UpdateHostHandler 更新主机
func (s *HostService) UpdateHostHandler(c *gin.Context) {
	id := c.Param("id")
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 更新主机信息
	result := s.db.Model(&models.Host{}).Where("id = ?", id).Updates(host)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update host",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	// 重新获取主机信息
	s.db.First(&host, id)

	// 测试连接并更新状态
	client, err := NewSSHClient(host)
	if err != nil {
		host.Status = "error"
	} else {
		host.Status = "running"
		client.Close()
	}

	// 更新主机状态
	s.db.Model(&host).Update("status", host.Status)

	c.JSON(http.StatusOK, gin.H{
		"message": "Host updated successfully",
		"status":  host.Status,
	})
}

// DeleteHostHandler 删除主机（彻底删除）
func (s *HostService) DeleteHostHandler(c *gin.Context) {
	id := c.Param("id")
	// 使用Unscoped()方法实现彻底删除，而不是软删除
	result := s.db.Unscoped().Delete(&models.Host{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete host",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Host deleted successfully",
	})
}

// TestConnectionHandler 测试主机连接
func (s *HostService) TestConnectionHandler(c *gin.Context) {
	id := c.Param("id")

	// 获取主机信息
	var host models.Host
	result := s.db.First(&host, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	// 创建SSH客户端，测试连接
	client, err := NewSSHClient(host)
	if err != nil {
		// 连接失败，更新状态为error
		s.db.Model(&host).Update("status", "error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to host",
		})
		return
	}
	defer client.Close()

	// 连接成功，更新状态为running
	s.db.Model(&host).Update("status", "running")

	c.JSON(http.StatusOK, gin.H{
		"message": "Host connection test successful",
		"status":  "running",
	})
}
