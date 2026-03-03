package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 系统配置结构体
type Config struct {
	System     SystemConfig     `yaml:"system"`
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	Logger     LoggerConfig     `yaml:"logger"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Email      EmailConfig      `yaml:"email"`
	SLS        SLSConfig        `yaml:"sls"`
	API        APIConfig        `yaml:"api"`
	JWT        JWTConfig        `yaml:"jwt"`
}

// SystemConfig 系统基础配置
type SystemConfig struct {
	Name        string          `yaml:"name"`
	Version     string          `yaml:"version"`
	Environment string          `yaml:"environment"`
	Debug       bool            `yaml:"debug"`
	Title       string          `yaml:"title"`
	Dashboard   DashboardConfig `yaml:"dashboard"`
}

// DashboardConfig 仪表盘配置
type DashboardConfig struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN                string        `yaml:"dsn"`
	MaxOpenConns       int           `yaml:"max_open_conns"`
	MaxIdleConns       int           `yaml:"max_idle_conns"`
	ConnMaxLifetime    time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `yaml:"conn_max_idle_time"`
	SlowQueryThreshold time.Duration `yaml:"slow_query_threshold"`
	LogLevel           string        `yaml:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// PrometheusConfig Prometheus配置
type PrometheusConfig struct {
	Address         string        `yaml:"addr"`
	Timeout         time.Duration `yaml:"timeout"`
	AlertManagerURL string        `yaml:"alert_manager_url"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPServer  string `yaml:"smtp_server"`
	SMTPPort    int    `yaml:"smtp_port"`
	SenderEmail string `yaml:"sender_email"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Encryption  string `yaml:"encryption"`
	Enabled     bool   `yaml:"enabled"`
}

// SLSConfig SLS日志服务配置
type SLSConfig struct {
	AccessKeyID     string        `yaml:"access_key_id"`
	AccessKeySecret string        `yaml:"access_key_secret"`
	RegionID        string        `yaml:"region_id"`
	Project         string        `yaml:"project"`
	Logstore        string        `yaml:"logstore"`
	Timeout         time.Duration `yaml:"timeout"`
	Enabled         bool          `yaml:"enabled"`
}

// APIConfig API配置
type APIConfig struct {
	Prefix         string   `yaml:"prefix"`
	RateLimit      int      `yaml:"rate_limit"`
	CORS           bool     `yaml:"cors"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string        `yaml:"secret"`
	ExpiresIn        time.Duration `yaml:"expires_in"`
	RefreshExpiresIn time.Duration `yaml:"refresh_expires_in"`
}

// LoadConfig 从yaml文件和环境变量加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 1. 创建默认配置
	defaultConfig := &Config{
		System: SystemConfig{
			Name:        "DeanTech",
			Version:     "1.0.0",
			Environment: "production",
			Debug:       false,
			Title:       "DeanTech 集群管理平台",
			Dashboard: DashboardConfig{
				Title:       "DeanTech 集群管理平台",
				Description: "Kubernetes 集群管理和监控平台",
			},
		},
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			DSN:                "deantech:jkFmTJhsTBKNZsXz@tcp(43.143.194.230:3306)/deantech?charset=utf8mb4&parseTime=True&loc=Local",
			MaxOpenConns:       100,
			MaxIdleConns:       20,
			ConnMaxLifetime:    1 * time.Hour,
			ConnMaxIdleTime:    30 * time.Minute,
			SlowQueryThreshold: 200 * time.Millisecond,
			LogLevel:           "info",
		},
		Redis: RedisConfig{
			Addr:         "127.0.0.1:6379",
			Password:     "",
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 5,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		Logger: LoggerConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			FilePath:   "./logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     7,
			Compress:   true,
		},
		Prometheus: PrometheusConfig{
			Address:         "http://localhost:9090",
			Timeout:         5 * time.Second,
			AlertManagerURL: "http://localhost:9093",
		},
		Email: EmailConfig{
			SMTPServer:  "smtp.example.com",
			SMTPPort:    587,
			SenderEmail: "noreply@example.com",
			Username:    "noreply@example.com",
			Password:    "password",
			Encryption:  "tls",
			Enabled:     false,
		},
		SLS: SLSConfig{
			AccessKeyID:     "",
			AccessKeySecret: "",
			RegionID:        "cn-hangzhou",
			Project:         "deantech-log",
			Logstore:        "cluster-logs",
			Timeout:         10 * time.Second,
			Enabled:         false,
		},
		API: APIConfig{
			Prefix:         "/api",
			RateLimit:      100,
			CORS:           true,
			AllowedOrigins: []string{"*"},
		},
		JWT: JWTConfig{
			Secret:           "your-secret-key",
			ExpiresIn:        24 * time.Hour,
			RefreshExpiresIn: 72 * time.Hour,
		},
	}

	// 2. 从yaml文件加载配置
	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		yamlDecoder := yaml.NewDecoder(file)
		if err := yamlDecoder.Decode(defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	// 3. 从环境变量加载配置，覆盖yaml文件中的配置
	// 服务器配置
	if port := getEnv("SERVER_PORT", ""); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			defaultConfig.Server.Port = p
		}
	}
	if host := getEnv("SERVER_HOST", ""); host != "" {
		defaultConfig.Server.Host = host
	}

	// 数据库配置
	if dsn := getEnv("DATABASE_DSN", ""); dsn != "" {
		defaultConfig.Database.DSN = dsn
	}

	// Redis配置
	if addr := getEnv("REDIS_ADDRESS", ""); addr != "" {
		defaultConfig.Redis.Addr = addr
	}
	if password := getEnv("REDIS_PASSWORD", ""); password != "" {
		defaultConfig.Redis.Password = password
	}
	if db := getEnv("REDIS_DB", ""); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			defaultConfig.Redis.DB = d
		}
	}

	// Prometheus配置
	if addr := getEnv("PROMETHEUS_ADDRESS", ""); addr != "" {
		defaultConfig.Prometheus.Address = addr
	}
	if timeout := getEnv("PROMETHEUS_TIMEOUT", ""); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			defaultConfig.Prometheus.Timeout = time.Duration(t) * time.Second
		}
	}
	if alertURL := getEnv("ALERTMANAGER_URL", ""); alertURL != "" {
		defaultConfig.Prometheus.AlertManagerURL = alertURL
	}

	// 日志配置
	if logLevel := getEnv("LOG_LEVEL", ""); logLevel != "" {
		defaultConfig.Logger.Level = logLevel
	}

	return defaultConfig, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
