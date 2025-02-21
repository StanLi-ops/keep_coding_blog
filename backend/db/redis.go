package db

import (
	"context"
	"encoding/json"
	"fmt"
	"keep_coding_blog/config"
	"keep_coding_blog/models"
	"time"

	"github.com/redis/go-redis/v9"
)

// 登录限制相关的常量
const (
	LoginAttemptsPrefix = "login_attempts:"
	LoginLockPrefix     = "login_lock:"
)

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

	_, err := RedisClient.Ping(ctx).Result()
	return err
}

// AddToBlacklist 将token加入黑名单
func AddToBlacklist(ctx context.Context, tokenID string, expiration time.Duration) error {
	return RedisClient.Set(ctx, fmt.Sprintf("blacklist:%s", tokenID), true, expiration).Err()
}

// IsBlacklisted 检查token是否在黑名单中
func IsBlacklisted(ctx context.Context, tokenID string) bool {
	exists, _ := RedisClient.Exists(ctx, fmt.Sprintf("blacklist:%s", tokenID)).Result()
	return exists > 0
}

// SetUserPermissions 缓存用户权限
func SetUserPermissions(ctx context.Context, userID uint, permissions []models.Permission, cfg *config.Config) error {
	// 序列化权限数据
	permissionsData, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	// 设置缓存,带过期时间
	key := fmt.Sprintf("%s%d", cfg.Redis.RBACPrefix, userID)
	return RedisClient.Set(ctx, key, permissionsData, cfg.Redis.RBACCacheTTL).Err()
}

// GetUserPermissions 获取用户权限缓存
func GetUserPermissions(ctx context.Context, userID uint, cfg *config.Config) ([]models.Permission, error) {
	key := fmt.Sprintf("%s%d", cfg.Redis.RBACPrefix, userID)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var permissions []models.Permission
	if err := json.Unmarshal(data, &permissions); err != nil {
		return nil, err
	}
	return permissions, nil
}

// DeleteUserPermissions 删除用户权限缓存
func DeleteUserPermissions(ctx context.Context, userID uint, cfg *config.Config) error {
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

	if success {
		// 登录成功，清除失败记录
		return RedisClient.Del(ctx, key).Err()
	}

	// 登录失败，增加计数
	attempts, err := RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	// 设置键的过期时间（如果还没有设置）
	RedisClient.Expire(ctx, key, lockDuration)

	// 如果达到最大尝试次数，设置锁定状态
	if attempts >= int64(maxAttempts) {
		lockKey := fmt.Sprintf("%s%s", LoginLockPrefix, identifier)
		return RedisClient.Set(ctx, lockKey, "locked", lockDuration).Err()
	}

	return nil
}
