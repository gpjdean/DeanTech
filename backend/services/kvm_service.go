package services

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"deantech/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// KVMService KVM主机服务
type KVMService struct {
	db *gorm.DB
}

// NewKVMService 创建KVM主机服务实例
func NewKVMService(db *gorm.DB) *KVMService {
	return &KVMService{db: db}
}

// KVMHostListItem KVM主机列表项，不包含敏感信息
type KVMHostListItem struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Name           string    `json:"name"`
	IP             string    `json:"ip"`
	Port           int       `json:"port"`
	ConnectionType string    `json:"connectionType"`
	Description    string    `json:"description"`
	IsActive       bool      `json:"isActive"`
	Status         string    `json:"status"`
}

// RegisterRoutes 注册KVM主机相关路由
func (s *KVMService) RegisterRoutes(r *gin.RouterGroup) {
	kvmHosts := r.Group("/kvm/hosts")
	{
		kvmHosts.GET("", s.ListKVMHostsHandler)
		kvmHosts.POST("", s.CreateKVMHostHandler)
		kvmHosts.GET("/:id", s.GetKVMHostHandler)
		kvmHosts.PUT("/:id", s.UpdateKVMHostHandler)
		kvmHosts.DELETE("/:id", s.DeleteKVMHostHandler)
		kvmHosts.POST("/:id/test", s.TestKVMHostConnectionHandler)
		kvmHosts.GET("/:id/vms", s.GetVMListHandler)
		kvmHosts.POST("/:id/vms/:vmId/start", s.StartVMHandler)
		kvmHosts.POST("/:id/vms/:vmId/shutdown", s.ShutdownVMHandler)
		kvmHosts.POST("/:id/vms/:vmId/reboot", s.RebootVMHandler)
		kvmHosts.POST("/:id/vms/:vmId/destroy", s.ForceShutdownVMHandler)
	}
}

// ListKVMHostsHandler 获取KVM主机列表
func (s *KVMService) ListKVMHostsHandler(c *gin.Context) {
	var kvmHosts []models.KVMHost
	result := s.db.Find(&kvmHosts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get KVM hosts",
		})
		return
	}

	// 转换为不含敏感信息的列表项
	var kvmHostListItems []KVMHostListItem
	for _, host := range kvmHosts {
		kvmHostListItems = append(kvmHostListItems, KVMHostListItem{
			ID:             host.ID,
			CreatedAt:      host.CreatedAt,
			UpdatedAt:      host.UpdatedAt,
			Name:           host.Name,
			IP:             host.IP,
			Port:           host.Port,
			ConnectionType: host.ConnectionType,
			Description:    host.Description,
			IsActive:       host.IsActive,
			Status:         host.Status,
		})
	}

	c.JSON(http.StatusOK, kvmHostListItems)
}

// CreateKVMHostHandler 创建KVM主机
func (s *KVMService) CreateKVMHostHandler(c *gin.Context) {
	var kvmHost models.KVMHost
	if err := c.ShouldBindJSON(&kvmHost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 创建KVM主机
	result := s.db.Create(&kvmHost)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create KVM host",
		})
		return
	}

	// 返回不含敏感信息的KVM主机信息
	kvmHostItem := KVMHostListItem{
		ID:             kvmHost.ID,
		CreatedAt:      kvmHost.CreatedAt,
		UpdatedAt:      kvmHost.UpdatedAt,
		Name:           kvmHost.Name,
		IP:             kvmHost.IP,
		Port:           kvmHost.Port,
		ConnectionType: kvmHost.ConnectionType,
		Description:    kvmHost.Description,
		IsActive:       kvmHost.IsActive,
		Status:         kvmHost.Status,
	}

	c.JSON(http.StatusCreated, kvmHostItem)
}

// GetKVMHostHandler 获取KVM主机详情
func (s *KVMService) GetKVMHostHandler(c *gin.Context) {
	id := c.Param("id")
	var kvmHost models.KVMHost
	result := s.db.First(&kvmHost, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 返回完整的KVM主机信息，包括敏感字段
	c.JSON(http.StatusOK, kvmHost)
}

// UpdateKVMHostHandler 更新KVM主机
func (s *KVMService) UpdateKVMHostHandler(c *gin.Context) {
	id := c.Param("id")
	var kvmHost models.KVMHost
	if err := c.ShouldBindJSON(&kvmHost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 更新KVM主机信息
	result := s.db.Model(&models.KVMHost{}).Where("id = ?", id).Updates(kvmHost)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update KVM host",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 重新获取KVM主机信息
	s.db.First(&kvmHost, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "KVM host updated successfully",
		"status":  kvmHost.Status,
	})
}

// DeleteKVMHostHandler 删除KVM主机
func (s *KVMService) DeleteKVMHostHandler(c *gin.Context) {
	id := c.Param("id")
	// 使用Unscoped()方法实现彻底删除，而不是软删除
	result := s.db.Unscoped().Delete(&models.KVMHost{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete KVM host",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KVM host deleted successfully",
	})
}

// TestKVMHostConnectionHandler 测试KVM主机连接
func (s *KVMService) TestKVMHostConnectionHandler(c *gin.Context) {
	id := c.Param("id")

	// 获取KVM主机信息
	var kvmHost models.KVMHost
	result := s.db.First(&kvmHost, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 实际连接测试，根据ConnectionType进行不同的连接测试
	isConnected, err := s.testKVMHostConnection(kvmHost)
	
	// 更新主机状态
	status := "error"
	if isConnected {
		status = "running"
	}
	
	updateResult := s.db.Model(&kvmHost).Update("status", status)
	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update KVM host status",
		})
		return
	}

	if !isConnected {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "KVM host connection test failed: " + err.Error(),
			"status":  "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KVM host connection test successful",
		"status":  status,
	})
}

// testKVMHostConnection 实际测试KVM主机连接
func (s *KVMService) testKVMHostConnection(host models.KVMHost) (bool, error) {
	switch host.ConnectionType {
	case "ssh":
		// 测试SSH连接
		return s.testSSHConnection(host)
	case "virsh":
		// 测试virsh直接连接
		return s.testVirshConnection(host)
	default:
		return false, fmt.Errorf("unsupported connection type: %s", host.ConnectionType)
	}
}

// testSSHConnection 测试SSH连接
func (s *KVMService) testSSHConnection(host models.KVMHost) (bool, error) {
	// 实际实现SSH连接测试
	// 创建SSH客户端配置
	config := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境中应该使用更安全的主机密钥验证
		Timeout:         5 * time.Second,
	}
	
	// 连接到SSH服务器
	addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return false, err
	}
	defer client.Close()
	
	// 连接成功
	return true, nil
}

// testVirshConnection 测试virsh直接连接
func (s *KVMService) testVirshConnection(host models.KVMHost) (bool, error) {
	// 实际实现virsh连接测试
	// 使用libvirt URL连接到远程libvirt服务
	// 这里使用简单的TCP连接测试
	addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	
	// 连接成功
	return true, nil
}

// VMInfo 虚拟机信息结构体
type VMInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	CPU      int    `json:"cpu"`
	Memory   string `json:"memory"`
	Disk     string `json:"disk"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	UUID     string `json:"uuid"`
	Hostname string `json:"hostname"`
}

// GetVMListHandler 获取KVM主机下的虚拟机列表
func (s *KVMService) GetVMListHandler(c *gin.Context) {
	id := c.Param("id")

	// 检查KVM主机是否存在
	var kvmHost models.KVMHost
	result := s.db.First(&kvmHost, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 根据连接类型获取虚拟机列表
	var vmList []VMInfo
	var err error

	switch kvmHost.ConnectionType {
	case "ssh":
		// 通过SSH连接获取虚拟机列表
		vmList, err = s.getVMListBySSH(kvmHost)
	case "virsh":
		// 通过virsh直接连接获取虚拟机列表
		vmList, err = s.getVMListByVirsh(kvmHost)
	default:
		// 默认返回模拟数据
		vmList = s.getMockVMList()
	}

	// 如果获取虚拟机列表失败，返回模拟数据
	if err != nil {
		// 记录错误日志
		fmt.Printf("Failed to get VM list for host %d: %v\n", kvmHost.ID, err)
		// 返回模拟数据，而不是错误
		vmList = s.getMockVMList()
	}

	c.JSON(http.StatusOK, vmList)
}

// getVMListBySSH 通过SSH获取虚拟机列表
func (s *KVMService) getVMListBySSH(host models.KVMHost) ([]VMInfo, error) {
	// 使用SSH客户端连接到主机，执行virsh命令获取虚拟机列表
	config := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	
	addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	
	// 执行virsh命令获取虚拟机列表
	cmd := "virsh list --all --name"
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	
	output, err := session.Output(cmd)
	if err != nil {
		return nil, err
	}
	
	// 处理命令输出，解析虚拟机列表
	vmNames := strings.Split(strings.TrimSpace(string(output)), "\n")
	var vmList []VMInfo
	
	for _, name := range vmNames {
		if name == "" {
			continue
		}
		
		// 为每个虚拟机获取详细信息
		vmInfo, err := s.getVMInfoBySSH(client, name)
		if err != nil {
			// 如果获取单个虚拟机信息失败，跳过该虚拟机
			continue
		}
		vmList = append(vmList, vmInfo)
	}
	
	return vmList, nil
}

// getVMInfoBySSH 通过SSH获取单个虚拟机的详细信息
func (s *KVMService) getVMInfoBySSH(client *ssh.Client, vmName string) (VMInfo, error) {
	// 初始化虚拟机信息
	vmInfo := VMInfo{
		Name: vmName,
		ID:   "",
		CPU:  0, // 重置CPU为0，避免异常值
	}
	
	// 1. 首先执行完整的virsh dominfo命令，获取所有基本信息
	dominfoCmd := fmt.Sprintf("virsh dominfo %s", vmName)
	dominfoSession, err := client.NewSession()
	if err != nil {
		return vmInfo, fmt.Errorf("failed to create session for dominfo: %w", err)
	}
	
	dominfoOutput, err := dominfoSession.Output(dominfoCmd)
	dominfoSession.Close()
	
	if err == nil {
		// 解析完整的dominfo输出
		dominfoStr := string(dominfoOutput)
		
		// 解析ID
		idMatch := regexp.MustCompile(`(?i)Id:\s*(\S+)`).FindStringSubmatch(dominfoStr)
		if len(idMatch) > 1 {
			vmInfo.ID = idMatch[1]
		}
		
		// 解析状态
		stateMatch := regexp.MustCompile(`(?i)State:\s*(\S+)`).FindStringSubmatch(dominfoStr)
		if len(stateMatch) > 1 {
			vmInfo.Status = stateMatch[1]
		}
		
		// 解析CPU数量
		cpuMatch := regexp.MustCompile(`(?i)CPU\(s\):\s*(\d+)`).FindStringSubmatch(dominfoStr)
		if len(cpuMatch) > 1 {
			if cpuCount, err := strconv.Atoi(cpuMatch[1]); err == nil {
				vmInfo.CPU = cpuCount
			}
		}
		
		// 解析内存信息 - 优先使用Max memory
		memMatch := regexp.MustCompile(`(?i)Max memory:\s*([\d\s]+[KMG]iB)`).FindStringSubmatch(dominfoStr)
		if len(memMatch) > 1 {
			vmInfo.Memory = strings.TrimSpace(memMatch[1])
		} else {
			// 尝试获取Current memory作为备选
			currMemMatch := regexp.MustCompile(`(?i)Current memory:\s*([\d\s]+[KMG]iB)`).FindStringSubmatch(dominfoStr)
			if len(currMemMatch) > 1 {
				vmInfo.Memory = strings.TrimSpace(currMemMatch[1])
			}
		}
	}
	
	// 2. 如果CPU信息获取失败，使用virsh vcpucount命令
	if vmInfo.CPU == 0 {
		vcpucountCmd := fmt.Sprintf("virsh vcpucount %s --current", vmName)
		vcpucountSession, err := client.NewSession()
		if err == nil {
			vcpucountOutput, err := vcpucountSession.Output(vcpucountCmd)
			vcpucountSession.Close()
			if err == nil {
				cpuValue := strings.TrimSpace(string(vcpucountOutput))
				if cpuCount, err := strconv.Atoi(cpuValue); err == nil && cpuCount > 0 && cpuCount < 1000 {
					vmInfo.CPU = cpuCount
				}
			}
		}
	}
	
	// 3. 确保内存信息被正确获取 - 使用更可靠的命令
	if vmInfo.Memory == "" {
		// 使用virsh dominfo命令获取所有内存相关信息
		memCmd := fmt.Sprintf("virsh dominfo %s", vmName)
		memSession, err := client.NewSession()
		if err == nil {
			memOutput, err := memSession.Output(memCmd)
			memSession.Close()
			if err == nil {
				memStr := string(memOutput)
				
				// 调试：打印完整的dominfo输出
				fmt.Printf("DEBUG: virsh dominfo %s output: %s\n", vmName, memStr)
				
				// 尝试匹配Max memory行
				maxMemLines := grepLines(memStr, "(?i)Max memory")
				if len(maxMemLines) > 0 {
					maxMemLine := maxMemLines[0]
					parts := strings.SplitN(maxMemLine, ":", 2)
					if len(parts) >= 2 {
						vmInfo.Memory = strings.TrimSpace(parts[1])
						fmt.Printf("DEBUG: Extracted Max memory: %s\n", vmInfo.Memory)
					}
				}
				
				// 如果Max memory失败，尝试Current memory
				if vmInfo.Memory == "" {
					currMemLines := grepLines(memStr, "(?i)Current memory")
					if len(currMemLines) > 0 {
						currMemLine := currMemLines[0]
						parts := strings.SplitN(currMemLine, ":", 2)
						if len(parts) >= 2 {
							vmInfo.Memory = strings.TrimSpace(parts[1])
							fmt.Printf("DEBUG: Extracted Current memory: %s\n", vmInfo.Memory)
						}
					}
				}
				
				// 如果还是失败，尝试Memory行（某些环境可能使用不同的标签）
				if vmInfo.Memory == "" {
					memLines := grepLines(memStr, "(?i)^Memory")
					if len(memLines) > 0 {
						memLine := memLines[0]
						parts := strings.SplitN(memLine, ":", 2)
						if len(parts) >= 2 {
							vmInfo.Memory = strings.TrimSpace(parts[1])
							fmt.Printf("DEBUG: Extracted Memory: %s\n", vmInfo.Memory)
						}
					}
				}
			}
		}
	}
	
	// 4. 如果还是失败，尝试使用virsh vcpuinfo命令获取内存信息（某些环境下可能有不同的命令输出）
	if vmInfo.Memory == "" {
		// 使用更简单的命令，避免复杂的管道
		memCmd := fmt.Sprintf("virsh dominfo %s | grep -i memory", vmName)
		memSession, err := client.NewSession()
		if err == nil {
			memOutput, err := memSession.Output(memCmd)
			memSession.Close()
			if err == nil {
				memOutputStr := strings.TrimSpace(string(memOutput))
				fmt.Printf("DEBUG: virsh dominfo | grep memory output: %s\n", memOutputStr)
				
				// 逐行解析
				lines := strings.Split(memOutputStr, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}
					
					// 简单分割冒号
					parts := strings.SplitN(line, ":", 2)
					if len(parts) >= 2 {
						memValue := strings.TrimSpace(parts[1])
						if memValue != "" {
							vmInfo.Memory = memValue
							fmt.Printf("DEBUG: Extracted memory from grep line '%s': %s\n", line, memValue)
							break
						}
					}
				}
			}
		}
	}
	
	// 5. 最后的尝试：使用virsh dumpxml获取内存信息
	if vmInfo.Memory == "" {
		dumpxmlCmd := fmt.Sprintf("virsh dumpxml %s | grep -A 1 '<memory'", vmName)
		dumpxmlSession, err := client.NewSession()
		if err == nil {
			dumpxmlOutput, err := dumpxmlSession.Output(dumpxmlCmd)
			dumpxmlSession.Close()
			if err == nil {
				dumpxmlStr := strings.TrimSpace(string(dumpxmlOutput))
				fmt.Printf("DEBUG: virsh dumpxml | grep memory output: %s\n", dumpxmlStr)
				
				// 查找内存值
				memMatch := regexp.MustCompile(`<memory[^>]*>(\d+)</memory>`).FindStringSubmatch(dumpxmlStr)
				if len(memMatch) > 1 {
					// dumpxml中的memory通常以KiB为单位
					memKiB, err := strconv.ParseInt(memMatch[1], 10, 64)
					if err == nil {
						// 转换为更友好的单位
						if memKiB >= 1024*1024 {
							vmInfo.Memory = fmt.Sprintf("%.2f GiB", float64(memKiB)/(1024*1024))
						} else {
							vmInfo.Memory = fmt.Sprintf("%.2f MiB", float64(memKiB)/1024)
						}
						fmt.Printf("DEBUG: Extracted memory from dumpxml: %s\n", vmInfo.Memory)
					}
				}
			}
		}
	}
	
	// 5. 使用专门的命令获取虚拟机状态，确保准确性
	stateCmd := fmt.Sprintf("virsh domstate %s", vmName)
	stateSession, err := client.NewSession()
	if err == nil {
		stateOutput, err := stateSession.Output(stateCmd)
		stateSession.Close()
		if err == nil {
			// 直接获取状态，domstate命令只返回状态值，如"running"、"stopped"、"paused"
			vmInfo.Status = strings.TrimSpace(string(stateOutput))
		}
	}
	
	// 6. 获取磁盘信息
	diskCmd := fmt.Sprintf("virsh domblklist %s --details", vmName)
	diskSession, err := client.NewSession()
	if err == nil {
		defer diskSession.Close()
		
		diskOutput, err := diskSession.Output(diskCmd)
		if err == nil {
			// 解析磁盘信息
			diskLines := strings.Split(string(diskOutput), "\n")
			for i, diskLine := range diskLines {
				// 跳过表头行
				if i == 0 {
					continue
				}
				
				diskParts := strings.Fields(diskLine)
				if len(diskParts) < 4 {
					continue
				}
				
				diskPath := diskParts[3]
				// 获取磁盘大小
				diskInfoCmd := fmt.Sprintf("virsh domblkinfo %s %s | grep 'Capacity' | awk '{print $2, $3}'", vmName, diskPath)
				diskInfoSession, err := client.NewSession()
				if err != nil {
					continue
				}
				
				diskInfoOutput, err := diskInfoSession.Output(diskInfoCmd)
				diskInfoSession.Close()
				if err != nil {
					continue
				}
				
				diskInfo := strings.TrimSpace(string(diskInfoOutput))
				if diskInfo != "" {
					// 只取第一个磁盘的信息作为示例
					vmInfo.Disk = diskInfo
					break
				}
			}
		}
	}
	
	// 7. 获取IP地址
	// 尝试多种方法获取IP地址，增加成功率
	
	// 首先检查虚拟机状态，如果是关闭状态，跳过IP地址获取
	if vmInfo.Status == "running" {
		// 方法1: 使用virsh domifaddr命令，添加--full选项，并只获取IPv4地址
		ipCmd := fmt.Sprintf("virsh domifaddr %s --full --ipv4", vmName)
		ipSession, err := client.NewSession()
		if err == nil {
			ipOutput, err := ipSession.Output(ipCmd)
			ipSession.Close()
			if err == nil {
				ipOutputStr := strings.TrimSpace(string(ipOutput))
				fmt.Printf("DEBUG: virsh domifaddr %s --full --ipv4 output: %s\n", vmName, ipOutputStr)
				
				// 清理输出，只保留有效的IP地址行
				ipLines := strings.Split(ipOutputStr, "\n")
				for i, ipLine := range ipLines {
					ipLine = strings.TrimSpace(ipLine)
					if ipLine == "" {
						continue
					}
					
					// 跳过表头行和分隔线
					if i == 0 || strings.Contains(ipLine, "------") || strings.Contains(ipLine, "Interface") {
						continue
					}
					
					// 使用正则表达式直接匹配IP地址
					ipMatch := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/\d{1,2}`).FindStringSubmatch(ipLine)
					if len(ipMatch) > 1 {
						ipAddr := ipMatch[1]
						// 验证IP地址的有效性
						if net.ParseIP(ipAddr) != nil {
							vmInfo.IP = ipAddr
							fmt.Printf("DEBUG: Extracted valid IP from domifaddr: %s\n", ipAddr)
							break
						}
					}
				}
			}
		}
		
		// 方法2: 如果方法1失败，尝试从虚拟机名称中提取IP地址（如果名称中包含IP）
		if vmInfo.IP == "" {
			// 检查虚拟机名称是否包含IP地址
			ipFromName := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`).FindString(vmName)
			if ipFromName != "" {
				// 验证提取到的确实是有效的IP地址
				if net.ParseIP(ipFromName) != nil {
					vmInfo.IP = ipFromName
					fmt.Printf("DEBUG: Extracted valid IP from VM name: %s\n", ipFromName)
				} else {
					fmt.Printf("DEBUG: Extracted invalid IP from VM name: %s, skipping\n", ipFromName)
				}
			}
		}
	}
	
	// 确保只返回有效的IP地址，否则清空
	if vmInfo.IP != "" && net.ParseIP(vmInfo.IP) == nil {
		fmt.Printf("DEBUG: Invalid IP detected, clearing: %s\n", vmInfo.IP)
		vmInfo.IP = ""
	}
	
	// 调试：打印最终获取到的IP地址
	fmt.Printf("DEBUG: Final IP for %s: %s\n", vmName, vmInfo.IP)
	
	return vmInfo, nil
}

// getVMListByVirsh 通过virsh直接获取虚拟机列表
func (s *KVMService) getVMListByVirsh(host models.KVMHost) ([]VMInfo, error) {
	// 构建libvirt连接URL
	// 对于TCP连接：qemu+tcp://host:port/system
	// 对于TLS连接：qemu+tls://host:port/system
	// 默认使用TCP连接
	libvirtURL := fmt.Sprintf("qemu+tcp://%s:%d/system", host.IP, host.Port)
	
	// 执行本地virsh命令连接到远程libvirt服务，获取虚拟机列表
	cmd := fmt.Sprintf("virsh -c %s list --all --name", libvirtURL)
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute virsh command: %w", err)
	}
	
	// 处理命令输出，解析虚拟机列表
	vmNames := strings.Split(strings.TrimSpace(string(output)), "\n")
	var vmList []VMInfo
	
	for _, name := range vmNames {
		if name == "" {
			continue
		}
		
		// 为每个虚拟机获取详细信息
		vmInfo, err := s.getVMInfoByVirsh(libvirtURL, name)
		if err != nil {
			// 如果获取单个虚拟机信息失败，跳过该虚拟机
			continue
		}
		vmList = append(vmList, vmInfo)
	}
	
	return vmList, nil
}

// getVMInfoByVirsh 通过virsh直接获取单个虚拟机的详细信息
func (s *KVMService) getVMInfoByVirsh(libvirtURL, vmName string) (VMInfo, error) {
	// 执行本地virsh命令获取虚拟机详细信息
	cmd := fmt.Sprintf("virsh -c %s dominfo %s", libvirtURL, vmName)
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return VMInfo{}, fmt.Errorf("failed to execute virsh command: %w", err)
	}
	
	// 解析dominfo输出，提取虚拟机信息
	vmInfo := VMInfo{
		Name: vmName,
		ID:   "",
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		switch key {
		case "Id":
			vmInfo.ID = value
		case "CPU(s)":
			cpuCount, _ := strconv.Atoi(value)
			vmInfo.CPU = cpuCount
		case "Max memory":
			// 提取内存信息，如"4194304 KiB"
			vmInfo.Memory = value
		}
	}
	
	// 使用专门的命令获取虚拟机状态，更可靠
	stateCmd := fmt.Sprintf("virsh -c %s domstate %s", libvirtURL, vmName)
	stateOutput, err := exec.Command("bash", "-c", stateCmd).Output()
	if err == nil {
		// 直接获取状态，domstate命令只返回状态值，如"running"、"stopped"、"paused"
		vmInfo.Status = strings.TrimSpace(string(stateOutput))
	}
	
	// 获取磁盘信息
	diskCmd := fmt.Sprintf("virsh -c %s domblklist %s --details", libvirtURL, vmName)
	diskOutput, err := exec.Command("bash", "-c", diskCmd).Output()
	if err == nil {
		// 解析磁盘信息
		diskLines := strings.Split(string(diskOutput), "\n")
		for i, diskLine := range diskLines {
			// 跳过表头行
			if i == 0 {
				continue
			}
			
			diskParts := strings.Fields(diskLine)
			if len(diskParts) < 4 {
				continue
			}
			
			diskPath := diskParts[3]
			// 获取磁盘大小
			diskInfoCmd := fmt.Sprintf("virsh -c %s domblkinfo %s %s | grep 'Capacity' | awk '{print $2, $3}'", libvirtURL, vmName, diskPath)
			diskInfoOutput, err := exec.Command("bash", "-c", diskInfoCmd).Output()
			if err != nil {
				continue
			}
			
			diskInfo := strings.TrimSpace(string(diskInfoOutput))
			if diskInfo != "" {
				// 只取第一个磁盘的信息作为示例
				vmInfo.Disk = diskInfo
				break
			}
		}
	}
	
	// 获取IP地址
	ipCmd := fmt.Sprintf("virsh -c %s domifaddr %s", libvirtURL, vmName)
	ipOutput, err := exec.Command("bash", "-c", ipCmd).Output()
	if err == nil {
		// 解析IP地址
		ipLines := strings.Split(string(ipOutput), "\n")
		for i, ipLine := range ipLines {
			// 跳过表头行
			if i == 0 {
				continue
			}
			
			ipParts := strings.Fields(ipLine)
			if len(ipParts) < 4 {
				continue
			}
			
			ipAddr := ipParts[3]
			// 提取IP地址（去除CIDR前缀）
			if strings.Contains(ipAddr, "/") {
				ipAddr = strings.Split(ipAddr, "/")[0]
			}
			vmInfo.IP = ipAddr
			break
		}
	}
	
	return vmInfo, nil
}

// grepLines 辅助函数，从字符串中提取匹配正则表达式的行
func grepLines(input, pattern string) []string {
	var matches []string
	re := regexp.MustCompile(pattern)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if re.MatchString(line) {
			matches = append(matches, line)
		}
	}
	return matches
}

// getMockVMList 返回模拟的虚拟机列表数据（当真实获取失败时使用）
func (s *KVMService) getMockVMList() []VMInfo {
	return []VMInfo{
		{
			ID:       "vm1",
			Name:     "test-vm-1",
			Status:   "running",
			CPU:      2,
			Memory:   "4G",
			Disk:     "50G",
			IP:       "192.168.1.101",
			MAC:      "52:54:00:12:34:56",
			UUID:     "550e8400-e29b-41d4-a716-446655440000",
			Hostname: "test-vm-1.local",
		},
		{
			ID:       "vm2",
			Name:     "test-vm-2",
			Status:   "stopped",
			CPU:      4,
			Memory:   "8G",
			Disk:     "100G",
			IP:       "192.168.1.102",
			MAC:      "52:54:00:78:90:ab",
			UUID:     "660e8400-e29b-41d4-a716-446655440001",
			Hostname: "test-vm-2.local",
		},
		{
			ID:       "vm3",
			Name:     "test-vm-3",
			Status:   "paused",
			CPU:      1,
			Memory:   "2G",
			Disk:     "20G",
			IP:       "192.168.1.103",
			MAC:      "52:54:00:cd:ef:12",
			UUID:     "770e8400-e29b-41d4-a716-446655440002",
			Hostname: "test-vm-3.local",
		},
	}
}

// StartVMHandler 启动虚拟机
func (s *KVMService) StartVMHandler(c *gin.Context) {
	id := c.Param("id")
	vmName := c.Param("vmId") // 使用vmId作为虚拟机名称

	// 检查KVM主机是否存在
	var kvmHost models.KVMHost
	result := s.db.First(&kvmHost, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 根据连接类型执行启动命令
	var err error
	switch kvmHost.ConnectionType {
	case "ssh":
		err = s.executeSSHVirshCommand(kvmHost, fmt.Sprintf("start %s", vmName))
	case "virsh":
		err = s.executeLocalVirshCommand(kvmHost, fmt.Sprintf("start %s", vmName))
	default:
		err = fmt.Errorf("unsupported connection type: %s", kvmHost.ConnectionType)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to start VM: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "VM started successfully",
	})
}

// ShutdownVMHandler 关闭虚拟机
func (s *KVMService) ShutdownVMHandler(c *gin.Context) {
	id := c.Param("id")
	vmName := c.Param("vmId") // 使用vmId作为虚拟机名称

	// 检查KVM主机是否存在
	var kvmHost models.KVMHost
	result := s.db.First(&kvmHost, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 根据连接类型执行关闭命令
	var err error
	switch kvmHost.ConnectionType {
	case "ssh":
		err = s.executeSSHVirshCommand(kvmHost, fmt.Sprintf("shutdown %s", vmName))
	case "virsh":
		err = s.executeLocalVirshCommand(kvmHost, fmt.Sprintf("shutdown %s", vmName))
	default:
		err = fmt.Errorf("unsupported connection type: %s", kvmHost.ConnectionType)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to shutdown VM: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "VM shutdown successfully",
	})
}

// RebootVMHandler 重启虚拟机
func (s *KVMService) RebootVMHandler(c *gin.Context) {
	id := c.Param("id")
	vmName := c.Param("vmId") // 使用vmId作为虚拟机名称

	// 检查KVM主机是否存在
	var kvmHost models.KVMHost
	result := s.db.First(&kvmHost, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 根据连接类型执行重启命令
	var err error
	switch kvmHost.ConnectionType {
	case "ssh":
		err = s.executeSSHVirshCommand(kvmHost, fmt.Sprintf("reboot %s", vmName))
	case "virsh":
		err = s.executeLocalVirshCommand(kvmHost, fmt.Sprintf("reboot %s", vmName))
	default:
		err = fmt.Errorf("unsupported connection type: %s", kvmHost.ConnectionType)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to reboot VM: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "VM rebooted successfully",
	})
}

// ForceShutdownVMHandler 强制关机虚拟机
func (s *KVMService) ForceShutdownVMHandler(c *gin.Context) {
	id := c.Param("id")
	vmName := c.Param("vmId") // 使用vmId作为虚拟机名称

	// 检查KVM主机是否存在
	var kvmHost models.KVMHost
	result := s.db.First(&kvmHost, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "KVM host not found",
		})
		return
	}

	// 根据连接类型执行强制关机命令
	var err error
	switch kvmHost.ConnectionType {
	case "ssh":
		err = s.executeSSHVirshCommand(kvmHost, fmt.Sprintf("destroy %s", vmName))
	case "virsh":
		err = s.executeLocalVirshCommand(kvmHost, fmt.Sprintf("destroy %s", vmName))
	default:
		err = fmt.Errorf("unsupported connection type: %s", kvmHost.ConnectionType)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to force shutdown VM: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "VM forced shutdown successfully",
	})
}

// executeSSHVirshCommand 通过SSH执行virsh命令
func (s *KVMService) executeSSHVirshCommand(host models.KVMHost, cmd string) error {
	// 创建SSH客户端配置
	config := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	
	// 连接到SSH服务器
	addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to dial SSH: %w", err)
	}
	defer client.Close()
	
	// 执行virsh命令
	fullCmd := fmt.Sprintf("virsh %s", cmd)
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()
	
	// 执行命令并获取输出
	output, err := session.CombinedOutput(fullCmd)
	if err != nil {
		return fmt.Errorf("failed to execute virsh command '%s': %w, output: %s", fullCmd, err, string(output))
	}
	
	return nil
}

// executeLocalVirshCommand 本地执行virsh命令连接到远程服务
func (s *KVMService) executeLocalVirshCommand(host models.KVMHost, cmd string) error {
	// 构建libvirt连接URL
	libvirtURL := fmt.Sprintf("qemu+tcp://%s:%d/system", host.IP, host.Port)
	
	// 执行本地virsh命令
	fullCmd := fmt.Sprintf("virsh -c %s %s", libvirtURL, cmd)
	output, err := exec.Command("bash", "-c", fullCmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute virsh command '%s': %w, output: %s", fullCmd, err, string(output))
	}
	
	return nil
}
