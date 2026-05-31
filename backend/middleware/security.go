package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SecurityMiddlewares 获取所有安全中间件
func SecurityMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CORSMiddleware(),
		XSSProtectionMiddleware(),
	}
}

// CORSMiddleware CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * 60 * 60, // 12小时
	})
}

// XSSProtectionMiddleware XSS 防护中间件
func XSSProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// XSS 防护
		c.Header("X-XSS-Protection", "1; mode=block")
		// 防止 MIME 类型混淆
		c.Header("X-Content-Type-Options", "nosniff")
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		// 严格传输安全（生产环境使用）
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}