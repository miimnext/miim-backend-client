package middlewares

import (
	"go_core/services"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 验证 JWT Token 是否有效
// AuthMiddleware 验证 JWT Token 是否有效
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 中提取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Next()
			return
		}

		// JWT 的格式通常是 "Bearer <token>"
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 如果 token 格式不正确，可以记录日志但不拦截请求
			c.Next()
			return
		}
		tokenString = parts[1]

		// 验证 token
		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			// 如果验证失败，也不拦截请求，只是打印错误日志
			// c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Next()
			return
		}

		// 将解析出来的 claims 存储到上下文中，方便后续的处理
		c.Set("user", claims)

		// 继续处理请求
		c.Next()
	}
}
