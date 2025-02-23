package middleware

import (
	"keep_learning_blog/config"
	"keep_learning_blog/db"
	"keep_learning_blog/models"
	"keep_learning_blog/service"
	"keep_learning_blog/utils/logger"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// userService 初始化
var userService *service.UserService

// RBACAuth RBAC权限控制中间件
func RBACAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取当前登录用户
		userID, exists := c.Get("user_id")
		if !exists {
			logger.Log.Warn("User not logged in or login expired")
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "not logged in or login has expired",
			})
			c.Abort()
			return
		}

		// 尝试从Redis获取权限
		permissions, err := db.GetUserPermissions(c, userID.(uint), cfg)
		if err != nil {
			logger.Log.WithError(err).Error("Failed to get permissions from Redis")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permissions"})
			c.Abort()
			return
		}

		// 如果Redis中没有,从数据库获取并缓存
		if permissions == nil {
			permissions, err = userService.GetUserPermissions(userID.(uint))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to get user permissions",
				})
				c.Abort()
				return
			}

			// 缓存到Redis
			if err := db.SetUserPermissions(c, userID.(uint), permissions, cfg); err != nil {
				// 仅记录日志,不中断请求
				log.Printf("Failed to cache permissions: %v", err)
			}
		}

		// 获取路径并去掉 /api 前缀
		fullPath := c.Request.URL.Path
		path := strings.TrimPrefix(fullPath, "/api")
		method := c.Request.Method

		// 检查权限
		if !checkPermission(permissions, method, path) {
			logger.Log.WithFields(logger.Fields(map[string]interface{}{
				"user_id": userID,
				"path":    path,
				"method":  method,
			})).Warn("Permission denied")

			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "no permission to access",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// checkPermission 检查是否有权限访问
func checkPermission(permissions []models.Permission, method, path string) bool {
	for _, permission := range permissions {
		if permission.Method == method && matchPath(permission.Path, path) {
			return true
		}
	}
	return false
}

// matchPath 检查请求路径是否匹配权限路径模式
func matchPath(pattern, path string) bool {
	// 将模式和路径按/分割
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// 如果路径部分比模式部分少，直接返回false
	if len(pathParts) < len(patternParts) {
		return false
	}

	// 逐段匹配
	for i := 0; i < len(patternParts); i++ {
		// 如果是参数部分（以:开头），则跳过检查
		if strings.HasPrefix(patternParts[i], ":") {
			continue
		}
		// 如果当前段不匹配，返回false
		if patternParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}
