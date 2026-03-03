package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username   string `gorm:"uniqueIndex;size:50"`
	Password   string `gorm:"size:255"`
	Email      string `gorm:"size:100"`
	Name       string `gorm:"size:50"`
	Phone      string `gorm:"size:20"`
	Department string `gorm:"size:50"`
	Position   string `gorm:"size:50"`
	Remark     string `gorm:"type:text"` // 备注
	Avatar     string `gorm:"type:text"` // 头像base64数据
	Role       string `gorm:"size:20;default:'user'"`
	Status     string `gorm:"size:20;default:'inactive'"` // active或inactive
}

// Alert 告警模型
type Alert struct {
	gorm.Model
	AlertName    string `gorm:"size:100"`
	Status       string `gorm:"size:20"`
	Severity     string `gorm:"size:20"`
	Labels       string `gorm:"type:text"`
	Annotations  string `gorm:"type:text"`
	StartsAt     time.Time
	EndsAt       time.Time
	GeneratorURL string `gorm:"size:255"`
	Fingerprint  string `gorm:"uniqueIndex;size:64"`
	ResolvedAt   time.Time
	ResolvedBy   string `gorm:"size:50"`
	AlertMediaID uint
	AlertMedia   AlertMedia
}

// AlertMedia 告警介质模型
type AlertMedia struct {
	gorm.Model
	Name       string `gorm:"uniqueIndex;size:50"`
	Type       string `gorm:"size:20"` // 卡片、自定义模板等
	Config     string `gorm:"type:text"`
	TemplateID uint
	Template   AlertTemplate
	IsEnabled  bool `gorm:"default:true"`
}

// AlertTemplate 告警模板模型
type AlertTemplate struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex;size:50"`
	Content   string `gorm:"type:text"`
	Type      string `gorm:"size:20"` // 邮件、短信、Webhook等
	IsDefault bool   `gorm:"default:false"`
}

// Silence 静默规则模型
type Silence struct {
	gorm.Model
	Matchers  string `gorm:"type:text"` // JSON格式的匹配规则
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedBy string `gorm:"size:50"`
	Comment   string `gorm:"size:255"`
	Status    string `gorm:"size:20"`
}

// InhibitionRule 抑制规则模型
type InhibitionRule struct {
	gorm.Model
	SourceMatch string `gorm:"type:text"`
	TargetMatch string `gorm:"type:text"`
	Equal       string `gorm:"type:text"`
	CreatedBy   string `gorm:"size:50"`
	Comment     string `gorm:"size:255"`
	IsEnabled   bool   `gorm:"default:true"`
}

// Host 主机模型
type Host struct {
	gorm.Model
	Name     string `gorm:"size:100"`
	IP       string `gorm:"size:50;index"`
	Port     int    `gorm:"default:22"`
	Username string `gorm:"size:50"`
	Password string `gorm:"size:255"`
	Remark   string `gorm:"type:text"`
	IsActive bool   `gorm:"default:true"`
	Status   string `gorm:"size:20;default:'unknown'"` // running, stopped, connecting, error
}

// KVMHost KVM主机模型
type KVMHost struct {
	gorm.Model
	Name           string `gorm:"size:100"`
	IP             string `gorm:"size:50;index"`
	Port           int    `gorm:"default:22"`
	ConnectionType string `gorm:"size:20;default:'ssh'"` // ssh, virsh
	Username       string `gorm:"size:50"`
	Password       string `gorm:"size:255"`
	Description    string `gorm:"type:text"`
	IsActive       bool   `gorm:"default:true"`
	Status         string `gorm:"size:20;default:'unknown'"` // running, stopped, connecting, error
}

// TableName 设置KVMHost模型的表名
func (KVMHost) TableName() string {
	return "kvm_hosts"
}

// Cluster K8s集群模型
type Cluster struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name           string         `gorm:"uniqueIndex;size:100" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	APIURL         string         `gorm:"size:255" json:"apiURL"`
	Token          string         `gorm:"size:255" json:"token"`
	Kubeconfig     string         `gorm:"type:text" json:"kubeconfig"`
	ConnectionType string         `gorm:"size:20;default:'unknown'" json:"connectionType"` // token, kubeconfig
	IsActive       *bool          `gorm:"default:true" json:"isActive"`
	Status         string         `gorm:"size:20;default:'unknown'" json:"status"` // connected, disconnected, error
	Version        string         `gorm:"size:20" json:"version"`
	Nodes          []Node         `json:"-"`
	Pods           []Pod          `json:"-"`
	Deployments    []Deployment   `json:"-"`
}

// Node K8s节点模型
type Node struct {
	gorm.Model
	ClusterID        uint    `gorm:"index" json:"clusterID"`
	Name             string  `gorm:"size:100" json:"name"`
	IP               string  `gorm:"size:50" json:"ip"`
	Role             string  `gorm:"size:20" json:"role"`   // master, worker
	Status           string  `gorm:"size:20" json:"status"` // Ready, NotReady, Unknown
	OS               string  `gorm:"size:50" json:"os"`
	KernelVersion    string  `gorm:"size:50" json:"kernelVersion"`
	ContainerRuntime string  `gorm:"size:50" json:"containerRuntime"`
	CPU              int     `gorm:"default:0" json:"cpu"`
	Memory           int64   `gorm:"default:0" json:"memory"` // 单位: MB
	Pods             int     `gorm:"default:0" json:"pods"`
	CPUUsage         float64 `gorm:"default:0" json:"cpuUsage"`    // CPU使用率，百分比
	MemoryUsage      float64 `gorm:"default:0" json:"memoryUsage"` // 内存使用率，百分比
}

// PodContainer Pod容器模型
type PodContainer struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
	Ready  bool   `json:"ready"`
}

// Pod K8s Pod模型
type Pod struct {
	gorm.Model
	ClusterID    uint           `gorm:"index" json:"clusterID"`
	Name         string         `gorm:"size:100" json:"name"`
	Namespace    string         `gorm:"size:50" json:"namespace"`
	Status       string         `gorm:"size:20" json:"status"` // Running, Pending, Succeeded, Failed, Unknown
	NodeName     string         `gorm:"size:100" json:"nodeName"`
	PodIP        string         `gorm:"size:50" json:"podIP"`
	HostIP       string         `gorm:"size:50" json:"hostIP"`
	CreatedAt    time.Time      `json:"createdAt"`
	RestartCount int            `gorm:"default:0" json:"restartCount"`
	Containers   []PodContainer `gorm:"-" json:"containers"` // 不存储到数据库，只用于API响应
}

// Deployment K8s Deployment模型
type Deployment struct {
	gorm.Model
	ClusterID     uint      `gorm:"index" json:"clusterID"`
	Name          string    `gorm:"size:100" json:"name"`
	Namespace     string    `gorm:"size:50" json:"namespace"`
	Replicas      int       `gorm:"default:0" json:"replicas"`
	ReadyReplicas int       `gorm:"default:0" json:"readyReplicas"`
	Status        string    `gorm:"size:20" json:"status"` // Available, Progressing, Degraded
	Image         string    `gorm:"size:255" json:"image"`
	CPU           int       `gorm:"default:0" json:"cpu"`    // 单位: millicpu
	Memory        int64     `gorm:"default:0" json:"memory"` // 单位: MB
	CreatedAt     time.Time `json:"createdAt"`
}

// ServicePort K8s Service端口模型
type ServicePort struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	Protocol   string `json:"protocol"`
	NodePort   int    `json:"nodePort"`
}

// Service K8s Service模型
type Service struct {
	gorm.Model
	ClusterID       uint              `gorm:"index" json:"clusterID"`
	Name            string            `gorm:"size:100" json:"name"`
	Namespace       string            `gorm:"size:50" json:"namespace"`
	Type            string            `gorm:"size:50" json:"type"` // ClusterIP, NodePort, LoadBalancer, ExternalName
	ClusterIP       string            `gorm:"size:50" json:"clusterIP"`
	ExternalIP      string            `gorm:"size:50" json:"externalIP"`
	ExternalName    string            `gorm:"size:100" json:"externalName"`
	Ports           []ServicePort     `gorm:"-" json:"ports"`    // 不存储到数据库，只用于API响应
	Selector        map[string]string `gorm:"-" json:"selector"` // 不存储到数据库，只用于API响应
	SessionAffinity string            `gorm:"size:50" json:"sessionAffinity"`
	Annotations     map[string]string `gorm:"-" json:"annotations"` // 不存储到数据库，只用于API响应
}

// ConfigMap K8s ConfigMap模型
type ConfigMap struct {
	gorm.Model
	ClusterID   uint              `gorm:"index" json:"clusterID"`
	Name        string            `gorm:"size:100" json:"name"`
	Namespace   string            `gorm:"size:50" json:"namespace"`
	Data        map[string]string `gorm:"-" json:"data"`        // 不存储到数据库，只用于API响应
	BinaryData  map[string]string `gorm:"-" json:"binaryData"`  // 不存储到数据库，只用于API响应
	Annotations map[string]string `gorm:"-" json:"annotations"` // 不存储到数据库，只用于API响应
	CreatedAt   time.Time         `json:"createdAt"`
}

// Secret K8s Secret模型
type Secret struct {
	gorm.Model
	ClusterID   uint              `gorm:"index" json:"clusterID"`
	Name        string            `gorm:"size:100" json:"name"`
	Namespace   string            `gorm:"size:50" json:"namespace"`
	Type        string            `gorm:"size:100" json:"type"` // Opaque, kubernetes.io/service-account-token, kubernetes.io/tls, etc.
	Data        map[string]string `gorm:"-" json:"data"`        // 不存储到数据库，只用于API响应
	BinaryData  map[string]string `gorm:"-" json:"binaryData"`  // 不存储到数据库，只用于API响应
	Annotations map[string]string `gorm:"-" json:"annotations"` // 不存储到数据库，只用于API响应
	CreatedAt   time.Time         `json:"createdAt"`
}

// SLSConfig SLS配置模型
type SLSConfig struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	AccessKeyID     string         `gorm:"size:100" json:"accessKeyId"`
	AccessKeySecret string         `gorm:"size:100" json:"accessKeySecret"`
	RegionID        string         `gorm:"size:50" json:"regionId"`
	DefaultProject  string         `gorm:"size:100" json:"defaultProject"`
	DefaultLogstore string         `gorm:"size:100" json:"defaultLogstore"`
	Timeout         int            `gorm:"default:30" json:"timeout"` // 单位: 秒
	IsActive        bool           `gorm:"default:true" json:"isActive"`
}

// SystemSetting 系统设置模型
type SystemSetting struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"uniqueIndex;size:100" json:"key"` // 键名，如：system.name, login.logo
	Value       string    `gorm:"type:text" json:"value"`          // 键值
	Description string    `gorm:"size:255" json:"description"`     // 描述
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// EmailSetting 邮箱配置模型
type EmailSetting struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SMTPServer     string    `gorm:"size:100" json:"smtpServer"`    // SMTP服务器地址
	SMTPPort       int       `json:"smtpPort"`                      // SMTP服务器端口
	FromEmail      string    `gorm:"size:100" json:"fromEmail"`     // 发件人邮箱地址
	SMTPUsername   string    `gorm:"size:100" json:"smtpUsername"`  // SMTP用户名
	SMTPPassword   string    `gorm:"size:100" json:"smtpPassword"`  // SMTP密码或授权码
	SMTPEncryption string    `gorm:"size:20" json:"smtpEncryption"` // 加密方式：ssl, tls, 空表示不加密
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// TableName 设置EmailSetting模型的表名
func (EmailSetting) TableName() string {
	return "email_settings"
}

// OperationLog 操作日志模型
type OperationLog struct {
	gorm.Model
	UserID    uint      `gorm:"index" json:"userId"`
	Username  string    `gorm:"size:50" json:"username"`
	Menu      string    `gorm:"size:100" json:"menu"`
	Operation string    `gorm:"size:100" json:"operation"`
	API       string    `gorm:"size:255" json:"api"`
	Method    string    `gorm:"size:10" json:"method"` // GET, POST, PUT, DELETE等
	Params    string    `gorm:"type:text" json:"params"`
	Result    string    `gorm:"type:text" json:"result"`
	Status    string    `gorm:"size:20" json:"status"` // success, failed
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"type:text" json:"userAgent"`
	TimeCost  int64     `json:"timeCost"` // 单位：毫秒
	CreatedAt time.Time `json:"createdAt"`
}

// LoginLog 登录日志模型
type LoginLog struct {
	gorm.Model
	UserID    uint      `gorm:"index" json:"userId"`
	Username  string    `gorm:"size:50" json:"username"`
	IP        string    `gorm:"size:50" json:"ip"`
	Location  string    `gorm:"size:100" json:"location"`
	UserAgent string    `gorm:"type:text" json:"userAgent"`
	Status    string    `gorm:"size:20" json:"status"` // success, failed
	Message   string    `gorm:"size:255" json:"message"`
	LoginTime time.Time `json:"loginTime"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName 设置LoginLog模型的表名
func (LoginLog) TableName() string {
	return "login_logs"
}

// TableName 设置OperationLog模型的表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// MD5History MD5加密历史模型
type MD5History struct {
	gorm.Model
	UserID    uint      `gorm:"index" json:"userId"`     // 用户ID
	Username  string    `gorm:"size:50" json:"username"` // 用户名
	Input     string    `gorm:"type:text" json:"input"`  // 原始输入文本
	Output    string    `gorm:"size:32" json:"output"`   // 生成的MD5哈希值
	Format    string    `gorm:"size:20" json:"format"`   // 输出格式：lowercase, uppercase
	Charset   string    `gorm:"size:20" json:"charset"`  // 字符集：utf8, gbk
	CreatedAt time.Time `json:"createdAt"`
}

// TableName 设置MD5History模型的表名
func (MD5History) TableName() string {
	return "md5_histories"
}

// ImageCompressHistory 图片压缩历史模型
type ImageCompressHistory struct {
	gorm.Model
	UserID          uint      `gorm:"index" json:"userId"`      // 用户ID
	Username        string    `gorm:"size:50" json:"username"`  // 用户名
	Filename        string    `gorm:"size:255" json:"filename"` // 原始文件名
	OriginalSize    int64     `json:"originalSize"`             // 原始文件大小（字节）
	CompressedSize  int64     `json:"compressedSize"`           // 压缩后文件大小（字节）
	CompressionRate float64   `json:"compressionRate"`          // 压缩率（百分比）
	Width           int       `json:"width"`                    // 压缩后宽度
	Height          int       `json:"height"`                   // 压缩后高度
	Quality         int       `json:"quality"`                  // 压缩质量
	Format          string    `gorm:"size:20" json:"format"`    // 输出格式
	CreatedAt       time.Time `json:"createdAt"`                // 创建时间
}

// TableName 设置ImageCompressHistory模型的表名
func (ImageCompressHistory) TableName() string {
	return "image_compress_histories"
}
