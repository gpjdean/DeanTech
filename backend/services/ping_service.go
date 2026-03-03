package services

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// PingService 提供Ping测试服务
type PingService struct{}

// NewPingService 创建一个新的PingService实例
func NewPingService() *PingService {
	return &PingService{}
}

// PingTestRequest 表示Ping测试请求的参数
type PingTestRequest struct {
	Target  string `json:"target" binding:"required"`
	Count   int    `json:"count" binding:"required,min=1,max=10"`
	Timeout int    `json:"timeout" binding:"required,min=100,max=5000"`
	Size    int    `json:"size" binding:"required,min=32,max=1500"`
}

// PingResult 表示单次Ping测试的结果
type PingResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	IP      string `json:"ip"`
	Bytes   int    `json:"bytes"`
	Time    string `json:"time"`
	TTL     int    `json:"ttl"`
}

// PingTestResponse 表示Ping测试的完整响应结果
type PingTestResponse struct {
	Results []PingResult `json:"results"`
	Stats   PingStats    `json:"stats"`
}

// PingStats 表示Ping测试的统计信息
type PingStats struct {
	Success    bool    `json:"success"`
	Sent       int     `json:"sent"`
	Received   int     `json:"received"`
	Lost       int     `json:"lost"`
	MinDelay   float64 `json:"minDelay"`
	AvgDelay   float64 `json:"avgDelay"`
	MaxDelay   float64 `json:"maxDelay"`
	PacketLoss float64 `json:"packetLoss"`
}

// TestPing 执行Ping测试
func (s *PingService) TestPing(c *gin.Context) {
	var req PingTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 执行Ping测试
	results, stats := s.executePing(req)

	// 返回结果
	c.JSON(200, PingTestResponse{
		Results: results,
		Stats:   stats,
	})
}

// executePing 执行实际的Ping命令
func (s *PingService) executePing(req PingTestRequest) ([]PingResult, PingStats) {
	// 初始化结果
	results := make([]PingResult, 0, req.Count)
	stats := PingStats{
		Sent:     req.Count,
		Received: 0,
		Lost:     0,
		MinDelay: 0,
		MaxDelay: 0,
		PacketLoss: 0,
	}

	// 初始化延迟数组用于计算统计信息
	delays := make([]float64, 0, req.Count)

	// 根据不同操作系统构建Ping命令
	cmdArgs := s.buildPingCommand(req)

	// 执行Ping命令，增加5秒的额外超时时间，确保命令有足够时间完成
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.Timeout*req.Count+5000)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	// 打印完整的Ping命令输出，方便调试
	outputStr := string(output)
	fmt.Printf("Ping command output: %s\n", outputStr)

	if err != nil {
		// 如果命令执行失败，仍然解析输出（可能部分成功）
		fmt.Printf("Ping command failed: %v\n", err)
	}

	// 解析Ping命令输出
	lines := strings.Split(outputStr, "\n")

	// 首先解析统计信息，确保获取到延迟数据
	// 处理统计行，macOS格式：4 packets transmitted, 2 packets received, 50.0% packet loss, 2 packets out of wait time
	// Linux格式：4 packets transmitted, 2 received, 50% packet loss, time 3001ms
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// 解析延迟统计行，macOS格式：round-trip min/avg/max/stddev = 192.185/263.544/348.514/59.564 ms
		if strings.Contains(line, "round-trip") && strings.Contains(line, "min/avg/max/stddev") {
			// 解析延迟统计，使用更灵活的正则表达式
			timeMatch := regexp.MustCompile(`round-trip min/avg/max/stddev\s*=\s*([\d.]+)/([\d.]+)/([\d.]+)/([\d.]+)\s*ms`).FindStringSubmatch(line)
			if len(timeMatch) > 4 {
				if min, err := strconv.ParseFloat(timeMatch[1], 64); err == nil {
					stats.MinDelay = min
				}
				if avg, err := strconv.ParseFloat(timeMatch[2], 64); err == nil {
					stats.AvgDelay = avg
				}
				if max, err := strconv.ParseFloat(timeMatch[3], 64); err == nil {
					stats.MaxDelay = max
				}
				// 如果有延迟统计，说明至少有部分成功
				stats.Success = true
			}
			continue
		}
		
		// 解析传输统计行
		if strings.Contains(line, "packets transmitted") {
			// 解析统计信息
			// 解析发送的数据包数
			if sentMatch := regexp.MustCompile(`(\d+) packets transmitted`).FindStringSubmatch(line); len(sentMatch) > 1 {
				if s, err := strconv.Atoi(sentMatch[1]); err == nil {
					stats.Sent = s
				}
			}
			
			// 解析收到的数据包数
			if recvMatch := regexp.MustCompile(`(\d+) packets? received`).FindStringSubmatch(line); len(recvMatch) > 1 {
				if r, err := strconv.Atoi(recvMatch[1]); err == nil {
					stats.Received = r
				}
			}
			
			// 解析丢包率
			if lossMatch := regexp.MustCompile(`(\d+(\.\d+)?)% packet loss`).FindStringSubmatch(line); len(lossMatch) > 1 {
				if pl, err := strconv.ParseFloat(lossMatch[1], 64); err == nil {
					stats.PacketLoss = pl
				}
			}
			
			// 不要使用break，因为延迟统计行可能在传输统计行之后
			continue
		}
	}

	// 处理每一行输出，解析单次Ping结果
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析单次Ping结果
		result := s.parsePingResult(line, req.Size)
		if result != nil {
			results = append(results, *result)
			
			// 更新统计信息
			if result.Success {
				timeVal, err := strconv.ParseFloat(result.Time, 64)
				if err == nil {
					delays = append(delays, timeVal)
				}
			}
		}
	}

	// 计算丢失的数据包数
	stats.Lost = stats.Sent - stats.Received

	// 处理macOS Ping命令的特殊情况：当使用-W选项时，可能不会显示每个数据包的详细信息
	// 只显示统计信息，此时需要为每个成功的数据包生成模拟结果
	if len(results) == 0 && stats.Received > 0 {
		// 从输出中提取目标IP
		var targetIP string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "PING") {
				ipMatch := regexp.MustCompile(`\(([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\)`).FindStringSubmatch(line)
				if len(ipMatch) > 1 {
					targetIP = ipMatch[1]
					break
				}
			}
		}
		
		// 为每个成功的数据包生成模拟结果
		for i := 0; i < stats.Received; i++ {
			results = append(results, PingResult{
				Success: true,
				Message: "请求成功",
				IP:      targetIP,
				Bytes:   req.Size,
				Time:    fmt.Sprintf("%.2f", stats.AvgDelay),
				TTL:     128, // 默认TTL值
			})
		}
		
		// 为丢失的数据包生成模拟结果
		for i := 0; i < stats.Lost; i++ {
			results = append(results, PingResult{
				Success: false,
				Message: "请求超时",
				IP:      "",
				Bytes:   0,
				Time:    "0",
				TTL:     0,
			})
		}
	}

	// 计算平均延迟（如果没有从延迟统计行获取到的话）
	if len(delays) > 0 && stats.AvgDelay == 0 {
		sum := 0.0
		for _, delay := range delays {
			sum += delay
		}
		stats.AvgDelay = sum / float64(len(delays))
		stats.Success = true
	}
	
	// 如果收到的数据包数大于0，说明至少有部分成功
	if stats.Received > 0 {
		stats.Success = true
	}

	return results, stats
}

// buildPingCommand 根据操作系统构建Ping命令
func (s *PingService) buildPingCommand(req PingTestRequest) []string {
	// 这里使用通用的Ping命令格式，实际项目中可能需要根据不同操作系统调整
	// Linux/macOS 使用 -c 指定次数，-s 指定大小，-W 指定超时
	// Windows 使用 -n 指定次数，-l 指定大小，-w 指定超时
	
	// 暂时使用Linux/macOS格式，实际部署时可能需要根据服务器操作系统调整
	return []string{
		"ping",
		"-c", strconv.Itoa(req.Count),
		"-s", strconv.Itoa(req.Size),
		"-W", strconv.Itoa(req.Timeout / 1000), // 转换为秒
		req.Target,
	}
}

// parsePingResult 解析单次Ping结果
func (s *PingService) parsePingResult(line string, expectedSize int) *PingResult {
	// 检查是否是成功的Ping响应，支持macOS和Linux格式
	// macOS格式：64 bytes from 8.8.8.8: icmp_seq=0 ttl=113 time=21.637 ms
	// Linux格式：64 bytes from 8.8.8.8: icmp_seq=1 ttl=113 time=21.6 ms
	if strings.Contains(line, "bytes from") {
		// 解析IP地址
		var ip string
		var ttl int = 0
		var timeStr string
		var bytes int = expectedSize

		// 解析IP地址，支持"from IP:"格式
		ipMatch := regexp.MustCompile(`from ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+):`).FindStringSubmatch(line)
		if len(ipMatch) > 1 {
			ip = ipMatch[1]
		}

		// 解析字节数
		bytesMatch := regexp.MustCompile(`^(\d+) bytes`).FindStringSubmatch(line)
		if len(bytesMatch) > 1 {
			if b, err := strconv.Atoi(bytesMatch[1]); err == nil {
				bytes = b
			}
		}

		// 解析TTL
		ttlMatch := regexp.MustCompile(`ttl=(\d+)`).FindStringSubmatch(line)
		if len(ttlMatch) > 1 {
			if t, err := strconv.Atoi(ttlMatch[1]); err == nil {
				ttl = t
			}
		}

		// 解析时间，支持time=xx.xxx ms格式
		timeMatch := regexp.MustCompile(`time=([\d.]+)\s*ms`).FindStringSubmatch(line)
		if len(timeMatch) > 1 {
			timeStr = timeMatch[1]
		}

		// 如果IP和时间都解析成功，返回成功结果
		if ip != "" && timeStr != "" {
			return &PingResult{
				Success: true,
				Message: "请求成功",
				IP:      ip,
				Bytes:   bytes,
				Time:    timeStr,
				TTL:     ttl,
			}
		}
	}

	return nil
}

// RegisterRoutes 注册Ping服务的路由
func (s *PingService) RegisterRoutes(api *gin.RouterGroup) {
	ping := api.Group("/ping")
	{
		ping.POST("/test", s.TestPing)
	}
}
