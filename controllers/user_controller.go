package controllers

import (
	"fmt"
	"go_core/models"
	"go_core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var request struct {
		NickName string `json:"nickname" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// Bind JSON input to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Create a new user object
	user := models.User{
		NickName: request.NickName,
		Username: request.Username,
		Password: request.Password,
	}

	// Call the service layer to create the user
	createdUser, err := services.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    createdUser,
	})
}

// Login 用户登录接口，生成 JWT Token
func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	// 查找用户
	dbUser, err := services.GetUserByUsername(user.Username)
	fmt.Println(dbUser, err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid email or password", "code": 404})
		return
	}

	// 验证密码
	if !services.CheckPassword(dbUser.Password, user.Password) {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid Password or password", "code": 404})
		return
	}

	// 生成 Token
	token, err := services.GenerateToken(*dbUser)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Could not generate token", "code": 404})
		return
	}
	type UserInfo struct {
		ID       uint
		NickName string
		Username string
		Balance  uint
	}
	userInfo := UserInfo{
		ID:       dbUser.ID,
		NickName: dbUser.NickName,
		Username: dbUser.Username,
		Balance:  dbUser.Balance,
	}
	// 返回 token
	c.JSON(http.StatusOK, gin.H{"token": token, "code": 200, "data": userInfo})
}

func GetAllUser(c *gin.Context) {
	// Call the service to get the game
	users, err := services.GetAllUser()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, users) // Return the game as JSON
}
