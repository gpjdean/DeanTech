package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"deantech/models"
)

// EmailService 邮件服务结构体
type EmailService struct {
	db *gorm.DB
}

// NewEmailService 创建邮件服务实例
func NewEmailService(db *gorm.DB) *EmailService {
	return &EmailService{db: db}
}

// RegisterRoutes 注册邮件相关路由
func (s *EmailService) RegisterRoutes(router *gin.RouterGroup) {
	emailGroup := router.Group("/email")
	{
		emailGroup.POST("/test", s.TestEmailSendHandler) // 测试发送邮件
	}
}

// TestEmailSendHandler 测试发送邮件
func (s *EmailService) TestEmailSendHandler(c *gin.Context) {
	var requestBody struct {
		ToEmail string `json:"toEmail" binding:"required,email"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 获取邮件配置
	emailConfig, err := s.getEmailConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get email config: " + err.Error(),
		})
		return
	}

	// 发送测试邮件
	if err := s.sendTestEmail(requestBody.ToEmail, emailConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send test email: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test email sent successfully",
	})
}

// EmailConfig 邮件配置结构体
type EmailConfig struct {
	SMTPServer   string
	SMTPPort     int
	FromEmail    string
	SMTPUsername string
	SMTPPassword string
	SMTPEncryption string
}

// getEmailConfig 获取邮件配置
func (s *EmailService) getEmailConfig() (*EmailConfig, error) {
	// 查询邮箱配置
	var emailSetting models.EmailSetting
	result := s.db.First(&emailSetting)
	if result.Error != nil {
		return nil, fmt.Errorf("email config not found: %w", result.Error)
	}

	// 构建并返回配置
	return &EmailConfig{
		SMTPServer:     emailSetting.SMTPServer,
		SMTPPort:       emailSetting.SMTPPort,
		FromEmail:      emailSetting.FromEmail,
		SMTPUsername:   emailSetting.SMTPUsername,
		SMTPPassword:   emailSetting.SMTPPassword,
		SMTPEncryption: emailSetting.SMTPEncryption,
	}, nil
}

// sendTestEmail 发送测试邮件
func (s *EmailService) sendTestEmail(toEmail string, config *EmailConfig) error {
	// 邮件内容
	subject := "[测试] 邮件发送功能测试"
	body := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>邮件发送测试</title>
</head>
<body>
    <h2>邮件发送测试成功！</h2>
    <p>这是一封测试邮件，用于验证系统的邮件发送功能。</p>
    <p>如果您收到这封邮件，说明系统的邮件发送功能已经配置正确。</p>
    <p>此邮件由系统自动发送，请勿回复。</p>
</body>
</html>
`

	// 构建邮件
	email := fmt.Sprintf("From: %s\r\n", config.FromEmail) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n\r\n" +
		body

	// 发送邮件
	return s.sendEmail(config, toEmail, email)
}

// sendEmail 发送邮件的核心功能
func (s *EmailService) sendEmail(config *EmailConfig, toEmail string, email string) error {
	// 构建SMTP地址
	addr := fmt.Sprintf("%s:%d", config.SMTPServer, config.SMTPPort)

	// 认证信息
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPServer)

	// 根据加密方式选择发送方式
	switch config.SMTPEncryption {
	case "ssl":
		// SSL加密
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			ServerName: config.SMTPServer,
			InsecureSkipVerify: false,
		})
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, config.SMTPServer)
		if err != nil {
			return err
		}
		defer client.Close()

		// 认证
		if err := client.Auth(auth); err != nil {
			return err
		}

		// 设置发件人和收件人
		if err := client.Mail(config.FromEmail); err != nil {
			return err
		}
		if err := client.Rcpt(toEmail); err != nil {
			return err
		}

		// 获取数据写入器
		w, err := client.Data()
		if err != nil {
			return err
		}
		defer w.Close()

		// 写入邮件内容
		_, err = w.Write([]byte(email))
		if err != nil {
			return err
		}

	case "tls":
		// TLS加密（STARTTLS）
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, config.SMTPServer)
		if err != nil {
			return err
		}
		defer client.Close()

		// 开启TLS
		if err := client.StartTLS(&tls.Config{
			ServerName: config.SMTPServer,
			InsecureSkipVerify: false,
		}); err != nil {
			return err
		}

		// 认证
		if err := client.Auth(auth); err != nil {
			return err
		}

		// 设置发件人和收件人
		if err := client.Mail(config.FromEmail); err != nil {
			return err
		}
		if err := client.Rcpt(toEmail); err != nil {
			return err
		}

		// 获取数据写入器
		w, err := client.Data()
		if err != nil {
			return err
		}
		defer w.Close()

		// 写入邮件内容
		_, err = w.Write([]byte(email))
		if err != nil {
			return err
		}

	default:
		// 无加密
		return smtp.SendMail(addr, auth, config.FromEmail, []string{toEmail}, []byte(email))
	}

	return nil
}
