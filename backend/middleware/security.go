package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 添加安全相关的 HTTP 响应头
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 基本安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 只在 HTTPS 环境下启用 HSTS
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// CSP 策略配置
		csp := []string{
			"default-src 'self'",                // 默认只允许同源
			"img-src 'self' data: https:",       // 允许加载 HTTPS 图片和 base64 图片
			"script-src 'self' 'unsafe-inline'", // 允许内联脚本（如果需要的话）
			"style-src 'self' 'unsafe-inline'",  // 允许内联样式
			"font-src 'self' https:",            // 允许加载字体
			"connect-src 'self'",                // API 请求限制
		}
		c.Header("Content-Security-Policy", join(csp))

		c.Next()
	}
}

// join 将 CSP 策略数组连接成字符串
func join(policies []string) string {
	result := ""
	for i, policy := range policies {
		if i > 0 {
			result += "; "
		}
		result += policy
	}
	return result
}
