package middleware

import (
	"keep_learning_blog/config"
	"net/http"
	"strconv"
	"strings"

	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
)

// CORS 创建 CORS 中间件
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 记录 CORS 请求
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"origin": origin,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})).Debug("Processing CORS request")

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

		// 记录预检请求
		if c.Request.Method == http.MethodOptions {
			logger.Log.WithFields(logger.Fields(map[string]interface{}{
				"origin": origin,
				"path":   c.Request.URL.Path,
			})).Debug("Handling CORS preflight request")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
