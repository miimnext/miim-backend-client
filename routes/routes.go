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

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		controllers.HandleWebSocket(c.Writer, c.Request)
	})

	// Protected routes
	protected := r.Group("/api")
	protected.POST("/register", controllers.RegisterUser)
	protected.POST("/login", controllers.LoginUser)
	protected.GET("/posts", controllers.GetAllPosts)
	protected.GET("/postsByUser/:id", controllers.GetPostsByUser)

	protected.GET("/tags", controllers.GetTags)
	protected.GET("/categorys", controllers.GetCategorys)
	protected.GET("/posts/:id", controllers.GetPostByID)
	protected.DELETE("/posts/:id", controllers.DeletePost)
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.Use(middlewares.TokenAuthMiddleware())
		protected.POST("/posts", controllers.CreatePost)
		protected.GET("/userinfo", controllers.GetUserInfo)
		protected.GET("/users", controllers.GetAllUser)
		protected.POST("/upload", controllers.UploadFile)
	}
	return r
}
