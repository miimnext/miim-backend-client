package controllers

import (
	"go_core/config"
	"go_core/models"
	"go_core/services"
	"go_core/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	ID        uint      `json:"id"`
	NickName  string    `json:"nickName"`
	Username  string    `json:"userName"`
	Avatar    string    `json:"avatar"`
	Balance   uint      `json:"balance"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func RegisterUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// Bind JSON input to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondFailed(c, "Invalid input")
		return
	}

	// Create a new user object
	user := models.User{
		Username: request.Username,
		Password: request.Password,
	}

	// Call the service layer to create the user
	createdUser, err := services.CreateUser(user)
	if err != nil {
		utils.RespondFailed(c, err.Error())
		return
	}

	// 生成 Token
	token, err := services.GenerateToken(createdUser)
	if err != nil {
		utils.RespondFailed(c, "Could not generate token")
		return
	}
	// 返回 token
	utils.RespondSuccess(c, gin.H{"token": token}, nil)
}

// Login 用户登录接口，生成 JWT Token
func LoginUser(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		utils.RespondFailed(c, "Invalid input")
		return
	}
	var user models.User
	user.Username = loginReq.Username
	user.Password = loginReq.Password

	// 查找用户
	dbUser, err := services.GetUserByUsername(user.Username)
	if err != nil {
		utils.RespondFailed(c, "用户名密码错误")
		return
	}
	// 验证密码
	if !services.CheckPassword(dbUser.Password, user.Password) {
		utils.RespondFailed(c, "用户名密码错误")
		return
	}

	// 生成 Token
	token, err := services.GenerateToken(*dbUser)
	if err != nil {
		utils.RespondFailed(c, "Could not generate token")
		return
	}

	// 返回 token
	utils.RespondSuccess(c, gin.H{"token": token}, nil)
}

func GetAllUser(c *gin.Context) {
	users, err := services.GetAllUser()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	c.JSON(http.StatusOK, users) // Return the game as JSON
}

// 根据token获取自己信息
func GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户信息
	user, exists := c.Get("user")
	if !exists {
		// 如果用户信息不存在，返回 404 错误
		utils.RespondFailed(c, "User not found")
		return
	}

	// 假设 user 是 *models.User 类型，进行类型断言
	userInfo, ok := user.(*models.User)
	if !ok {
		// 如果类型断言失败，返回 400 错误
		utils.RespondFailed(c, "Invalid user data")
		return
	}
	data := UserInfo{
		ID:        userInfo.ID,
		Username:  userInfo.Username,
		NickName:  userInfo.NickName,
		Balance:   userInfo.Balance,
		Avatar:    userInfo.Avatar,
		CreatedAt: userInfo.CreatedAt,
	}
	utils.RespondSuccess(c, data, nil)
}

func GetAuthorByID(c *gin.Context) {
	id := c.Param("id") // 从路径参数中获取 id
	user, err := services.GetUserByID(id)
	if err != nil {
		utils.RespondFailed(c, "Invalid id ")
		return
	}
	utils.RespondSuccess(c, user, nil)

}

func UpdateUserInfo(c *gin.Context) {
	// 从上下文中获取用户信息
	user, exists := c.Get("user")
	if !exists {
		// 如果用户信息不存在，返回 404 错误
		utils.RespondFailed(c, "User not found")
		return
	}

	// 假设 user 是 *models.User 类型，进行类型断言
	userInfo, ok := user.(*models.User)
	if !ok {
		// 如果类型断言失败，返回 400 错误
		utils.RespondFailed(c, "Invalid user data")
		return
	}

	// 获取用户传来的更新数据（只包含更新的字段）
	var updateData map[string]interface{} // 使用 map 来接收传入的字段
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.RespondFailed(c, "Invalid input data")
		return
	}

	// 使用更新的数据直接更新数据库
	if err := config.DB.Model(&models.User{}).Where("id = ?", userInfo.ID).Updates(updateData).Error; err != nil {
		// 如果数据库更新失败，返回 500 错误
		utils.RespondFailed(c, "Failed to update user information")
		return
	}

	// 更新成功，返回成功响应
	utils.RespondSuccess(c, gin.H{
		"message": "User information updated successfully",
	}, nil)
}
