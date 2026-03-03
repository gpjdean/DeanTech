package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"deantech/models"
)

// AiService AI服务结构体
type AiService struct {
	db *gorm.DB
}

// NewAiService 创建AI服务实例
func NewAiService(db *gorm.DB) *AiService {
	return &AiService{db: db}
}

// RegisterRoutes 注册AI相关路由
func (s *AiService) RegisterRoutes(router *gin.RouterGroup) {
	aiGroup := router.Group("/ai")
	{
		aiGroup.POST("/chat", s.ChatHandler)           // AI聊天
		aiGroup.POST("/test-connection", s.TestConnectionHandler) // 测试AI连接
	}
}

// ChatRequest 聊天请求结构体
type ChatRequest struct {
	Model     string `json:"model" binding:"required"`     // AI模型
	Message   string `json:"message" binding:"required"`   // 消息内容
	History   []map[string]string `json:"history"`         // 聊天历史
}

// ChatResponse 聊天响应结构体
type ChatResponse struct {
	Response string `json:"response"` // AI回复内容
}

// ChatHandler AI聊天处理函数
func (s *AiService) ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 从数据库获取AI配置
	var aiConfig models.SystemSetting
	result := s.db.Where("`key` = ?", "aiConfig").First(&aiConfig)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get AI config: " + result.Error.Error(),
		})
		return
	}

	// 这里需要实现具体的AI模型调用逻辑
	// 暂时先返回模拟数据
	response := ChatResponse{
		Response: "这是" + req.Model + "模型对您的回复：" + req.Message,
	}

	c.JSON(http.StatusOK, response)
}

// TestConnectionRequest 测试连接请求结构体
type TestConnectionRequest struct {
	Model   string `json:"model" binding:"required"`   // AI模型
	ApiKey  string `json:"apiKey" binding:"required"`  // API Key
}

// TestConnectionResponse 测试连接响应结构体
type TestConnectionResponse struct {
	Success bool   `json:"success"` // 连接是否成功
	Message string `json:"message"` // 响应消息
}

// TestConnectionHandler 测试AI连接处理函数
func (s *AiService) TestConnectionHandler(c *gin.Context) {
	var req TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 这里需要实现具体的AI连接测试逻辑
	// 暂时先返回模拟数据
	response := TestConnectionResponse{
		Success: true,
		Message: req.Model + "连接成功！",
	}

	c.JSON(http.StatusOK, response)
}
