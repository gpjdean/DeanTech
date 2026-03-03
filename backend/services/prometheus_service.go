package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"deantech/config"
	"time"

	"github.com/gin-gonic/gin"
)

type PrometheusService struct {
	config *config.PrometheusConfig
	client *http.Client
}

func NewPrometheusService(config *config.PrometheusConfig) *PrometheusService {
	return &PrometheusService{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// PrometheusAlert 表示Prometheus告警结构
type PrometheusAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// PrometheusAlertResponse 表示Prometheus告警API响应
type PrometheusAlertResponse struct {
	Status string            `json:"status"`
	Data   PrometheusAlertData `json:"data"`
}

type PrometheusAlertData struct {
	Alerts []PrometheusAlert `json:"alerts"`
}

// GetPrometheusAlerts 从Prometheus获取告警列表
func (s *PrometheusService) GetPrometheusAlerts(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/alerts", s.config.Address)
	resp, err := s.client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var alertResp PrometheusAlertResponse
	if err := json.Unmarshal(body, &alertResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alertResp.Data.Alerts)
}

// GetPrometheusRules 从Prometheus获取规则列表
func (s *PrometheusService) GetPrometheusRules(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/rules", s.config.Address)
	resp, err := s.client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

// GetPrometheusTargets 从Prometheus获取目标列表
func (s *PrometheusService) GetPrometheusTargets(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/targets", s.config.Address)
	resp, err := s.client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}