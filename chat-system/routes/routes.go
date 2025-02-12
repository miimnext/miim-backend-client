package routes

import (
	"chat-system/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	// 创建 UserController 实例
	userController := &controllers.UserController{}

	// 注册路由
	router.POST("/register", userController.Register) // 绑定注册接口
	router.POST("/login", userController.Login)       // 绑定登录接口

	router.GET("/ws/:user_id", controllers.WebSocketHandler)
}
