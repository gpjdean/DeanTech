package services

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// TelnetService 提供Telnet测试服务
type TelnetService struct{}

// NewTelnetService 创建一个新的TelnetService实例
func NewTelnetService() *TelnetService {
	return &TelnetService{}
}

// TelnetTestRequest 表示Telnet测试请求的参数
type TelnetTestRequest struct {
	Target  string `json:"target" binding:"required"`
	Port    int    `json:"port" binding:"required,min=1,max=65535"`
	Timeout int    `json:"timeout" binding:"required,min=100,max=5000"`
}

// TelnetTestResponse 表示Telnet测试的响应结果
type TelnetTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details"`
	Elapsed int64  `json:"elapsed"`
}

// TestTelnet 执行Telnet测试
func (s *TelnetService) TestTelnet(c *gin.Context) {
	var req TelnetTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 记录开始时间
	startTime := time.Now()

	// 执行Telnet连接测试
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(req.Target, strconv.Itoa(req.Port)), time.Duration(req.Timeout)*time.Millisecond)
	
	// 计算耗时
	elapsed := time.Since(startTime).Milliseconds()

	response := &TelnetTestResponse{
		Elapsed: elapsed,
	}

	if err != nil {
		// 连接失败
		response.Success = false
		response.Message = "连接失败: " + req.Target + ":" + strconv.Itoa(req.Port)
		response.Details = "正在连接 " + req.Target + ":" + strconv.Itoa(req.Port) + "...\n" +
			"连接到 " + req.Target + " 失败。\n" +
			"错误信息: " + err.Error() + "\n" +
			"耗时: " + strconv.FormatInt(elapsed, 10) + " 毫秒"
	} else {
		// 连接成功
		defer conn.Close()
		
		response.Success = true
		response.Message = "连接成功: " + req.Target + ":" + strconv.Itoa(req.Port)
		response.Details = "正在连接 " + req.Target + ":" + strconv.Itoa(req.Port) + "...\n" +
			"连接到 " + req.Target + " 成功。\n" +
			"Escape 字符是 '^]'\n" +
			"耗时: " + strconv.FormatInt(elapsed, 10) + " 毫秒"
	}

	c.JSON(200, response)
}

// RegisterRoutes 注册Telnet服务的路由
func (s *TelnetService) RegisterRoutes(api *gin.RouterGroup) {
	telnet := api.Group("/telnet")
	{
		telnet.POST("/test", s.TestTelnet)
	}
}
