package services

import (
	"bytes"
	"deantech/models"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Register 用户注册
func (s *UserService) Register(username, password, email, name, phone, department, position string) (*models.User, error) {
	// 检查用户名是否已存在（包括软删除记录）
	var existingUser models.User
	if err := s.db.Unscoped().Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 检查邮箱是否已存在（包括软删除记录）
	if err := s.db.Unscoped().Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username:   username,
		Password:   string(hashedPassword),
		Email:      email,
		Name:       name,
		Phone:      phone,
		Department: department,
		Position:   position,
		Role:       "user",     // 默认角色为普通用户
		Status:     "inactive", // 新注册用户默认禁用
	}

	// 调试日志：打印创建的用户
	log.Printf("Creating user: %+v", user)

	// 创建用户
	if err := s.db.Create(user).Error; err != nil {
		// 处理唯一约束错误
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				return nil, errors.New("用户名已存在")
			} else if strings.Contains(err.Error(), "email") {
				return nil, errors.New("邮箱已存在")
			}
		}
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*models.User, error) {
	// 查找用户
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, password, email, name, phone, department, position string) (*models.User, error) {
	// 检查用户名是否已存在（包括软删除记录）
	var count int64
	if err := s.db.Unscoped().Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return nil, err
	} else if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在（包括软删除记录）
	if err := s.db.Unscoped().Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return nil, err
	} else if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username:   username,
		Password:   string(hashedPassword),
		Email:      email,
		Name:       name,
		Phone:      phone,
		Department: department,
		Position:   position,
		Role:       "user",     // 默认角色为普通用户
		Status:     "inactive", // 新注册用户默认禁用
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

// GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile 更新当前用户信息
func (s *UserService) UpdateProfile(userID uint, email, name, phone, department, position, remark, avatar string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 更新用户信息
	user.Email = email
	user.Name = name
	user.Phone = phone
	user.Department = department
	user.Position = position
	user.Remark = remark
	// 只有当avatar不为空时才更新头像
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser 删除用户（硬删除，直接从数据库中删除记录）
func (s *UserService) DeleteUser(userID uint) error {
	return s.db.Unscoped().Delete(&models.User{}, userID).Error
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 设置新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.db.Save(&user).Error
}

// GenerateRandomPassword 生成6位随机密码
func (s *UserService) GenerateRandomPassword() string {
	const charset = "0123456789"
	const length = 6
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

// VerifyEmail 验证邮箱是否存在
func (s *UserService) VerifyEmail(email string) error {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("邮箱不存在")
		}
		return err
	}
	return nil
}

// VerifyUsername 验证用户名是否存在
func (s *UserService) VerifyUsername(username string) error {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户名不存在")
		}
		return err
	}
	return nil
}

// ResetPasswordByEmailAndUsername 根据邮箱和用户名重置密码
func (s *UserService) ResetPasswordByEmailAndUsername(email, username string) (string, error) {
	// 首先验证用户名是否存在
	if err := s.VerifyUsername(username); err != nil {
		return "", err
	}

	// 然后验证邮箱和用户名是否匹配
	var user models.User
	if err := s.db.Where("email = ? AND username = ?", email, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("邮箱与用户名不匹配")
		}
		return "", err
	}

	// 生成6位随机密码
	newPassword := s.GenerateRandomPassword()

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return "", err
	}

	return newPassword, nil
}

// SendResetPasswordEmail 发送重置密码邮件
func (s *UserService) SendResetPasswordEmail(toEmail, username, newPassword string) error {
	// 获取邮件配置
	emailService := NewEmailService(s.db)
	emailConfig, err := emailService.getEmailConfig()
	if err != nil {
		return fmt.Errorf("获取邮件配置失败: %w", err)
	}

	// 邮件内容
	subject := "密码重置通知"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>密码重置通知</title>
</head>
<body>
    <h2>密码重置成功！</h2>
    <p>尊敬的 %s 用户：</p>
    <p>您的密码已成功重置，新密码为：<strong>%s</strong></p>
    <p>请使用新密码登录系统，建议登录后立即修改密码。</p>
    <p>此邮件由系统自动发送，请勿回复。</p>
</body>
</html>
`, username, newPassword)

	// 构建邮件
	email := fmt.Sprintf("From: %s\r\n", emailConfig.FromEmail) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n\r\n" +
		body

	// 发送邮件
	return emailService.sendEmail(emailConfig, toEmail, email)
}

// API处理方法

// RegisterHandler 处理用户注册请求
func (s *UserService) RegisterHandler(c *gin.Context) {
	var req struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required,min=6"`
		Name       string `json:"name" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
		Position   string `json:"position"`
	}

	// 调试日志：打印请求体
	body, _ := io.ReadAll(c.Request.Body)
	log.Printf("Register request body: %s", string(body))
	// 重置请求体，以便ShouldBindJSON可以重新读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调试日志：打印请求参数
	log.Printf("Register request: Username=%s, Password=%s, Email=%s, Name=%s, Phone=%s, Department=%s, Position=%s",
		req.Username, req.Password, req.Email, req.Name, req.Phone, req.Department, req.Position)

	// 检查用户名是否已存在（包括软删除记录）
	var existingUser models.User
	if err := s.db.Unscoped().Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "注册失败，请稍后重试"})
		return
	}

	// 检查邮箱是否已存在（包括软删除记录）
	if err := s.db.Unscoped().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "注册失败，请稍后重试"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "注册失败，请稍后重试"})
		return
	}

	// 创建用户
	user := &models.User{
		Username:   req.Username,
		Password:   string(hashedPassword),
		Email:      req.Email,
		Name:       req.Name,
		Phone:      req.Phone,
		Department: req.Department,
		Position:   req.Position,
		Role:       "user",     // 默认角色为普通用户
		Status:     "inactive", // 新注册用户默认禁用
	}

	// 调试日志：打印创建的用户
	log.Printf("Creating user: %+v", user)

	// 直接创建用户
	if err := s.db.Create(user).Error; err != nil {
		// 处理唯一约束错误
		errorMsg := "注册失败，请稍后重试"
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "idx_users_username") {
				errorMsg = "用户名已存在"
			} else if strings.Contains(err.Error(), "idx_users_email") {
				errorMsg = "邮箱已存在"
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"department": user.Department,
		"position":   user.Position,
		"role":       user.Role,
	})
}

// LoginHandler 处理用户登录请求
func (s *UserService) LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录失败的登录日志
		loginLog := models.LoginLog{
			Username:  req.Username,
			IP:        c.ClientIP(),
			Location:  "", // 后续可以添加IP地址解析
			UserAgent: c.Request.UserAgent(),
			Status:    "failed",
			Message:   "请求参数错误: " + err.Error(),
			LoginTime: time.Now(),
		}
		s.db.Create(&loginLog)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.Login(req.Username, req.Password)
	if err != nil {
		// 记录失败的登录日志
		loginLog := models.LoginLog{
			Username:  req.Username,
			IP:        c.ClientIP(),
			Location:  "", // 后续可以添加IP地址解析
			UserAgent: c.Request.UserAgent(),
			Status:    "failed",
			Message:   err.Error(),
			LoginTime: time.Now(),
		}
		s.db.Create(&loginLog)

		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 记录成功的登录日志
	loginLog := models.LoginLog{
		UserID:    user.ID,
		Username:  user.Username,
		IP:        c.ClientIP(),
		Location:  "", // 后续可以添加IP地址解析
		UserAgent: c.Request.UserAgent(),
		Status:    "success",
		Message:   "登录成功",
		LoginTime: time.Now(),
	}
	s.db.Create(&loginLog)

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// GetUserHandler 处理获取用户信息请求
func (s *UserService) GetUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := s.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// UpdateUserHandler 处理更新用户信息请求
func (s *UserService) UpdateUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		Username   string `json:"username"`
		Password   string `json:"password" binding:"omitempty,min=6"`
		Email      string `json:"email" binding:"omitempty,email"`
		Status     string `json:"status"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
		Position   string `json:"position"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Department != "" {
		user.Department = req.Department
	}
	if req.Position != "" {
		user.Position = req.Position
	}

	// 如果密码不为空，则更新密码
	if req.Password != "" {
		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := s.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"phone":      user.Phone,
		"department": user.Department,
		"position":   user.Position,
		"role":       user.Role,
		"status":     user.Status,
	})
}

// ChangePasswordHandler 处理修改密码请求
func (s *UserService) ChangePasswordHandler(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"user_id" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.ChangePassword(req.UserID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// ListUsersHandler 处理获取用户列表请求
func (s *UserService) ListUsersHandler(c *gin.Context) {
	users, err := s.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateUserHandler 处理创建用户请求
func (s *UserService) CreateUserHandler(c *gin.Context) {
	var req struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required,min=6"`
		Email      string `json:"email" binding:"required,email"`
		Name       string `json:"name" binding:"required"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
		Position   string `json:"position"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.CreateUser(req.Username, req.Password, req.Email, req.Name, req.Phone, req.Department, req.Position)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"phone":      user.Phone,
		"department": user.Department,
		"position":   user.Position,
		"role":       user.Role,
	})
}

// GetCurrentUserHandler 处理获取当前用户信息请求
func (s *UserService) GetCurrentUserHandler(c *gin.Context) {
	// 从请求头中获取用户名（简单的认证方式）
	username := c.GetHeader("X-Username")
	var user models.User

	if username == "" {
		// 如果没有用户名，尝试获取第一个用户（兼容旧逻辑）
		if err := s.db.First(&user).Error; err != nil {
			// 如果没有用户，创建一个默认用户
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			user = models.User{
				Username:   "admin",
				Password:   string(hashedPassword),
				Email:      "admin@example.com",
				Role:       "admin",
				Name:       "管理员",
				Phone:      "13800138000",
				Department: "技术部",
				Position:   "系统管理员",
				Status:     "active", // 设置默认用户为激活状态
			}
			s.db.Create(&user)
		}
	} else {
		// 根据用户名获取用户信息
		if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
	}

	// 确保用户状态为active
	if user.Status != "active" {
		user.Status = "active"
		s.db.Save(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"phone":      user.Phone,
		"department": user.Department,
		"position":   user.Position,
		"remark":     user.Remark,
		"role":       user.Role,
		"avatar":     user.Avatar,
		"status":     user.Status,
	})
}

// UpdateProfileHandler 处理更新当前用户信息请求
func (s *UserService) UpdateProfileHandler(c *gin.Context) {
	var req struct {
		UserID     uint   `json:"user_id" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Name       string `json:"name" binding:"required"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
		Position   string `json:"position"`
		Remark     string `json:"remark"` // 备注
		Avatar     string `json:"avatar"` // 头像base64数据
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// 自定义错误信息
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag") {
			errorMsg = "姓名为必填项"
		} else if strings.Contains(errorMsg, "Key: 'Email' Error:Field validation for 'Email' failed on the 'required' tag") {
			errorMsg = "邮箱为必填项"
		} else if strings.Contains(errorMsg, "Key: 'Email' Error:Field validation for 'Email' failed on the 'email' tag") {
			errorMsg = "请输入有效的邮箱地址"
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	user, err := s.UpdateProfile(req.UserID, req.Email, req.Name, req.Phone, req.Department, req.Position, req.Remark, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"phone":      user.Phone,
		"department": user.Department,
		"position":   user.Position,
		"remark":     user.Remark,
		"role":       user.Role,
		"avatar":     user.Avatar,
	})
}

// DeleteUserHandler 处理删除用户请求
func (s *UserService) DeleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := s.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// VerifyEmailHandler 处理验证邮箱请求
func (s *UserService) VerifyEmailHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邮箱地址"})
		return
	}

	if err := s.VerifyEmail(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邮箱验证通过"})
}

// ResetPasswordHandler 处理重置密码请求
func (s *UserService) ResetPasswordHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 重置密码
	newPassword, err := s.ResetPasswordByEmailAndUsername(req.Email, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 发送重置密码邮件
	if err := s.SendResetPasswordEmail(req.Email, req.Username, newPassword); err != nil {
		log.Printf("发送重置密码邮件失败: %v", err)
		// 密码已经重置成功，邮件发送失败不影响整体流程
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功，新密码已发送到您的邮箱"})
}
