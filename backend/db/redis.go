package db

import (
	"context"
	"encoding/json"
	"fmt"
	"keep_learning_blog/config"
	"keep_learning_blog/models"
	"keep_learning_blog/utils/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

// 登录限制相关的常量
const (
	LoginAttemptsPrefix = "login_attempts:"
	LoginLockPrefix     = "login_lock:"
)

// RedisClient 全局Redis客户端
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		logger.Log.WithError(err).Error("Failed to connect to Redis")
		return err
	}

	logger.Log.Info("Redis connected successfully")
	return nil
}

// AddToBlacklist 将token加入黑名单
func AddToBlacklist(ctx context.Context, tokenID string, expiration time.Duration) error {
	// 将token加入黑名单
	err := RedisClient.Set(ctx, fmt.Sprintf("blacklist:%s", tokenID), true, expiration).Err()
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"token_id": tokenID,
			"error":    err,
		})).Error("Failed to add token to blacklist")
		return err
	}

	logger.Log.WithField("token_id", tokenID).Info("Token added to blacklist")
	return nil
}

// IsBlacklisted 检查token是否在黑名单中
func IsBlacklisted(ctx context.Context, tokenID string) bool {
	// 检查token是否在黑名单中
	exists, _ := RedisClient.Exists(ctx, fmt.Sprintf("blacklist:%s", tokenID)).Result()
	return exists > 0
}

// SetUserPermissions 缓存用户权限
func SetUserPermissions(ctx context.Context, userID uint, permissions []models.Permission, cfg *config.Config) error {
	// 缓存用户权限
	permissionsData, err := json.Marshal(permissions)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to marshal permissions")
		return err
	}

	// 设置缓存
	key := fmt.Sprintf("%s%d", cfg.Redis.RBACPrefix, userID)
	if err := RedisClient.Set(ctx, key, permissionsData, cfg.Redis.RBACCacheTTL).Err(); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"user_id": userID,
			"error":   err,
		})).Error("Failed to cache user permissions")
		return err
	}

	logger.Log.WithField("user_id", userID).Info("User permissions cached successfully")
	return nil
}

// GetUserPermissions 获取用户权限缓存
func GetUserPermissions(ctx context.Context, userID uint, cfg *config.Config) ([]models.Permission, error) {
	// 获取用户权限缓存
	key := fmt.Sprintf("%s%d", cfg.Redis.RBACPrefix, userID)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			logger.Log.WithField("user_id", userID).Debug("No cached permissions found")
			return nil, nil
		}
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"user_id": userID,
			"error":   err,
		})).Error("Failed to get cached permissions")
		return nil, err
	}

	// 解析权限数据
	var permissions []models.Permission
	if err := json.Unmarshal(data, &permissions); err != nil {
		logger.Log.WithError(err).Error("Failed to unmarshal permissions")
		return nil, err
	}

	return permissions, nil
}

// DeleteUserPermissions 删除用户权限缓存
func DeleteUserPermissions(ctx context.Context, userID uint, cfg *config.Config) error {
	// 删除用户权限缓存
	key := fmt.Sprintf("%s%d", cfg.Redis.RBACPrefix, userID)
	return RedisClient.Del(ctx, key).Err()
}

// GetRateLimit 获取请求计数
func GetRateLimit(ctx context.Context, key string) (int, error) {
	count, err := RedisClient.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

// SetRateLimit 设置初始请求计数和过期时间
func SetRateLimit(ctx context.Context, key string, duration time.Duration) error {
	return RedisClient.Set(ctx, key, 1, duration).Err()
}

// IncrRateLimit 增加请求计数
func IncrRateLimit(ctx context.Context, key string) error {
	return RedisClient.Incr(ctx, key).Err()
}

// IsLoginLocked 检查用户是否被锁定
func IsLoginLocked(ctx context.Context, identifier string) bool {
	lockKey := fmt.Sprintf("%s%s", LoginLockPrefix, identifier)
	exists, _ := RedisClient.Exists(ctx, lockKey).Result()
	return exists > 0
}

// GetLoginLockRemainingTime 获取锁定剩余时间
func GetLoginLockRemainingTime(ctx context.Context, identifier string) time.Duration {
	lockKey := fmt.Sprintf("%s%s", LoginLockPrefix, identifier)
	ttl, _ := RedisClient.TTL(ctx, lockKey).Result()
	return ttl
}

// RecordLoginAttempt 记录登录尝试
func RecordLoginAttempt(ctx context.Context, identifier string, success bool, maxAttempts int, lockDuration time.Duration) error {
	key := fmt.Sprintf("%s%s", LoginAttemptsPrefix, identifier)

	// 如果登录成功，则清除登录尝试
	if success {
		if err := RedisClient.Del(ctx, key).Err(); err != nil {
			logger.Log.WithFields(logger.Fields(map[string]interface{}{
				"identifier": identifier,
				"error":      err,
			})).Error("Failed to clear login attempts")
			return err
		}
		logger.Log.WithField("identifier", identifier).Info("Login attempts cleared after successful login")
		return nil
	}

	// 增加登录尝试次数
	attempts, err := RedisClient.Incr(ctx, key).Result()
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"identifier": identifier,
			"error":      err,
		})).Error("Failed to record login attempt")
		return err
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"identifier": identifier,
		"attempts":   attempts,
	})).Info("Login attempt recorded")

	// 如果尝试次数达到上限，则锁定账户
	if attempts >= int64(maxAttempts) {
		lockKey := fmt.Sprintf("%s%s", LoginLockPrefix, identifier)
		if err := RedisClient.Set(ctx, lockKey, "locked", lockDuration).Err(); err != nil {
			logger.Log.WithFields(logger.Fields(map[string]interface{}{
				"identifier": identifier,
				"error":      err,
			})).Error("Failed to set account lock")
			return err
		}
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"identifier":    identifier,
			"lock_duration": lockDuration,
		})).Warning("Account locked due to too many failed attempts")
	}

	return nil
}
