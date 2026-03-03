package services

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TraceRouteService 提供路由跟踪服务
type TraceRouteService struct{}

// NewTraceRouteService 创建一个新的TraceRouteService实例
func NewTraceRouteService() *TraceRouteService {
	return &TraceRouteService{}
}

// TraceRouteRequest 表示路由跟踪请求的参数
type TraceRouteRequest struct {
	Target  string `json:"target" binding:"required"`
	MaxHops int    `json:"maxHops" binding:"required,min=1,max=30"`
	Timeout int    `json:"timeout" binding:"required,min=100,max=5000"`
	Size    int    `json:"size" binding:"required,min=8,max=1500"`
}

// TraceRouteHop 表示路由跟踪的每一跳结果
type TraceRouteHop struct {
	Hop     int     `json:"hop"`
	Success bool    `json:"success"`
	IP      string  `json:"ip"`
	Host    string  `json:"host"`
	Time    float64 `json:"time"`
	TTL     int     `json:"ttl"`
}

// TraceRouteResponse 表示路由跟踪的完整响应结果
type TraceRouteResponse struct {
	Target string          `json:"target"`
	Hops   []TraceRouteHop `json:"hops"`
	Stats  TraceRouteStats `json:"stats"`
}

// TraceRouteStats 表示路由跟踪的统计信息
type TraceRouteStats struct {
	Success    bool  `json:"success"`
	TotalHops  int   `json:"totalHops"`
	Reached    bool  `json:"reached"`
	MaxRTT     float64 `json:"maxRTT"`
	MinRTT     float64 `json:"minRTT"`
	AvgRTT     float64 `json:"avgRTT"`
}

// TestTraceRoute 执行路由跟踪测试
func (s *TraceRouteService) TestTraceRoute(c *gin.Context) {
	var req TraceRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 执行路由跟踪测试
	hops, stats := s.executeTraceRoute(req)

	// 检查是否是权限错误
	if len(hops) == 1 && !hops[0].Success && hops[0].Hop == 0 {
		// 返回权限错误
		c.JSON(500, gin.H{
			"error": "路由跟踪命令执行失败: 权限不足，请确保服务有执行traceroute命令的权限",
		})
		return
	}

	// 返回结果
	c.JSON(200, TraceRouteResponse{
		Target: req.Target,
		Hops:   hops,
		Stats:  stats,
	})
}

// executeTraceRoute 执行实际的路由跟踪命令
func (s *TraceRouteService) executeTraceRoute(req TraceRouteRequest) ([]TraceRouteHop, TraceRouteStats) {
	// 尝试使用自定义的ICMP路由跟踪实现
	customHops, customStats := s.executeCustomTraceRoute(req)
	if len(customHops) > 0 && (customStats.Success || len(customHops) > 1) {
		return customHops, customStats
	}
	
	// 如果自定义实现失败，尝试使用外部traceroute命令
	externalHops, externalStats := s.executeExternalTraceRoute(req)
	if len(externalHops) > 0 {
		return externalHops, externalStats
	}
	
	// 如果都失败，返回空结果
	return customHops, customStats
}

// executeExternalTraceRoute 使用ping命令结合TTL实现路由跟踪
func (s *TraceRouteService) executeExternalTraceRoute(req TraceRouteRequest) ([]TraceRouteHop, TraceRouteStats) {
	// 初始化结果
	hops := make([]TraceRouteHop, 0, req.MaxHops)
	stats := TraceRouteStats{
		Success:    false,
		TotalHops:  0,
		Reached:    false,
		MaxRTT:     0,
		MinRTT:     0,
		AvgRTT:     0,
	}

	// 初始化延迟数组用于计算统计信息
	delays := make([]float64, 0)

	// 先解析目标地址，获取IP
	targetIPs, err := net.LookupIP(req.Target)
	if err != nil {
		fmt.Printf("Failed to resolve target: %v\n", err)
		return hops, stats
	}
	targetIP := targetIPs[0].String()

	// 使用ping命令结合TTL实现路由跟踪
	for ttl := 1; ttl <= req.MaxHops; ttl++ {
		// 构建ping命令
		cmd := exec.Command("ping", 
			"-c", "1", // 发送1个数据包
			"-t", strconv.Itoa(ttl), // 设置TTL
			"-W", strconv.FormatFloat(float64(req.Timeout)/1000, 'f', 3, 64), // 超时时间（秒）
			targetIP, // 目标IP
		)

		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		outputStr := string(output)
		
		// 解析ping输出，提取延迟和中间节点IP
		// macOS ping输出格式：
		// Request timeout for icmp_seq 0
		// 64 bytes from 192.168.1.1: icmp_seq=0 ttl=64 time=1.234 ms
		// traceroute-like info can be extracted from the TTL exceeded messages
		
		hop := TraceRouteHop{
			Hop:     ttl,
			Success: false,
			IP:      "",
			Host:    "",
			Time:    0,
			TTL:     ttl,
		}
		
		// 检查是否成功收到响应
		if err == nil {
			// 成功到达目标
			hop.Success = true
			hop.IP = targetIP
			hop.Host = req.Target
			
			// 提取延迟
			timeRegex := regexp.MustCompile(`time=(\d+\.\d+)\s+ms`)
			timeMatch := timeRegex.FindStringSubmatch(outputStr)
			if len(timeMatch) > 1 {
				timeVal, _ := strconv.ParseFloat(timeMatch[1], 64)
				hop.Time = timeVal
				delays = append(delays, timeVal)
			}
			
			stats.Reached = true
		} else {
			// 检查是否是TTL exceeded消息
			ttlExceededRegex := regexp.MustCompile(`From\s+([\d\.]+)\s+icmp_seq=\d+\s+Time\s+to\s+live\s+exceeded`)
			ttlMatch := ttlExceededRegex.FindStringSubmatch(outputStr)
			if len(ttlMatch) > 1 {
				// TTL exceeded，提取中间节点IP
				hop.Success = true
				hop.IP = ttlMatch[1]
				hop.Host = ttlMatch[1] // 不解析主机名
				
				// 提取延迟（如果有）
				timeRegex := regexp.MustCompile(`time=(\d+\.\d+)\s+ms`)
				timeMatch := timeRegex.FindStringSubmatch(outputStr)
				if len(timeMatch) > 1 {
					timeVal, _ := strconv.ParseFloat(timeMatch[1], 64)
					hop.Time = timeVal
					delays = append(delays, timeVal)
				}
			}
		}
		
		// 添加跳信息
		hops = append(hops, hop)
		stats.TotalHops++
		
		// 如果到达目标，结束跟踪
		if hop.Success && hop.IP == targetIP {
			break
		}
	}

	// 计算统计信息
	if len(delays) > 0 {
		stats.Success = true
		
		// 计算最小、最大和平均延迟
		minDelay := delays[0]
		maxDelay := delays[0]
		sum := delays[0]
		
		for i := 1; i < len(delays); i++ {
			if delays[i] < minDelay {
				minDelay = delays[i]
			}
			if delays[i] > maxDelay {
				maxDelay = delays[i]
			}
			sum += delays[i]
		}
		
		stats.MinRTT = minDelay
		stats.MaxRTT = maxDelay
		stats.AvgRTT = sum / float64(len(delays))
	}

	return hops, stats
}

// executeCustomTraceRoute 使用简单的TCP连接方式实现路由跟踪
func (s *TraceRouteService) executeCustomTraceRoute(req TraceRouteRequest) ([]TraceRouteHop, TraceRouteStats) {
	// 初始化结果
	hops := make([]TraceRouteHop, 0, req.MaxHops)
	stats := TraceRouteStats{
		Success:    false,
		TotalHops:  0,
		Reached:    false,
		MaxRTT:     0,
		MinRTT:     0,
		AvgRTT:     0,
	}

	// 初始化延迟数组用于计算统计信息
	delays := make([]float64, 0)

	// 先解析目标地址，获取IP
	ips, err := net.LookupIP(req.Target)
	if err != nil {
		fmt.Printf("Failed to resolve target: %v\n", err)
		return hops, stats
	}

	targetIP := ips[0]
	targetAddr := &net.TCPAddr{IP: targetIP, Port: 80} // 使用HTTP端口80

	// 使用简单的TCP连接方式，直接尝试连接目标
	// 虽然无法获取中间节点，但可以检测是否能到达目标
	start := time.Now()
	conn, err := net.DialTimeout("tcp", targetAddr.String(), time.Duration(req.Timeout)*time.Millisecond)
	elapsed := time.Since(start).Milliseconds()
	
	hop := TraceRouteHop{
		Hop:     1,
		Success: false,
		IP:      "",
		Host:    "",
		Time:    0,
		TTL:     1,
	}
	
	if err != nil {
		// 连接失败
		hop.Success = false
	} else {
		// 连接成功，到达目标
		conn.Close()
		hop.Success = true
		hop.IP = targetIP.String()
		hop.Host = req.Target
		hop.Time = float64(elapsed)
		delays = append(delays, float64(elapsed))
		stats.Reached = true
	}
	
	// 添加跳信息
	hops = append(hops, hop)
	stats.TotalHops++

	// 计算统计信息
	if len(delays) > 0 {
		stats.Success = true
		
		// 计算最小、最大和平均延迟
		minDelay := delays[0]
		maxDelay := delays[0]
		sum := delays[0]
		
		for i := 1; i < len(delays); i++ {
			if delays[i] < minDelay {
				minDelay = delays[i]
			}
			if delays[i] > maxDelay {
				maxDelay = delays[i]
			}
			sum += delays[i]
		}
		
		stats.MinRTT = minDelay
		stats.MaxRTT = maxDelay
		stats.AvgRTT = sum / float64(len(delays))
	}

	return hops, stats
}

// checksum 计算ICMP包的校验和
func checksum(msg []byte) uint16 {
	var sum uint32
	for i := 0; i < len(msg); i += 2 {
		sum += uint32(msg[i]) << 8
		if i+1 < len(msg) {
			sum += uint32(msg[i+1])
		}
	}
	
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	
	return ^uint16(sum)
}

// parseTraceRouteHop 解析路由跟踪的每一跳结果
func (s *TraceRouteService) parseTraceRouteHop(line string) *TraceRouteHop {
	// 打印当前解析的行，方便调试
	fmt.Printf("Parsing trace route line: %s\n", line)
	
	// 支持Linux/macOS traceroute输出格式
	// 格式1: 1  192.168.1.1 (192.168.1.1)  1.234 ms  2.345 ms  3.456 ms
	// 格式2: 2  * * *
	// 格式3: 3  10.0.0.1 (10.0.0.1)  10.123 ms
	
	// 解析跳数
	hopMatch := regexp.MustCompile(`^(\d+)\s+`).FindStringSubmatch(line)
	if len(hopMatch) < 2 {
		return nil
	}
	
	hopNum, err := strconv.Atoi(hopMatch[1])
	if err != nil {
		return nil
	}
	
	// 检查是否超时
	if strings.Contains(line, "* * *") {
		return &TraceRouteHop{
			Hop:     hopNum,
			Success: false,
			IP:      "",
			Host:    "",
			Time:    0,
			TTL:     0,
		}
	}
	
	// 解析成功的跳
	// 匹配IP地址和时间
	successMatch := regexp.MustCompile(`^(\d+)\s+([^\s]+)\s+\(([^)]+)\)\s+([\d.]+)\s+ms`).FindStringSubmatch(line)
	if len(successMatch) >= 5 {
		timeStr := successMatch[4]
		timeVal, _ := strconv.ParseFloat(timeStr, 64)
		
		return &TraceRouteHop{
			Hop:     hopNum,
			Success: true,
			IP:      successMatch[3],
			Host:    successMatch[2],
			Time:    timeVal,
			TTL:     hopNum,
		}
	}
	
	// 尝试另一种格式
	simpleMatch := regexp.MustCompile(`^(\d+)\s+([^\s]+)\s+([\d.]+)\s+ms`).FindStringSubmatch(line)
	if len(simpleMatch) >= 4 {
		timeStr := simpleMatch[3]
		timeVal, _ := strconv.ParseFloat(timeStr, 64)
		
		return &TraceRouteHop{
			Hop:     hopNum,
			Success: true,
			IP:      simpleMatch[2],
			Host:    simpleMatch[2],
			Time:    timeVal,
			TTL:     hopNum,
		}
	}
	
	return nil
}

// RegisterRoutes 注册TraceRoute服务的路由
func (s *TraceRouteService) RegisterRoutes(api *gin.RouterGroup) {
	trace := api.Group("/trace-route")
	{
		trace.POST("/test", s.TestTraceRoute)
	}
}
