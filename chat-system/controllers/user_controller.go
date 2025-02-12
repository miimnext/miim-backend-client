package controllers

import (
	"chat-system/config"
	"chat-system/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 定义 JWT 密钥，用于签名生成的 token
var jwtSecret = []byte("your_secret_key")

// UserController 用户控制器
// 处理用户注册和登录功能
type UserController struct{}

// Register 用户注册处理
func (uc *UserController) Register(c *gin.Context) {
	// 定义输入数据结构
	var userInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&userInput); err != nil {
		// 错误处理：请求参数无效
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", userInput.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username already exists",
		})
		return
	}

	// 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// 创建新用户
	newUser := models.User{
		UserID:   fmt.Sprintf("%d", 10000), // 从10000开始递增
		Username: userInput.Username,
		Password: string(hashedPassword), // 保存加密后的密码
	}

	// 将用户信息插入数据库
	if err := config.DB.Create(&newUser).Error; err != nil {
		// 错误处理：数据库插入失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// 成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user_id": newUser.UserID,
	})
}

// Login 用户登录处理
func (uc *UserController) Login(c *gin.Context) {
	// 定义输入数据结构
	var loginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		// 错误处理：请求参数无效
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("username = ?", loginInput.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// 生成 JWT Token
	token, err := generateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// 成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// 生成 JWT Token
func generateJWT(userID string) (string, error) {
	// 定义过期时间
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时过期

	// 创建 JWT 声明
	claims := &jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "chat-system",
	}

	// 创建 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成 token
	return token.SignedString(jwtSecret)
}
