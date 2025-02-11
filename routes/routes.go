package routes

import (
	"go_core/controllers"
	"go_core/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 配置跨域中间件
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许的域名，可以是前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		AllowCredentials: true,                                                // 是否允许发送 cookies
	}

	// 使用 CORS 中间件
	r.Use(cors.New(corsConfig))
	r.Use(middlewares.Logger())
	// Public routes

	// Protected routes
	protected := r.Group("/api")
	protected.POST("/register", controllers.RegisterUser)
	protected.POST("/login", controllers.LoginUser)

	// 文章相关的路由
	protected.POST("/posts", controllers.CreatePost)     // 创建文章
	protected.GET("/posts", controllers.GetAllPosts)     // 获取所有文章
	protected.GET("/posts/:id", controllers.GetPostByID) // 获取单篇文章
	protected.DELETE("/posts/:id", controllers.DeletePost)
	protected.Use(middlewares.AuthMiddleware())

	{
		protected.Use(middlewares.TokenAuthMiddleware())
		protected.GET("/users", controllers.GetAllUser)
		protected.POST("/upload", controllers.UploadFile)
	}
	return r
}
