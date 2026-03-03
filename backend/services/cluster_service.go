package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	"deantech/models"
)

// ClusterService 集群服务结构体
type ClusterService struct {
	db *gorm.DB
}

// NewClusterService 创建集群服务实例
func NewClusterService(db *gorm.DB) *ClusterService {
	return &ClusterService{db: db}
}

// RegisterRoutes 注册集群相关路由
func (s *ClusterService) RegisterRoutes(router *gin.RouterGroup) {
	clusterGroup := router.Group("/clusters")
	{
		clusterGroup.GET("", s.ListClustersHandler)
		clusterGroup.GET("/:id", s.GetClusterHandler)
		clusterGroup.POST("", s.CreateClusterHandler)
		clusterGroup.PUT("/:id", s.UpdateClusterHandler)
		clusterGroup.DELETE("/:id", s.DeleteClusterHandler)
		clusterGroup.GET("/:id/nodes", s.ListNodesHandler)
		clusterGroup.GET("/:id/pods", s.ListPodsHandler)
		clusterGroup.GET("/:id/pods/:name/logs", s.GetPodLogsHandler)
		clusterGroup.GET("/:id/deployments", s.ListDeploymentsHandler)
		clusterGroup.POST("/:id/deployments", s.CreateDeploymentHandler)
		clusterGroup.POST("/:id/deployments/precheck", s.PrecheckDeploymentHandler)
		clusterGroup.PUT("/:id/deployments/:name/scale", s.ScaleDeploymentHandler)
		clusterGroup.POST("/:id/deployments/:name/restart", s.RestartDeploymentHandler)
		clusterGroup.GET("/:id/statefulsets", s.ListStatefulSetsHandler)
		clusterGroup.GET("/:id/daemonsets", s.ListDaemonSetsHandler)
		clusterGroup.GET("/:id/jobs", s.ListJobsHandler)
		clusterGroup.GET("/:id/cronjobs", s.ListCronJobsHandler)
		clusterGroup.GET("/:id/services", s.ListServicesHandler)
		clusterGroup.GET("/:id/configmaps", s.ListConfigMapsHandler)
		clusterGroup.GET("/:id/secrets", s.ListSecretsHandler)
		clusterGroup.GET("/:id/ingresses", s.ListIngressesHandler)
		clusterGroup.GET("/:id/namespaces", s.ListNamespacesHandler)
		clusterGroup.POST("/:id/namespaces", s.CreateNamespaceHandler)
		clusterGroup.DELETE("/:id/namespaces/:name", s.DeleteNamespaceHandler)
		clusterGroup.GET("/:id/replicasets", s.ListReplicaSetsHandler)
		clusterGroup.GET("/:id/pvcs", s.ListPVCsHandler)
		clusterGroup.GET("/:id/pvs", s.ListPVsHandler)
		clusterGroup.GET("/:id/storageclasses", s.ListStorageClassesHandler)
		clusterGroup.GET("/:id/storageclasses/:name/yaml", s.GetStorageClassYamlHandler)
		clusterGroup.DELETE("/:id/storageclasses/:name", s.DeleteStorageClassHandler)
		clusterGroup.GET("/:id/events", s.ListEventsHandler)
		clusterGroup.GET("/:id/dashboard", s.GetDashboardHandler)
		clusterGroup.POST("/:id/test-connection", s.TestConnectionHandler)
		clusterGroup.GET("/:id/pods/:name/exec", s.ExecPodHandler)
	}
}

// ListClustersHandler 获取集群列表
func (s *ClusterService) ListClustersHandler(c *gin.Context) {
	var clusters []models.Cluster
	result := s.db.Find(&clusters)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, clusters)
}

// GetClusterHandler 获取单个集群
func (s *ClusterService) GetClusterHandler(c *gin.Context) {
	id := c.Param("id")
	var cluster models.Cluster
	result := s.db.First(&cluster, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}
	c.JSON(http.StatusOK, cluster)
}

// CreateClusterHandler 创建集群
func (s *ClusterService) CreateClusterHandler(c *gin.Context) {
	var cluster models.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查集群名称是否已存在
	var existingCluster models.Cluster
	result := s.db.Where("name = ?", cluster.Name).First(&existingCluster)
	if result.Error == nil {
		// 集群名称已存在
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster name already exists", "message": "集群名称已存在"})
		return
	}

	// 如果IsActive为nil，设置默认值为true
	if cluster.IsActive == nil {
		trueVal := true
		cluster.IsActive = &trueVal
	}

	// 设置连接类型
	if cluster.Kubeconfig != "" {
		cluster.ConnectionType = "kubeconfig"
		// 从KubeConfig文件中解析API地址
		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
		if err == nil {
			cluster.APIURL = config.Host
		}
	} else if cluster.Token != "" {
		cluster.ConnectionType = "token"
	}

	// 保存集群到数据库
	result = s.db.Create(&cluster)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error(), "message": "创建集群失败"})
		return
	}

	// 测试集群连接并更新状态
	testClusterConnection(&cluster)
	s.db.Save(&cluster)

	c.JSON(http.StatusCreated, cluster)
}

// UpdateClusterHandler 更新集群
func (s *ClusterService) UpdateClusterHandler(c *gin.Context) {
	id := c.Param("id")
	var cluster models.Cluster
	result := s.db.First(&cluster, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found", "message": "集群不存在"})
		return
	}

	var updatedCluster models.Cluster
	if err := c.ShouldBindJSON(&updatedCluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "请求参数错误"})
		return
	}

	// 检查集群名称是否已被其他集群使用
	if updatedCluster.Name != cluster.Name {
		var existingCluster models.Cluster
		result := s.db.Where("name = ? AND id != ?", updatedCluster.Name, id).First(&existingCluster)
		if result.Error == nil {
			// 集群名称已存在
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster name already exists", "message": "集群名称已存在"})
			return
		}
	}

	// 更新集群信息
	cluster.Name = updatedCluster.Name
	cluster.Description = updatedCluster.Description
	cluster.Token = updatedCluster.Token
	cluster.Kubeconfig = updatedCluster.Kubeconfig
	// 只有当updatedCluster.IsActive不为nil时才更新，否则保持原状态
	if updatedCluster.IsActive != nil {
		cluster.IsActive = updatedCluster.IsActive
	}

	// 设置连接类型
	if cluster.Kubeconfig != "" {
		cluster.ConnectionType = "kubeconfig"
		// 从KubeConfig文件中解析API地址
		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
		if err == nil {
			cluster.APIURL = config.Host
		}
	} else if cluster.Token != "" {
		cluster.ConnectionType = "token"
		// 如果没有KubeConfig，使用手动输入的API地址
		cluster.APIURL = updatedCluster.APIURL
	}

	// 测试集群连接并更新状态
	testClusterConnection(&cluster)

	result = s.db.Save(&cluster)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error(), "message": "更新集群失败"})
		return
	}

	c.JSON(http.StatusOK, cluster)
}

// DeleteClusterHandler 删除集群
func (s *ClusterService) DeleteClusterHandler(c *gin.Context) {
	id := c.Param("id")

	// 直接彻底删除集群本身
	if err := s.db.Unscoped().Delete(&models.Cluster{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "删除集群失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cluster deleted successfully", "success": true})
}

// ListNodesHandler 获取集群节点列表
func (s *ClusterService) ListNodesHandler(c *gin.Context) {
	clusterID := c.Param("id")
	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s节点获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取节点数据
	k8sNodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list nodes: " + err.Error(),
		})
		return
	}

	// 获取节点使用率数据
	nodeUsageMap := make(map[string]struct {
		CPUUsage    float64
		MemoryUsage float64
	})

	// 检查kubectl命令是否存在
	kubectlPath, err := exec.LookPath("kubectl")
	if err == nil {
		// 构建kubectl命令
		var cmd *exec.Cmd
		// 设置环境变量，确保kubectl能正确执行
		env := os.Environ()

		if cluster.Kubeconfig != "" {
			// 使用kubeconfig文件
			tempFile, err := os.CreateTemp("", "kubeconfig-")
			if err == nil {
				defer func() {
					tempFile.Close()
					os.Remove(tempFile.Name())
				}()
				tempFile.WriteString(cluster.Kubeconfig)
				tempFile.Sync()
				// 构建命令，使用kubeconfig文件
				cmd = exec.Command(kubectlPath, "top", "node")
				// 添加环境变量
				env = append(env, "KUBECONFIG="+tempFile.Name())
				cmd.Env = env
				fmt.Printf("Executing kubectl with kubeconfig: %s\n", tempFile.Name())
			}
		} else {
			// 使用token方式
			cmd = exec.Command(kubectlPath, "top", "node", "--server", cluster.APIURL, "--token", cluster.Token, "--insecure-skip-tls-verify")
			cmd.Env = env
			fmt.Printf("Executing kubectl with token: %s\n", cluster.APIURL)
		}

		if cmd != nil {
			// 执行命令获取输出
			output, err := cmd.CombinedOutput()
			fmt.Printf("kubectl command output: %s\n", string(output))
			if err == nil {
				// 解析kubectl top node输出
				lines := strings.Split(string(output), "\n")
				for i, line := range lines {
					// 跳过表头
					if i == 0 {
						continue
					}

					// 跳过空行
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}

					// 分割行数据
					parts := strings.Fields(line)
					if len(parts) < 5 {
						fmt.Printf("Invalid line format: %s\n", line)
						continue
					}

					// 节点名称
					nodeName := parts[0]

					// 解析CPU使用率 - 使用第3列的百分比数据
					cpuUsageStr := parts[2] // CPU% 列
					cpuUsageStr = strings.TrimSuffix(cpuUsageStr, "%")
					cpuUsage, err := strconv.ParseFloat(cpuUsageStr, 64)
					if err != nil {
						fmt.Printf("Failed to parse CPU usage: %s\n", cpuUsageStr)
						continue
					}

					// 解析内存使用率 - 使用第5列的百分比数据
					memoryUsageStr := parts[4] // MEMORY% 列
					memoryUsageStr = strings.TrimSuffix(memoryUsageStr, "%")
					memoryUsage, err := strconv.ParseFloat(memoryUsageStr, 64)
					if err != nil {
						fmt.Printf("Failed to parse memory usage: %s\n", memoryUsageStr)
						continue
					}

					// 存储到map中
					nodeUsageMap[nodeName] = struct {
						CPUUsage    float64
						MemoryUsage float64
					}{CPUUsage: cpuUsage, MemoryUsage: memoryUsage}
					fmt.Printf("Node %s: CPU %.2f%%, Memory %.2f%%\n", nodeName, cpuUsage, memoryUsage)
				}
			} else {
				// 打印命令执行错误，便于调试
				fmt.Printf("kubectl top node command error: %s, output: %s\n", err.Error(), string(output))
			}
		}
	}

	// 转换为我们的节点模型
	var nodes []models.Node
	for _, node := range k8sNodes.Items {
		// 从节点标签和污点中获取角色
		role := "worker"

		// 检查节点标签
		if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
			role = "master"
		} else if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
			role = "master"
		}

		// 检查节点污点，如果是master节点通常会有taint
		if role == "worker" {
			for _, taint := range node.Spec.Taints {
				if strings.Contains(taint.Key, "master") || strings.Contains(taint.Key, "control-plane") {
					role = "master"
					break
				}
			}
		}

		// 获取节点状态
		status := "NotReady"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status == corev1.ConditionTrue {
					status = "Ready"
				}
				break
			}
		}

		// 解析CPU数量
		cpu := 0
		if node.Status.Capacity.Cpu().Value() > 0 {
			cpu = int(node.Status.Capacity.Cpu().Value())
		}

		// 解析内存大小（转换为MB）
		memory := int64(0)
		if node.Status.Capacity.Memory().Value() > 0 {
			memory = node.Status.Capacity.Memory().Value() / (1024 * 1024)
		}

		// 解析容器运行时
		containerRuntime := node.Status.NodeInfo.ContainerRuntimeVersion

		// 解析Pod容量
		pods := 0
		if node.Status.Capacity.Pods().Value() > 0 {
			pods = int(node.Status.Capacity.Pods().Value())
		}

		// 获取节点使用率
		cpuUsage := 0.0
		memoryUsage := 0.0
		if usage, exists := nodeUsageMap[node.Name]; exists {
			cpuUsage = usage.CPUUsage
			memoryUsage = usage.MemoryUsage
			fmt.Printf("Using real metrics for node %s: CPU %.2f%%, Memory %.2f%%\n", node.Name, cpuUsage, memoryUsage)
		} else {
			// 没有真实数据时，使用率保持为0
			fmt.Printf("No real metrics available for node %s, using 0%% for CPU and Memory\n", node.Name)
		}

		nodes = append(nodes, models.Node{
			ClusterID:        cluster.ID,
			Name:             node.Name,
			IP:               node.Status.Addresses[0].Address,
			Role:             role,
			Status:           status,
			OS:               node.Status.NodeInfo.OSImage,
			KernelVersion:    node.Status.NodeInfo.KernelVersion,
			ContainerRuntime: containerRuntime,
			CPU:              cpu,
			Memory:           memory,
			Pods:             pods,
			CPUUsage:         cpuUsage,
			MemoryUsage:      memoryUsage,
		})
	}

	c.JSON(http.StatusOK, nodes)
}

// ListPodsHandler 获取集群Pod列表
func (s *ClusterService) ListPodsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s Pod获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从请求中获取命名空间参数
	namespace := c.Query("namespace")

	// 从K8s集群获取Pod数据
	// 如果提供了命名空间，则只获取该命名空间的Pod，否则获取所有命名空间的Pod
	var pods *corev1.PodList
	if namespace != "" {
		// 获取指定命名空间的Pod
		pods, err = clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的Pod
		pods, err = clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list pods: " + err.Error(),
		})
		return
	}

	// 转换为我们的Pod模型
	var resultPods []models.Pod
	for _, pod := range pods.Items {
		// 获取Pod状态
		status := string(pod.Status.Phase)

		// 转换重启次数
		restartCount := int32(0)
		for _, container := range pod.Status.ContainerStatuses {
			restartCount += container.RestartCount
		}

		// 转换容器列表
		var containers []models.PodContainer
		for _, container := range pod.Spec.Containers {
			// 获取容器状态
			containerStatus := "Waiting"
			containerReady := false
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.Name == container.Name {
					if cs.Ready {
						containerStatus = "Running"
						containerReady = true
					} else if cs.State.Terminated != nil {
						containerStatus = "Terminated"
					}
					break
				}
			}

			containers = append(containers, models.PodContainer{
				Name:   container.Name,
				Image:  container.Image,
				Status: containerStatus,
				Ready:  containerReady,
			})
		}

		resultPods = append(resultPods, models.Pod{
			ClusterID:    cluster.ID,
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Status:       status,
			NodeName:     pod.Spec.NodeName,
			RestartCount: int(restartCount),
			PodIP:        pod.Status.PodIP,
			HostIP:       pod.Status.HostIP,
			CreatedAt:    pod.CreationTimestamp.Time,
			Containers:   containers,
		})
	}

	// 如果resultPods为空，初始化一个空切片，确保返回[]而不是null
	if resultPods == nil {
		resultPods = []models.Pod{}
	}

	c.JSON(http.StatusOK, resultPods)
}

// ListDeploymentsHandler 获取集群Deployment列表
func (s *ClusterService) ListDeploymentsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s Deployment获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取Deployment数据
	// 如果提供了命名空间，则只获取该命名空间的Deployment，否则获取所有命名空间的Deployment
	var deployments *appsv1.DeploymentList
	if namespace != "" {
		// 获取指定命名空间的Deployment
		deployments, err = clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的Deployment
		deployments, err = clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list deployments: " + err.Error(),
		})
		return
	}

	// 转换为我们的Deployment模型
	var resultDeployments []models.Deployment
	for _, deployment := range deployments.Items {
		// 获取Deployment状态
		status := "Unknown"
		if deployment.Status.AvailableReplicas == deployment.Status.Replicas {
			status = "Available"
		} else if deployment.Status.UpdatedReplicas > 0 {
			status = "Progressing"
		} else {
			status = "Degraded"
		}

		// 获取第一个容器的镜像
		image := ""
		if len(deployment.Spec.Template.Spec.Containers) > 0 {
			image = deployment.Spec.Template.Spec.Containers[0].Image
		}

		// 处理Replicas指针类型
		replicas := 0
		if deployment.Spec.Replicas != nil {
			replicas = int(*deployment.Spec.Replicas)
		}

		resultDeployments = append(resultDeployments, models.Deployment{
			ClusterID:     cluster.ID,
			Name:          deployment.Name,
			Namespace:     deployment.Namespace,
			Replicas:      replicas,
			ReadyReplicas: int(deployment.Status.ReadyReplicas),
			Status:        status,
			Image:         image,
			CreatedAt:     deployment.CreationTimestamp.Time,
		})
	}

	// 如果resultDeployments为空，初始化一个空切片，确保返回[]而不是null
	if resultDeployments == nil {
		resultDeployments = []models.Deployment{}
	}

	c.JSON(http.StatusOK, resultDeployments)
}

// ListNamespacesHandler 获取集群命名空间列表
func (s *ClusterService) ListNamespacesHandler(c *gin.Context) {
	clusterID := c.Param("id")
	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s命名空间获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取命名空间数据
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list namespaces: " + err.Error(),
		})
		return
	}

	// 转换为完整的命名空间对象列表
	var namespaceList []map[string]interface{}
	for _, ns := range namespaces.Items {
		// 获取命名空间状态
		status := "Active"
		if ns.Status.Phase == corev1.NamespaceTerminating {
			status = "Terminating"
		}

		// 构建命名空间对象，将时间转换为 ISO 格式字符串
		namespace := map[string]interface{}{
			"name":        ns.Name,
			"status":      status,
			"createdAt":   ns.CreationTimestamp.Time.Format(time.RFC3339),
			"labels":      ns.Labels,
			"annotations": ns.Annotations,
		}

		namespaceList = append(namespaceList, namespace)
	}

	// 处理搜索和筛选
	keyword := c.Query("keyword")
	statusFilter := c.Query("status")

	// 过滤命名空间
	filteredNamespaces := namespaceList
	if keyword != "" {
		var temp []map[string]interface{}
		for _, ns := range filteredNamespaces {
			if name, ok := ns["name"].(string); ok {
				if strings.Contains(name, keyword) {
					temp = append(temp, ns)
				}
			}
		}
		filteredNamespaces = temp
	}

	if statusFilter != "" {
		var temp []map[string]interface{}
		for _, ns := range filteredNamespaces {
			if status, ok := ns["status"].(string); ok {
				if status == statusFilter {
					temp = append(temp, ns)
				}
			}
		}
		filteredNamespaces = temp
	}

	// 处理分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 计算总条数
	total := len(filteredNamespaces)

	// 计算分页起始位置
	start := (page - 1) * pageSize
	end := start + pageSize

	// 确保起始位置不小于0
	if start < 0 {
		start = 0
	}

	// 确保结束位置不大于总条数
	if end > total {
		end = total
	}

	// 截取分页数据
	var paginatedNamespaces []map[string]interface{}
	if start < total {
		paginatedNamespaces = filteredNamespaces[start:end]
	} else {
		paginatedNamespaces = []map[string]interface{}{}
	}

	// 返回分页数据
	c.JSON(http.StatusOK, gin.H{
		"data":     paginatedNamespaces,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateNamespaceHandler 创建集群命名空间
func (s *ClusterService) CreateNamespaceHandler(c *gin.Context) {
	clusterID := c.Param("id")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 解析请求体
	var namespaceReq struct {
		Name        string            `json:"name" binding:"required"`
		Status      string            `json:"status"`
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
	}

	if err := c.ShouldBindJSON(&namespaceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 实现实际的K8s命名空间创建逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 创建命名空间
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespaceReq.Name,
			Labels:      namespaceReq.Labels,
			Annotations: namespaceReq.Annotations,
		},
	}

	createdNamespace, err := clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create namespace: " + err.Error(),
		})
		return
	}

	// 返回创建成功的命名空间信息
	c.JSON(http.StatusCreated, gin.H{
		"message": "Namespace created successfully",
		"success": true,
		"data": map[string]interface{}{
			"name":        createdNamespace.Name,
			"status":      string(createdNamespace.Status.Phase),
			"createdAt":   createdNamespace.CreationTimestamp.Time.Format(time.RFC3339),
			"labels":      createdNamespace.Labels,
			"annotations": createdNamespace.Annotations,
		},
	})
}

// DeleteNamespaceHandler 删除集群命名空间
func (s *ClusterService) DeleteNamespaceHandler(c *gin.Context) {
	clusterID := c.Param("id")
	namespaceName := c.Param("name")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s命名空间删除逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 删除命名空间
	deletePolicy := metav1.DeletePropagationForeground
	err = clientset.CoreV1().Namespaces().Delete(context.Background(), namespaceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete namespace: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Namespace deleted successfully",
		"success": true,
	})
}

// ScaleDeploymentHandler 扩缩容Deployment
func (s *ClusterService) ScaleDeploymentHandler(c *gin.Context) {
	clusterID := c.Param("id")
	name := c.Param("name")

	// 获取请求参数
	var scaleRequest struct {
		Namespace string `json:"namespace" binding:"required"`
		Replicas  int    `json:"replicas" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&scaleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found: " + result.Error.Error()})
		return
	}

	// 使用kubectl命令扩缩容Deployment
	var kubeconfigPath string
	var err error
	var tempFile *os.File

	if cluster.Kubeconfig != "" {
		// 创建临时kubeconfig文件
		tempFile, err = os.CreateTemp("", "kubeconfig-*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary kubeconfig file: " + err.Error(), "details": err.Error(), "step": "create_temp_file"})
			return
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		// 写入kubeconfig内容
		if _, err := tempFile.WriteString(cluster.Kubeconfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write kubeconfig file: " + err.Error(), "details": err.Error(), "step": "write_kubeconfig"})
			return
		}

		// 确保文件内容已写入
		if err := tempFile.Sync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync kubeconfig file: " + err.Error(), "details": err.Error(), "step": "sync_kubeconfig"})
			return
		}

		kubeconfigPath = tempFile.Name()
	} else {
		// 如果没有kubeconfig，使用默认位置
		kubeconfigPath = ""
	}

	// 检查kubectl命令是否存在
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kubectl command not found: " + err.Error(), "details": err.Error(), "step": "check_kubectl"})
		return
	}

	// 构建kubectl命令
	cmd := exec.Command(kubectlPath, "scale", "deployment", name, "--replicas", strconv.Itoa(scaleRequest.Replicas), "-n", scaleRequest.Namespace)

	// 设置KUBECONFIG环境变量
	if kubeconfigPath != "" {
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute kubectl scale command: " + err.Error(), "output": string(output), "details": err.Error(), "command": cmd.String(), "step": "execute_command"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deployment scaled successfully", "output": string(output), "command": cmd.String()})
}

// RestartDeploymentHandler 重启Deployment
func (s *ClusterService) RestartDeploymentHandler(c *gin.Context) {
	clusterID := c.Param("id")
	name := c.Param("name")

	// 获取请求参数
	var restartRequest struct {
		Namespace string `json:"namespace" binding:"required"`
	}

	if err := c.ShouldBindJSON(&restartRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 使用kubectl命令重启Deployment
	var kubeconfigPath string
	var err error
	var tempFile *os.File

	if cluster.Kubeconfig != "" {
		// 创建临时kubeconfig文件
		tempFile, err = os.CreateTemp("", "kubeconfig-*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create temporary kubeconfig file: " + err.Error(),
			})
			return
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		// 写入kubeconfig内容
		if _, err := tempFile.WriteString(cluster.Kubeconfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to write kubeconfig file: " + err.Error(),
			})
			return
		}

		kubeconfigPath = tempFile.Name()
	} else {
		// 如果没有kubeconfig，使用默认位置
		kubeconfigPath = ""
	}

	// 构建kubectl命令
	cmd := exec.Command("kubectl", "rollout", "restart", "deployment", name, "-n", restartRequest.Namespace)

	// 设置KUBECONFIG环境变量
	if kubeconfigPath != "" {
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to restart deployment: " + err.Error(),
			"output": string(output),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deployment restarted successfully", "output": string(output)})
}

// ListServicesHandler 获取集群Service列表
func (s *ClusterService) ListServicesHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s Service获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取Service数据
	// 如果提供了命名空间，则只获取该命名空间的Service，否则获取所有命名空间的Service
	var services *corev1.ServiceList
	if namespace != "" {
		// 获取指定命名空间的Service
		services, err = clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的Service
		services, err = clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list services: " + err.Error(),
		})
		return
	}

	// 转换为我们的Service模型
	var resultServices []models.Service
	for _, service := range services.Items {
		// 转换端口列表
		var ports []models.ServicePort
		for _, port := range service.Spec.Ports {
			targetPort := 0
			if port.TargetPort.IntValue() > 0 {
				targetPort = int(port.TargetPort.IntValue())
			}

			ports = append(ports, models.ServicePort{
				Name:       port.Name,
				Port:       int(port.Port),
				TargetPort: targetPort,
				Protocol:   string(port.Protocol),
				NodePort:   int(port.NodePort),
			})
		}

		// 获取外部IP，考虑所有可能的情况
		externalIP := ""
		// 1. 对于ExternalName类型的Service，EXTERNAL-IP就是externalName
		if service.Spec.Type == "ExternalName" && service.Spec.ExternalName != "" {
			externalIP = service.Spec.ExternalName
			// 2. 对于LoadBalancer类型的Service，从Ingress获取
		} else if service.Spec.Type == "LoadBalancer" && len(service.Status.LoadBalancer.Ingress) > 0 {
			// 检查是否有IP地址
			if service.Status.LoadBalancer.Ingress[0].IP != "" {
				externalIP = service.Status.LoadBalancer.Ingress[0].IP
				// 否则检查是否有主机名
			} else if service.Status.LoadBalancer.Ingress[0].Hostname != "" {
				externalIP = service.Status.LoadBalancer.Ingress[0].Hostname
			}
			// 3. 从Spec.ExternalIPs获取
		} else if len(service.Spec.ExternalIPs) > 0 {
			externalIP = service.Spec.ExternalIPs[0]
		}

		// 确保类型值不包含任何点
		serviceType := string(service.Spec.Type)
		// 移除所有点
		serviceType = strings.ReplaceAll(serviceType, ".", "")

		// 创建Service模型，使用Kubernetes对象的实际创建时间
		serviceModel := models.Service{
			ClusterID:       cluster.ID,
			Name:            service.Name,
			Namespace:       service.Namespace,
			Type:            serviceType,
			ClusterIP:       service.Spec.ClusterIP,
			ExternalIP:      externalIP,
			ExternalName:    service.Spec.ExternalName,
			Ports:           ports,
			Selector:        service.Spec.Selector,
			SessionAffinity: string(service.Spec.SessionAffinity),
		}

		// 显式设置CreatedAt字段为Kubernetes对象的创建时间
		serviceModel.CreatedAt = service.CreationTimestamp.Time

		resultServices = append(resultServices, serviceModel)
	}

	// 如果resultServices为空，初始化一个空切片，确保返回[]而不是null
	if resultServices == nil {
		resultServices = []models.Service{}
	}

	c.JSON(http.StatusOK, resultServices)
}

// ListConfigMapsHandler 获取集群ConfigMap列表
func (s *ClusterService) ListConfigMapsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s ConfigMap获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取ConfigMap数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的ConfigMap，否则获取所有命名空间的ConfigMap
	var configMaps *corev1.ConfigMapList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的ConfigMap
		configMaps, err = clientset.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的ConfigMap
		configMaps, err = clientset.CoreV1().ConfigMaps("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list configmaps: " + err.Error(),
		})
		return
	}

	// 转换为我们的ConfigMap模型
	var resultConfigMaps []models.ConfigMap
	for _, cm := range configMaps.Items {
		// 转换BinaryData从map[string][]byte到map[string]string
		binaryData := make(map[string]string)
		for key, value := range cm.BinaryData {
			binaryData[key] = string(value)
		}

		resultConfigMaps = append(resultConfigMaps, models.ConfigMap{
			ClusterID:   cluster.ID,
			Name:        cm.Name,
			Namespace:   cm.Namespace,
			Data:        cm.Data,
			BinaryData:  binaryData,
			Annotations: cm.Annotations,
			CreatedAt:   cm.CreationTimestamp.Time,
		})
	}

	// 如果resultConfigMaps为空，初始化一个空切片，确保返回[]而不是null
	if resultConfigMaps == nil {
		resultConfigMaps = []models.ConfigMap{}
	}

	c.JSON(http.StatusOK, resultConfigMaps)
}

// ListSecretsHandler 获取集群Secret列表
func (s *ClusterService) ListSecretsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s Secret获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取Secret数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的Secret，否则获取所有命名空间的Secret
	var secrets *corev1.SecretList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的Secret
		secrets, err = clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的Secret
		secrets, err = clientset.CoreV1().Secrets("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list secrets: " + err.Error(),
		})
		return
	}

	// 转换为我们的Secret模型
	var resultSecrets []models.Secret
	for _, secret := range secrets.Items {
		// 将map[string][]byte转换为map[string]string
		secretData := make(map[string]string)
		for key, value := range secret.Data {
			secretData[key] = string(value)
		}

		// K8s Secret没有单独的BinaryData字段，只有Data字段
		// 这里返回空的BinaryData，以匹配前端期望的结构
		binaryData := make(map[string]string)

		resultSecrets = append(resultSecrets, models.Secret{
			ClusterID:   cluster.ID,
			Name:        secret.Name,
			Namespace:   secret.Namespace,
			Type:        string(secret.Type),
			Data:        secretData,
			BinaryData:  binaryData,
			Annotations: secret.Annotations,
			CreatedAt:   secret.CreationTimestamp.Time,
		})
	}

	// 如果resultSecrets为空，初始化一个空切片，确保返回[]而不是null
	if resultSecrets == nil {
		resultSecrets = []models.Secret{}
	}

	c.JSON(http.StatusOK, resultSecrets)
}

// ListIngressesHandler 获取集群Ingress列表
func (s *ClusterService) ListIngressesHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s Ingress获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取Ingress数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的Ingress，否则获取所有命名空间的Ingress
	var ingresses *networkingv1.IngressList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的Ingress
		ingresses, err = clientset.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的Ingress
		ingresses, err = clientset.NetworkingV1().Ingresses("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list ingresses: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultIngresses []map[string]interface{}
	for _, ingress := range ingresses.Items {
		// 处理规则
		rules := make([]map[string]interface{}, 0)
		for _, rule := range ingress.Spec.Rules {
			httpPaths := make([]map[string]interface{}, 0)
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					// 处理PathType
					pathTypeStr := ""
					if path.PathType != nil {
						pathTypeStr = string(*path.PathType)
					}
					httpPaths = append(httpPaths, map[string]interface{}{
						"path":     path.Path,
						"pathType": pathTypeStr,
						"backend": map[string]interface{}{
							"service": map[string]interface{}{
								"name": path.Backend.Service.Name,
								"port": map[string]interface{}{
									"number": path.Backend.Service.Port.Number,
								},
							},
						},
					})
				}
			}
			rules = append(rules, map[string]interface{}{
				"host": rule.Host,
				"http": map[string]interface{}{
					"paths": httpPaths,
				},
			})
		}

		// 构建Ingress响应数据
		resultIngress := map[string]interface{}{
			"name":        ingress.Name,
			"namespace":   ingress.Namespace,
			"rules":       rules,
			"annotations": ingress.Annotations,
			"createdAt":   ingress.CreationTimestamp.Time,
		}
		resultIngresses = append(resultIngresses, resultIngress)
	}

	// 如果resultIngresses为空，初始化一个空切片，确保返回[]而不是null
	if resultIngresses == nil {
		resultIngresses = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultIngresses)
}

// ListStatefulSetsHandler 获取集群StatefulSet列表
func (s *ClusterService) ListStatefulSetsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s StatefulSet获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取StatefulSet数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的StatefulSet，否则获取所有命名空间的StatefulSet
	var statefulsets *appsv1.StatefulSetList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的StatefulSet
		statefulsets, err = clientset.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的StatefulSet
		statefulsets, err = clientset.AppsV1().StatefulSets("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list statefulsets: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultStatefulSets []map[string]interface{}
	for _, statefulset := range statefulsets.Items {
		// 获取第一个容器的镜像
		image := ""
		if len(statefulset.Spec.Template.Spec.Containers) > 0 {
			image = statefulset.Spec.Template.Spec.Containers[0].Image
		}

		// 获取StatefulSet状态
		status := "Unknown"
		if statefulset.Status.ReadyReplicas == statefulset.Status.Replicas {
			status = "Available"
		} else if statefulset.Status.UpdatedReplicas > 0 {
			status = "Progressing"
		} else {
			status = "Degraded"
		}

		// 处理Replicas指针类型
		replicas := 0
		if statefulset.Spec.Replicas != nil {
			replicas = int(*statefulset.Spec.Replicas)
		}

		// 构建StatefulSet响应数据
		resultStatefulSet := map[string]interface{}{
			"name":          statefulset.Name,
			"namespace":     statefulset.Namespace,
			"status":        status,
			"replicas":      replicas,
			"readyReplicas": int(statefulset.Status.ReadyReplicas),
			"image":         image,
			"createdAt":     statefulset.CreationTimestamp.Time,
		}
		resultStatefulSets = append(resultStatefulSets, resultStatefulSet)
	}

	// 如果resultStatefulSets为空，初始化一个空切片，确保返回[]而不是null
	if resultStatefulSets == nil {
		resultStatefulSets = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultStatefulSets)
}

// ListDaemonSetsHandler 获取集群DaemonSet列表
func (s *ClusterService) ListDaemonSetsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s DaemonSet获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取DaemonSet数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的DaemonSet，否则获取所有命名空间的DaemonSet
	var daemonsets *appsv1.DaemonSetList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的DaemonSet
		daemonsets, err = clientset.AppsV1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的DaemonSet
		daemonsets, err = clientset.AppsV1().DaemonSets("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list daemonsets: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultDaemonSets []map[string]interface{}
	for _, daemonset := range daemonsets.Items {
		// 获取第一个容器的镜像
		image := ""
		if len(daemonset.Spec.Template.Spec.Containers) > 0 {
			image = daemonset.Spec.Template.Spec.Containers[0].Image
		}

		// 获取DaemonSet状态
		status := "Unknown"
		if daemonset.Status.NumberReady == daemonset.Status.DesiredNumberScheduled {
			status = "Available"
		} else if daemonset.Status.UpdatedNumberScheduled > 0 {
			status = "Progressing"
		} else {
			status = "Degraded"
		}

		// 构建DaemonSet响应数据
		resultDaemonSet := map[string]interface{}{
			"name":      daemonset.Name,
			"namespace": daemonset.Namespace,
			"status":    status,
			"image":     image,
			"createdAt": daemonset.CreationTimestamp.Time,
		}
		resultDaemonSets = append(resultDaemonSets, resultDaemonSet)
	}

	// 如果resultDaemonSets为空，初始化一个空切片，确保返回[]而不是null
	if resultDaemonSets == nil {
		resultDaemonSets = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultDaemonSets)
}

// ListJobsHandler 获取集群Job列表
func (s *ClusterService) ListJobsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s Job获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取Job数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的Job，否则获取所有命名空间的Job
	var jobs *batchv1.JobList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的Job
		jobs, err = clientset.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的Job
		jobs, err = clientset.BatchV1().Jobs("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list jobs: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultJobs []map[string]interface{}
	for _, job := range jobs.Items {
		// 获取第一个容器的镜像
		image := ""
		if len(job.Spec.Template.Spec.Containers) > 0 {
			image = job.Spec.Template.Spec.Containers[0].Image
		}

		// 获取Job状态
		status := "Unknown"
		if job.Status.Succeeded > 0 {
			status = "Succeeded"
		} else if job.Status.Failed > 0 {
			status = "Failed"
		} else if job.Status.Active > 0 {
			status = "Running"
		} else {
			status = "Pending"
		}

		// 处理Completions指针类型
		completions := int32(1)
		if job.Spec.Completions != nil {
			completions = *job.Spec.Completions
		}

		// 构建Job响应数据
		resultJob := map[string]interface{}{
			"name":        job.Name,
			"namespace":   job.Namespace,
			"status":      status,
			"completions": int(completions),
			"succeeded":   int(job.Status.Succeeded),
			"image":       image,
			"createdAt":   job.CreationTimestamp.Time,
		}
		resultJobs = append(resultJobs, resultJob)
	}

	// 如果resultJobs为空，初始化一个空切片，确保返回[]而不是null
	if resultJobs == nil {
		resultJobs = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultJobs)
}

// ListCronJobsHandler 获取集群CronJob列表
func (s *ClusterService) ListCronJobsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s CronJob获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取CronJob数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的CronJob，否则获取所有命名空间的CronJob
	var cronjobs *batchv1.CronJobList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的CronJob
		cronjobs, err = clientset.BatchV1().CronJobs(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的CronJob
		cronjobs, err = clientset.BatchV1().CronJobs("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list cronjobs: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultCronJobs []map[string]interface{}
	for _, cronjob := range cronjobs.Items {
		// 获取第一个容器的镜像
		image := ""
		if len(cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers) > 0 {
			image = cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
		}

		// 获取CronJob状态
		status := "Unknown"
		if cronjob.Status.LastScheduleTime != nil {
			status = "Active"
		} else {
			status = "Inactive"
		}

		// 构建CronJob响应数据
		resultCronJob := map[string]interface{}{
			"name":         cronjob.Name,
			"namespace":    cronjob.Namespace,
			"status":       status,
			"schedule":     cronjob.Spec.Schedule,
			"lastSchedule": cronjob.Status.LastScheduleTime,
			"image":        image,
			"createdAt":    cronjob.CreationTimestamp.Time,
		}
		resultCronJobs = append(resultCronJobs, resultCronJob)
	}

	// 如果resultCronJobs为空，初始化一个空切片，确保返回[]而不是null
	if resultCronJobs == nil {
		resultCronJobs = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultCronJobs)
}

// testClusterConnection 测试集群连接并更新集群状态
func testClusterConnection(cluster *models.Cluster) {
	// 实现实际的K8s连接测试逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		cluster.Status = "error"
		cluster.Version = ""
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		cluster.Status = "error"
		cluster.Version = ""
		return
	}

	// 测试连接，获取服务器版本
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		cluster.Status = "disconnected"
		cluster.Version = ""
		return
	}

	// 连接成功，更新状态和版本
	cluster.Status = "connected"
	cluster.Version = version.String()
}

// TestConnectionHandler 测试集群连接
func (s *ClusterService) TestConnectionHandler(c *gin.Context) {
	id := c.Param("id")
	var cluster models.Cluster
	result := s.db.First(&cluster, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 测试集群连接并更新状态
	testClusterConnection(&cluster)
	s.db.Save(&cluster)

	if cluster.Status == "connected" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "连接测试成功",
			"version": cluster.Version,
		})
	} else if cluster.Status == "error" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Kubeconfig解析失败，请检查文件格式是否正确",
			"version": "",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "连接测试失败，请检查集群配置",
			"version": "",
		})
	}
}

// GetDashboardHandler 获取集群仪表盘数据
func (s *ClusterService) GetDashboardHandler(c *gin.Context) {
	clusterID := c.Param("id")
	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s仪表盘数据获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 初始化统计数据
	stats := struct {
		Pods struct {
			Total   int `json:"total"`
			Running int `json:"running"`
			Pending int `json:"pending"`
			Failed  int `json:"failed"`
		} `json:"pods"`
		Services struct {
			Total        int `json:"total"`
			ClusterIP    int `json:"clusterIP"`
			NodePort     int `json:"nodePort"`
			LoadBalancer int `json:"loadBalancer"`
		} `json:"services"`
		PVC struct {
			Total   int `json:"total"`
			Bound   int `json:"bound"`
			Pending int `json:"pending"`
			Failed  int `json:"failed"`
		} `json:"pvc"`
		PV struct {
			Total     int `json:"total"`
			Available int `json:"available"`
			Bound     int `json:"bound"`
			Failed    int `json:"failed"`
		} `json:"pv"`
		ConfigMaps struct {
			Total int `json:"total"`
		} `json:"configMaps"`
		Secrets struct {
			Total int `json:"total"`
		} `json:"secrets"`
		Ingresses struct {
			Total int `json:"total"`
		} `json:"ingresses"`
	}{}

	// 1. 获取Pod统计
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list pods: " + err.Error(),
		})
		return
	}

	stats.Pods.Total = len(pods.Items)
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			stats.Pods.Running++
		case corev1.PodPending:
			stats.Pods.Pending++
		case corev1.PodFailed:
			stats.Pods.Failed++
		}
	}

	// 2. 获取Service统计
	services, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list services: " + err.Error(),
		})
		return
	}

	stats.Services.Total = len(services.Items)
	for _, service := range services.Items {
		switch service.Spec.Type {
		case corev1.ServiceTypeClusterIP:
			stats.Services.ClusterIP++
		case corev1.ServiceTypeNodePort:
			stats.Services.NodePort++
		case corev1.ServiceTypeLoadBalancer:
			stats.Services.LoadBalancer++
		}
	}

	// 3. 获取PVC统计
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list pvcs: " + err.Error(),
		})
		return
	}

	stats.PVC.Total = len(pvcs.Items)
	for _, pvc := range pvcs.Items {
		switch pvc.Status.Phase {
		case corev1.ClaimBound:
			stats.PVC.Bound++
		case corev1.ClaimPending:
			stats.PVC.Pending++
			// PVC没有Failed状态，这里简化处理
		}
	}

	// 4. 获取PV统计
	pvs, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list pvs: " + err.Error(),
		})
		return
	}

	stats.PV.Total = len(pvs.Items)
	for _, pv := range pvs.Items {
		switch pv.Status.Phase {
		case corev1.VolumeAvailable:
			stats.PV.Available++
		case corev1.VolumeBound:
			stats.PV.Bound++
			// PV没有Failed状态，这里简化处理
		}
	}

	// 5. 获取ConfigMap统计
	configMaps, err := clientset.CoreV1().ConfigMaps("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list configmaps: " + err.Error(),
		})
		return
	}
	stats.ConfigMaps.Total = len(configMaps.Items)

	// 6. 获取Secret统计
	secrets, err := clientset.CoreV1().Secrets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list secrets: " + err.Error(),
		})
		return
	}
	stats.Secrets.Total = len(secrets.Items)

	// 7. 获取Ingress统计
	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list ingresses: " + err.Error(),
		})
		return
	}
	stats.Ingresses.Total = len(ingresses.Items)

	// 8. 获取资源使用率（这里简化处理，实际应该从metrics-server获取）
	resourceUsage := struct {
		CPU     int `json:"cpu"`
		Memory  int `json:"memory"`
		Storage int `json:"storage"`
	}{
		CPU:     45, // 模拟数据
		Memory:  68, // 模拟数据
		Storage: 32, // 模拟数据
	}

	// 6. 获取命名空间分布数据
	// 初始化命名空间计数
	namespacePodCount := make(map[string]int)
	namespaceServiceCount := make(map[string]int)
	namespaceDeploymentCount := make(map[string]int)

	// 统计Pod分布
	for _, pod := range pods.Items {
		namespacePodCount[pod.Namespace]++
	}

	// 统计Service分布
	for _, service := range services.Items {
		namespaceServiceCount[service.Namespace]++
	}

	// 统计Deployment分布
	deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list deployments: " + err.Error(),
		})
		return
	}

	for _, deployment := range deployments.Items {
		namespaceDeploymentCount[deployment.Namespace]++
	}

	// 转换为图表数据格式
	var namespacePodData []map[string]interface{}
	for ns, count := range namespacePodCount {
		namespacePodData = append(namespacePodData, map[string]interface{}{
			"name":  ns,
			"value": count,
		})
	}

	var namespaceServiceData []map[string]interface{}
	for ns, count := range namespaceServiceCount {
		namespaceServiceData = append(namespaceServiceData, map[string]interface{}{
			"name":  ns,
			"value": count,
		})
	}

	var namespaceDeploymentData []map[string]interface{}
	for ns, count := range namespaceDeploymentCount {
		namespaceDeploymentData = append(namespaceDeploymentData, map[string]interface{}{
			"name":  ns,
			"value": count,
		})
	}

	// 7. 最近活动（模拟数据）
	recentActivities := []map[string]interface{}{
		{
			"time":        "2026-02-06 14:30:00",
			"type":        "info",
			"title":       "Pod创建",
			"description": "default命名空间中创建了新的Pod: nginx-7f489b75c9-8xzv9",
		},
		{
			"time":        "2026-02-06 14:25:00",
			"type":        "success",
			"title":       "Deployment更新",
			"description": "kube-system命名空间中的coredns Deployment已成功更新",
		},
		{
			"time":        "2026-02-06 14:20:00",
			"type":        "warning",
			"title":       "节点状态变化",
			"description": "节点worker-1状态变为NotReady",
		},
		{
			"time":        "2026-02-06 14:15:00",
			"type":        "info",
			"title":       "Service创建",
			"description": "default命名空间中创建了新的Service: nginx-service",
		},
		{
			"time":        "2026-02-06 14:10:00",
			"type":        "success",
			"title":       "集群连接成功",
			"description": "已成功连接到集群: " + cluster.Name,
		},
	}

	// 构建响应数据
	dashboardData := gin.H{
		"stats":                   stats,
		"resourceUsage":           resourceUsage,
		"namespacePodData":        namespacePodData,
		"namespaceServiceData":    namespaceServiceData,
		"namespaceDeploymentData": namespaceDeploymentData,
		"recentActivities":        recentActivities,
	}

	c.JSON(http.StatusOK, dashboardData)
}

// GetPodLogsHandler 获取Pod日志
func (s *ClusterService) GetPodLogsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	podName := c.Param("name")

	// 获取请求参数
	namespace := c.Query("namespace")
	container := c.Query("container")
	linesStr := c.Query("lines")

	// 解析日志行数，默认为100行
	lines := int64(100)
	if linesStr != "" {
		linesInt, err := strconv.ParseInt(linesStr, 10, 64)
		if err == nil && linesInt > 0 {
			lines = linesInt
		}
	}

	// 获取集群信息
	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 创建K8s配置
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 获取Pod日志
	logOptions := &corev1.PodLogOptions{
		Container: container,
		TailLines: &lines,
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	logStream, err := req.Stream(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pod logs: " + err.Error(),
		})
		return
	}
	defer logStream.Close()

	// 读取日志内容
	logData, err := io.ReadAll(logStream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read pod logs: " + err.Error(),
		})
		return
	}

	// 返回日志内容，设置Content-Type为text/plain
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, string(logData))
}

// ExecPodHandler 处理Pod终端WebSocket连接
func (s *ClusterService) ExecPodHandler(c *gin.Context) {
	clusterID := c.Param("id")
	podName := c.Param("name")

	// 获取集群信息
	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 创建K8s配置
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 升级HTTP连接为WebSocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 直接返回错误，不使用JSON，因为连接已经被升级
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte("Failed to upgrade to WebSocket: " + err.Error()))
		return
	}
	defer ws.Close()

	// 读取初始消息获取容器信息
	var initMsg struct {
		Namespace string `json:"namespace"`
		Container string `json:"container"`
	}

	// 设置读取超时
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))

	msgType, msg, err := ws.ReadMessage()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Failed to read container info: "+err.Error()))
		return
	}

	if msgType != websocket.TextMessage {
		ws.WriteMessage(websocket.TextMessage, []byte("Invalid message type"))
		return
	}

	if err := json.Unmarshal(msg, &initMsg); err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Failed to parse container info: "+err.Error()))
		return
	}

	// 重置读取超时
	ws.SetReadDeadline(time.Time{})

	// 发送调试信息
	ws.WriteMessage(websocket.TextMessage, []byte("正在执行命令...\n"))

	// 执行Pod命令
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(initMsg.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: initMsg.Container,
			Command:   []string{"/bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	// 发送调试信息
	ws.WriteMessage(websocket.TextMessage, []byte("创建executor...\n"))

	// 创建executor
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Failed to create executor: "+err.Error()))
		return
	}

	// 发送调试信息
	ws.WriteMessage(websocket.TextMessage, []byte("executor创建成功，准备执行命令...\n"))

	// 创建管道
	stdin, stdinWriter := io.Pipe()
	stdout := &WebSocketWriter{ws: ws}
	stderr := &WebSocketWriter{ws: ws}

	// 运行命令
	go func() {
		defer stdin.Close()
		defer stdinWriter.Close()

		// 读取WebSocket消息并写入stdin
		for {
			msgType, msg, err := ws.ReadMessage()
			if err != nil {
				break
			}

			if msgType != websocket.TextMessage {
				continue
			}

			if _, err := stdinWriter.Write(msg); err != nil {
				break
			}
		}
	}()

	// 执行命令
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    true,
	})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Command execution failed: "+err.Error()))
	}
}

// WebSocketWriter 是一个实现了io.Writer接口的类型，用于将输出写入WebSocket
type WebSocketWriter struct {
	ws *websocket.Conn
}

// Write 将数据写入WebSocket
func (w *WebSocketWriter) Write(p []byte) (n int, err error) {
	err = w.ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// upgrader 是WebSocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
	Subprotocols: []string{"terminal"},
}

// ListReplicaSetsHandler 获取集群ReplicaSet列表
func (s *ClusterService) ListReplicaSetsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")
	deployment := c.Query("deployment")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s ReplicaSet获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取ReplicaSet数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的ReplicaSet，否则获取所有命名空间的ReplicaSet
	var replicasets *appsv1.ReplicaSetList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的ReplicaSet
		replicasets, err = clientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的ReplicaSet
		replicasets, err = clientset.AppsV1().ReplicaSets("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list replicasets: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultReplicaSets []map[string]interface{}
	for _, rs := range replicasets.Items {
		// 如果指定了Deployment，只返回关联的ReplicaSet
		if deployment != "" {
			// 检查ReplicaSet的所有者是否是指定的Deployment
			found := false
			for _, owner := range rs.OwnerReferences {
				if owner.Kind == "Deployment" && owner.Name == deployment {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// 获取第一个容器的镜像
		image := ""
		if len(rs.Spec.Template.Spec.Containers) > 0 {
			image = rs.Spec.Template.Spec.Containers[0].Image
		}

		// 处理Replicas指针类型
		replicas := 0
		if rs.Spec.Replicas != nil {
			replicas = int(*rs.Spec.Replicas)
		}

		// 构建ReplicaSet响应数据
		resultRS := map[string]interface{}{
			"Name":          rs.Name,
			"Namespace":     rs.Namespace,
			"Replicas":      replicas,
			"ReadyReplicas": int(rs.Status.ReadyReplicas),
			"Image":         image,
			"CreatedAt":     rs.CreationTimestamp.Time,
		}
		resultReplicaSets = append(resultReplicaSets, resultRS)
	}

	// 如果resultReplicaSets为空，初始化一个空切片，确保返回[]而不是null
	if resultReplicaSets == nil {
		resultReplicaSets = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultReplicaSets)
}

// CreateDeploymentHandler 创建Deployment
func (s *ClusterService) CreateDeploymentHandler(c *gin.Context) {
	clusterID := c.Param("id")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 检查kubectl命令是否存在
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kubectl command not found"})
		return
	}

	// 创建临时kubeconfig文件（如果需要）
	var kubeconfigPath string
	var tempFile *os.File

	if cluster.Kubeconfig != "" {
		// 创建临时kubeconfig文件
		tempFile, err = os.CreateTemp("", "kubeconfig-*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary kubeconfig file"})
			return
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		// 写入kubeconfig内容
		if _, err := tempFile.WriteString(cluster.Kubeconfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write kubeconfig file"})
			return
		}

		// 确保文件内容已写入
		if err := tempFile.Sync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync kubeconfig file"})
			return
		}

		kubeconfigPath = tempFile.Name()
	}

	// 解析请求体
	var req struct {
		Name          string `json:"name,omitempty"`
		Namespace     string `json:"namespace,omitempty"`
		Replicas      int    `json:"replicas,omitempty"`
		Image         string `json:"image,omitempty"`
		ContainerName string `json:"containerName,omitempty"`
		Ports         []struct {
			ContainerPort int    `json:"containerPort"`
			Protocol      string `json:"protocol"`
		} `json:"ports,omitempty"`
		Labels      map[string]string `json:"labels,omitempty"`
		Annotations map[string]string `json:"annotations,omitempty"`
		Resources   struct {
			Requests struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"requests"`
			Limits struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"limits"`
		} `json:"resources,omitempty"`
		YAML string `json:"yaml,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	var cmd *exec.Cmd
	var output []byte

	if req.YAML != "" {
		// YAML模式：直接应用YAML
		// 创建临时YAML文件
		yamlFile, err := os.CreateTemp("", "deployment-*.yaml")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary YAML file"})
			return
		}
		defer func() {
			yamlFile.Close()
			os.Remove(yamlFile.Name())
		}()

		// 写入YAML内容
		if _, err := yamlFile.WriteString(req.YAML); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write YAML file"})
			return
		}

		// 确保文件内容已写入
		if err := yamlFile.Sync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync YAML file"})
			return
		}

		// 构建kubectl命令
		cmd = exec.Command(kubectlPath, "apply", "-f", yamlFile.Name())
	} else {
		// 表单模式：生成YAML并应用
		// 生成YAML内容
		yamlContent := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + req.Name + `
  namespace: ` + req.Namespace + `
`

		// 添加标签
		if len(req.Labels) > 0 {
			yamlContent += `  labels:
`
			for key, value := range req.Labels {
				yamlContent += `    ` + key + `: ` + value + `
`
			}
		}

		// 添加注解
		if len(req.Annotations) > 0 {
			yamlContent += `  annotations:
`
			for key, value := range req.Annotations {
				yamlContent += `    ` + key + `: ` + value + `
`
			}
		}

		// 添加spec部分
		yamlContent += `spec:
  replicas: ` + strconv.Itoa(req.Replicas) + `
  selector:
    matchLabels:
`

		// 添加selector标签
		if len(req.Labels) > 0 {
			for key, value := range req.Labels {
				yamlContent += `      ` + key + `: ` + value + `
`
			}
		} else {
			// 默认标签
			yamlContent += `      app: ` + req.Name + `
`
		}

		// 添加template部分
		yamlContent += `  template:
    metadata:
      labels:
`

		// 添加template标签
		if len(req.Labels) > 0 {
			for key, value := range req.Labels {
				yamlContent += `        ` + key + `: ` + value + `
`
			}
		} else {
			// 默认标签
			yamlContent += `        app: ` + req.Name + `
`
		}

		// 添加spec.containers部分
		yamlContent += `    spec:
      containers:
      - name: ` + req.ContainerName + `
        image: ` + req.Image + `
`

		// 添加端口
		if len(req.Ports) > 0 {
			yamlContent += `        ports:
`
			for _, port := range req.Ports {
				yamlContent += `        - containerPort: ` + strconv.Itoa(port.ContainerPort) + `
          protocol: ` + port.Protocol + `
`
			}
		}

		// 添加资源配置
		yamlContent += `        resources:
          requests:
            cpu: ` + req.Resources.Requests.CPU + `
            memory: ` + req.Resources.Requests.Memory + `
          limits:
            cpu: ` + req.Resources.Limits.CPU + `
            memory: ` + req.Resources.Limits.Memory + `
`

		// 创建临时YAML文件
		yamlFile, err := os.CreateTemp("", "deployment-*.yaml")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary YAML file"})
			return
		}
		defer func() {
			yamlFile.Close()
			os.Remove(yamlFile.Name())
		}()

		// 写入YAML内容
		if _, err := yamlFile.WriteString(yamlContent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write YAML file"})
			return
		}

		// 确保文件内容已写入
		if err := yamlFile.Sync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync YAML file"})
			return
		}

		// 构建kubectl命令
		cmd = exec.Command(kubectlPath, "apply", "-f", yamlFile.Name())
	}

	// 设置KUBECONFIG环境变量
	if kubeconfigPath != "" {
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	}

	// 执行命令
	output, err = cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create deployment: " + err.Error(),
			"output": string(output),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deployment created successfully",
		"output":  string(output),
	})
}

// PrecheckDeploymentHandler 预检Deployment YAML
func (s *ClusterService) PrecheckDeploymentHandler(c *gin.Context) {
	clusterID := c.Param("id")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 检查kubectl命令是否存在
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kubectl command not found"})
		return
	}

	// 创建临时kubeconfig文件（如果需要）
	var kubeconfigPath string
	var tempFile *os.File

	if cluster.Kubeconfig != "" {
		// 创建临时kubeconfig文件
		tempFile, err = os.CreateTemp("", "kubeconfig-*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary kubeconfig file"})
			return
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		// 写入kubeconfig内容
		if _, err := tempFile.WriteString(cluster.Kubeconfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write kubeconfig file"})
			return
		}

		// 确保文件内容已写入
		if err := tempFile.Sync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync kubeconfig file"})
			return
		}

		kubeconfigPath = tempFile.Name()
	}

	// 解析请求体
	var req struct {
		YAML      string `json:"yaml" binding:"required"`
		Namespace string `json:"namespace,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 创建临时YAML文件
	yamlFile, err := os.CreateTemp("", "deployment-*.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary YAML file"})
		return
	}
	defer func() {
		yamlFile.Close()
		os.Remove(yamlFile.Name())
	}()

	// 写入YAML内容
	if _, err := yamlFile.WriteString(req.YAML); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write YAML file"})
		return
	}

	// 确保文件内容已写入
	if err := yamlFile.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync YAML file"})
		return
	}

	// 构建kubectl命令，使用dry-run验证YAML
	cmd := exec.Command(kubectlPath, "apply", "-f", yamlFile.Name(), "--dry-run=client")

	// 设置KUBECONFIG环境变量
	if kubeconfigPath != "" {
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "YAML 验证失败了！" + err.Error(),
			"output": string(output),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "YAML 验证通过了！",
		"output":  string(output),
	})
}

// ListPVCsHandler 获取集群PVC列表
func (s *ClusterService) ListPVCsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s PVC获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取PVC数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的PVC，否则获取所有命名空间的PVC
	var pvcs *corev1.PersistentVolumeClaimList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的PVC
		pvcs, err = clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的PVC
		pvcs, err = clientset.CoreV1().PersistentVolumeClaims("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list pvcs: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultPVCs []map[string]interface{}
	for _, pvc := range pvcs.Items {
		// 获取容量信息
		capacity := ""
		if pvc.Status.Capacity != nil {
			if storageCapacity, ok := pvc.Status.Capacity["storage"]; ok {
				capacity = storageCapacity.String()
			}
		}

		// 获取访问模式
		accessModes := make([]string, 0)
		for _, mode := range pvc.Spec.AccessModes {
			accessModes = append(accessModes, string(mode))
		}

		// 获取存储类名称
		storageClassName := ""
		if pvc.Spec.StorageClassName != nil {
			storageClassName = *pvc.Spec.StorageClassName
		}

		// 构建PVC响应数据
		resultPVC := map[string]interface{}{
			"name":             pvc.Name,
			"namespace":        pvc.Namespace,
			"status":           string(pvc.Status.Phase),
			"volumeName":       pvc.Spec.VolumeName,
			"capacity":         capacity,
			"accessModes":      accessModes,
			"storageClassName": storageClassName,
			"createdAt":        pvc.CreationTimestamp.Time,
		}
		resultPVCs = append(resultPVCs, resultPVC)
	}

	// 如果resultPVCs为空，初始化一个空切片，确保返回[]而不是null
	if resultPVCs == nil {
		resultPVCs = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultPVCs)
}

// ListPVsHandler 获取集群PV列表
func (s *ClusterService) ListPVsHandler(c *gin.Context) {
	clusterID := c.Param("id")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s PV获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取PV数据
	// PV是集群级别的资源，不需要指定命名空间
	pvs, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list pvs: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultPVs []map[string]interface{}
	for _, pv := range pvs.Items {
		// 获取容量信息
		capacity := ""
		if pv.Spec.Capacity != nil {
			if storageCapacity, ok := pv.Spec.Capacity["storage"]; ok {
				capacity = storageCapacity.String()
			}
		}

		// 获取访问模式
		accessModes := make([]string, 0)
		for _, mode := range pv.Spec.AccessModes {
			accessModes = append(accessModes, string(mode))
		}

		// 获取存储类名称
		storageClassName := pv.Spec.StorageClassName

		// 获取回收策略
		reclaimPolicy := string(pv.Spec.PersistentVolumeReclaimPolicy)

		// 获取绑定的PVC信息
		var claimRef map[string]interface{}
		if pv.Spec.ClaimRef != nil {
			claimRef = map[string]interface{}{
				"namespace": pv.Spec.ClaimRef.Namespace,
				"name":      pv.Spec.ClaimRef.Name,
			}
		}

		// 获取PV状态
		status := string(pv.Status.Phase)

		// 构建PV响应数据
		resultPV := map[string]interface{}{
			"name":             pv.Name,
			"status":           status,
			"claimRef":         claimRef,
			"storageClassName": storageClassName,
			"capacity":         capacity,
			"accessModes":      accessModes,
			"reclaimPolicy":    reclaimPolicy,
			"labels":           pv.Labels,
			"annotations":      pv.Annotations,
			"createdAt":        pv.CreationTimestamp.Time,
		}
		resultPVs = append(resultPVs, resultPV)
	}

	// 如果resultPVs为空，初始化一个空切片，确保返回[]而不是null
	if resultPVs == nil {
		resultPVs = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultPVs)
}

// ListStorageClassesHandler 获取集群存储类列表
func (s *ClusterService) ListStorageClassesHandler(c *gin.Context) {
	clusterID := c.Param("id")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s存储类获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取存储类数据
	// 存储类是集群级别的资源，不需要指定命名空间
	storageClasses, err := clientset.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list storageclasses: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultStorageClasses []map[string]interface{}
	for _, sc := range storageClasses.Items {
		// 获取回收策略
		reclaimPolicy := ""
		if sc.ReclaimPolicy != nil {
			reclaimPolicy = string(*sc.ReclaimPolicy)
		}

		// 获取卷绑定模式
		volumeBindingMode := string(*sc.VolumeBindingMode)

		// 获取允许卷扩展
		allowVolumeExpansion := false
		if sc.AllowVolumeExpansion != nil {
			allowVolumeExpansion = *sc.AllowVolumeExpansion
		}

		// 检查是否为默认存储类
		isDefault := false
		if sc.Annotations != nil {
			if val, ok := sc.Annotations["storageclass.kubernetes.io/is-default-class"]; ok && val == "true" {
				isDefault = true
			}
		}

		// 构建存储类响应数据
		resultSC := map[string]interface{}{
			"name":                 sc.Name,
			"provisioner":          sc.Provisioner,
			"reclaimPolicy":        reclaimPolicy,
			"volumeBindingMode":    volumeBindingMode,
			"allowVolumeExpansion": allowVolumeExpansion,
			"isDefault":            isDefault,
			"labels":               sc.Labels,
			"annotations":          sc.Annotations,
			"createdAt":            sc.CreationTimestamp.Time.Format(time.RFC3339),
		}
		resultStorageClasses = append(resultStorageClasses, resultSC)
	}

	// 如果resultStorageClasses为空，初始化一个空切片，确保返回[]而不是null
	if resultStorageClasses == nil {
		resultStorageClasses = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, resultStorageClasses)
}

// DeleteStorageClassHandler 删除集群存储类
func (s *ClusterService) DeleteStorageClassHandler(c *gin.Context) {
	clusterID := c.Param("id")
	name := c.Param("name")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s存储类删除逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 删除存储类
	err = clientset.StorageV1().StorageClasses().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete storageclass: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "StorageClass deleted successfully",
		"success": true,
	})
}

// GetStorageClassYamlHandler 获取集群存储类YAML内容
func (s *ClusterService) GetStorageClassYamlHandler(c *gin.Context) {
	clusterID := c.Param("id")
	name := c.Param("name")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s存储类YAML获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 获取存储类对象
	storageClass, err := clientset.StorageV1().StorageClasses().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get storageclass: " + err.Error(),
		})
		return
	}

	// 返回存储类的主要信息，因为直接使用yaml.v2.Marshal复杂的K8s对象会有问题
	c.JSON(http.StatusOK, gin.H{
		"name":                 storageClass.Name,
		"provisioner":          storageClass.Provisioner,
		"reclaimPolicy":        storageClass.ReclaimPolicy,
		"volumeBindingMode":    storageClass.VolumeBindingMode,
		"allowVolumeExpansion": storageClass.AllowVolumeExpansion,
		"labels":               storageClass.Labels,
		"annotations":          storageClass.Annotations,
		"createdAt":            storageClass.CreationTimestamp,
	})
}

// ListEventsHandler 获取集群事件列表
func (s *ClusterService) ListEventsHandler(c *gin.Context) {
	clusterID := c.Param("id")
	// 获取请求参数
	namespace := c.Query("namespace")
	eventType := c.Query("type")
	searchKey := c.Query("search")

	var cluster models.Cluster
	result := s.db.First(&cluster, clusterID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// 实现实际的K8s事件获取逻辑
	var config *rest.Config
	var err error

	if cluster.Kubeconfig != "" {
		// 使用kubeconfig方式连接
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Kubeconfig))
	} else {
		// 使用Token方式连接
		config = &rest.Config{
			Host:        cluster.APIURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 允许不安全的连接，生产环境应该设置为false
			},
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s config: " + err.Error(),
		})
		return
	}

	// 创建K8s客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// 从K8s集群获取事件数据
	// 如果提供了命名空间且不为空，则只获取该命名空间的事件，否则获取所有命名空间的事件
	var events *corev1.EventList
	if strings.TrimSpace(namespace) != "" {
		// 获取指定命名空间的事件
		events, err = clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	} else {
		// 获取所有命名空间的事件
		events, err = clientset.CoreV1().Events("").List(context.Background(), metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list events: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的数据格式
	var resultEvents []map[string]interface{}
	for _, event := range events.Items {
		// 构建事件响应数据
		eventData := map[string]interface{}{
			"name":           event.Name,
			"namespace":      event.Namespace,
			"type":           event.Type,
			"reason":         event.Reason,
			"message":        event.Message,
			"count":          event.Count,
			"firstTimestamp": event.FirstTimestamp.Time.Format(time.RFC3339),
			"lastTimestamp":  event.LastTimestamp.Time.Format(time.RFC3339),
			"source": map[string]interface{}{
				"component": event.Source.Component,
				"host":      event.Source.Host,
			},
			"involvedObject": map[string]interface{}{
				"kind":            event.InvolvedObject.Kind,
				"namespace":       event.InvolvedObject.Namespace,
				"name":            event.InvolvedObject.Name,
				"uid":             string(event.InvolvedObject.UID),
				"apiVersion":      event.InvolvedObject.APIVersion,
				"resourceVersion": event.InvolvedObject.ResourceVersion,
				"fieldPath":       event.InvolvedObject.FieldPath,
			},
		}
		resultEvents = append(resultEvents, eventData)
	}

	// 处理筛选
	filteredEvents := resultEvents

	// 按事件类型筛选
	if eventType != "" {
		var temp []map[string]interface{}
		for _, event := range filteredEvents {
			if eventType == event["type"] {
				temp = append(temp, event)
			}
		}
		filteredEvents = temp
	}

	// 按搜索关键词筛选
	if searchKey != "" {
		var temp []map[string]interface{}
		for _, event := range filteredEvents {
			// 搜索事件信息、原因、资源名称等字段
			match := false
			if reason, ok := event["reason"].(string); ok && strings.Contains(reason, searchKey) {
				match = true
			}
			if message, ok := event["message"].(string); ok && strings.Contains(message, searchKey) {
				match = true
			}
			if involvedObj, ok := event["involvedObject"].(map[string]interface{}); ok {
				if name, ok := involvedObj["name"].(string); ok && strings.Contains(name, searchKey) {
					match = true
				}
				if kind, ok := involvedObj["kind"].(string); ok && strings.Contains(kind, searchKey) {
					match = true
				}
			}
			if match {
				temp = append(temp, event)
			}
		}
		filteredEvents = temp
	}

	// 处理分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 计算总条数
	total := len(filteredEvents)

	// 计算分页起始位置
	start := (page - 1) * pageSize
	end := start + pageSize

	// 确保起始位置不小于0
	if start < 0 {
		start = 0
	}

	// 确保结束位置不大于总条数
	if end > total {
		end = total
	}

	// 截取分页数据
	var paginatedEvents []map[string]interface{}
	if start < total {
		paginatedEvents = filteredEvents[start:end]
	} else {
		paginatedEvents = []map[string]interface{}{}
	}

	// 返回分页数据
	c.JSON(http.StatusOK, gin.H{
		"items":    paginatedEvents,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
