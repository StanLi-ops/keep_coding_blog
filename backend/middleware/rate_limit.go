package middleware

import (
	"context"
	"fmt"
	"net/http"

	"keep_coding_blog/config"
	"keep_coding_blog/db"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	config *config.Config
}

func NewRateLimiter(config *config.Config) *RateLimiter {
	return &RateLimiter{
		config: config,
	}
}

// RateLimit 创建限流中间件
func (rl *RateLimiter) RateLimit(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端标识（优先使用用户ID，其次使用IP）
		identifier := getClientIdentifier(c)

		// 构造 Redis key
		key := fmt.Sprintf("%s%s:%s",
			rl.config.Redis.RatePrefix,
			c.FullPath(),
			identifier,
		)

		// 获取当前请求数
		count, err := db.GetRateLimit(context.Background(), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		// 如果是第一次请求，设置初始值和过期时间
		if count == 0 {
			err = db.SetRateLimit(context.Background(), key, rl.config.RateLimit.Duration)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit set failed"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// 检查是否超过限制
		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": rl.config.RateLimit.Duration.Seconds(),
			})
			c.Abort()
			return
		}

		// 增加计数
		err = db.IncrRateLimit(context.Background(), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit increment failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PublicAPILimit 公开 API 限流
func (rl *RateLimiter) PublicAPILimit() gin.HandlerFunc {
	return rl.RateLimit(rl.config.RateLimit.PublicAPILimit)
}

// PrivateAPILimit 私有 API 限流
func (rl *RateLimiter) PrivateAPILimit() gin.HandlerFunc {
	return rl.RateLimit(rl.config.RateLimit.PrivateAPILimit)
}

// AuthAPILimit 认证 API 限流
func (rl *RateLimiter) AuthAPILimit() gin.HandlerFunc {
	return rl.RateLimit(rl.config.RateLimit.AuthAPILimit)
}

// getClientIdentifier 获取客户端标识
func getClientIdentifier(c *gin.Context) string {
	// 如果用户已登录，使用用户ID
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%v", userID)
	}

	// 否则使用IP地址
	return fmt.Sprintf("ip:%s", c.ClientIP())
}
