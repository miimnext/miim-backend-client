package main

import (
	"chat-system/config"
	"chat-system/models"
	"chat-system/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()
	// 自动迁移
	models.Migrate()

	// 初始化 Gin 引擎
	r := gin.Default()

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
