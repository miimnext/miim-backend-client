package routes

import (
	"go_core/controllers"
	"go_core/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// 初始化 Gin 引擎
	r := gin.Default()

	// 配置跨域中间件（CORS）
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许所有来源（可以修改为前端地址如 "http://localhost:3000"）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		AllowCredentials: true,                                                // 是否允许发送 cookies
	}
	// 静态文件目录
	r.Static("/static", "./static")
	// 使用 CORS 中间件配置
	r.Use(cors.New(corsConfig))

	// 使用自定义的日志中间件
	r.Use(middlewares.Logger())

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		controllers.HandleWebSocket(c.Writer, c.Request) // 处理 WebSocket 连接
	})

	// 不需要身份验证的路由（公共路由）
	// 用户注册和登录
	protected := r.Group("/api") // 创建一个路由组，路径前缀为 /api
	protected.Use(middlewares.AuthMiddleware())
	protected.POST("/register", controllers.RegisterUser)   // 用户注册
	protected.POST("/login", controllers.LoginUser)         // 用户登录
	protected.GET("/author/:id", controllers.GetAuthorByID) // 根据 ID 获取作者信息

	// 文章相关路由
	protected.GET("/posts", controllers.GetAllPosts)              // 获取所有文章
	protected.GET("/postsByUser/:id", controllers.GetPostsByUser) // 根据用户 ID 获取文章
	protected.GET("/posts/:id", controllers.GetPostByID)          // 根据文章 ID 获取单篇文章
	protected.DELETE("/posts/:id", controllers.DeletePost)        // 根据文章 ID 删除文章

	// 获取标签和分类
	protected.GET("/tags", controllers.GetTags)           // 获取所有标签
	protected.GET("/categorys", controllers.GetCategorys) // 获取所有分类
	{

		// 需要 token 验证的路由
		protected.Use(middlewares.TokenAuthMiddleware())            // 验证用户的 token
		protected.POST("/posts", controllers.CreatePost)            // 创建文章
		protected.POST("/posts/reaction", controllers.PostReaction) // 点赞/取消点赞
		protected.GET("/userinfo", controllers.GetUserInfo)         // 获取当前用户信息
		protected.GET("/users", controllers.GetAllUser)             // 获取所有用户
		protected.POST("/upload", controllers.UploadFile)           // 上传文件
	}

	return r // 返回配置好的路由引擎
}
