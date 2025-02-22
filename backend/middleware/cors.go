package middleware

import (
	"keep_coding_blog/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 创建 CORS 中间件
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查请求源是否在允许列表中
		allowOrigin := ""
		for _, o := range cfg.CORS.AllowOrigins {
			if o == origin {
				allowOrigin = o
				break
			}
		}

		// 如果请求源在允许列表中，设置相应的 CORS 头
		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowMethods, ","))
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowHeaders, ","))
			c.Header("Access-Control-Expose-Headers", strings.Join(cfg.CORS.ExposeHeaders, ","))
			c.Header("Access-Control-Max-Age", strconv.Itoa(cfg.CORS.MaxAge))

			if cfg.CORS.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
