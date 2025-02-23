package middleware

import (
	"net/http"
	"time"

	"keep_learning_blog/config"
	"keep_learning_blog/db"
	"keep_learning_blog/models"
	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
)

const (
	maxLoginAttempts = 5                // 最大登录失败次数
	lockDuration     = 15 * time.Minute // 锁定时间
)

type LoginLimiter struct {
	config *config.Config
}

func NewLoginLimiter(config *config.Config) *LoginLimiter {
	return &LoginLimiter{
		config: config,
	}
}

// CheckLoginAttempts 检查登录尝试次数中间件
func (l *LoginLimiter) CheckLoginAttempts() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户名或邮箱（从请求体中获取）
		var loginRequest models.LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			logger.Log.WithError(err).Warn("Invalid login request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		// 确定用户标识符（用户名或邮箱）
		identifier := loginRequest.Username
		if identifier == "" {
			identifier = loginRequest.Email
		}

		// 检查是否被锁定
		if db.IsLoginLocked(c.Request.Context(), identifier) {
			remaining := db.GetLoginLockRemainingTime(c.Request.Context(), identifier)

			logger.Log.WithFields(logger.Fields(map[string]interface{}{
				"identifier": identifier,
				"remaining":  remaining,
			})).Warn("Account is temporarily locked")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":             "Account is temporarily locked",
				"remaining_minutes": int(remaining.Minutes()),
			})
			c.Abort()
			return
		}

		// 将用户标识符存储在上下文中，供后续使用
		c.Set("login_identifier", identifier)
		c.Set("login_request", loginRequest)

		c.Next()
	}
}

// RecordLoginAttempt 记录登录尝试结果
func (l *LoginLimiter) RecordLoginAttempt(ctx *gin.Context, success bool, identifier string) error {
	return db.RecordLoginAttempt(ctx.Request.Context(), identifier, success, maxLoginAttempts, lockDuration)
}
