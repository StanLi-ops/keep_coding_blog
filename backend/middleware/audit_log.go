package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"keep_learning_blog/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLog 审计日志中间件
func AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 获取请求信息
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 获取请求体
		var requestBody interface{}
		if c.Request.Body != nil {
			if err := c.ShouldBindJSON(&requestBody); err == nil {
				// 重置请求体，以便后续中间件和处理函数可以读取
				bodyBytes, _ := json.Marshal(requestBody)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 处理请求
		c.Next()

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录审计日志
		logger.AuditLog.WithFields(logger.Fields(map[string]interface{}{
			"user_id":     userID,
			"username":    username,
			"path":        path,
			"method":      method,
			"status_code": statusCode,
			"client_ip":   clientIP,
			"duration":    duration.String(),
			"request":     requestBody,
		})).Info("Audit log")
	}
}
