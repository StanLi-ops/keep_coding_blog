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
func SetUserPermissions(ctx context.Context, userID uint, permissions []models.Permission) error {
	// 序列化权限数据
	permissionsData, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	// 设置缓存,带过期时间
	key := fmt.Sprintf("%s%d", config.GetConfig().Redis.RBACPrefix, userID)
	return RedisClient.Set(ctx, key, permissionsData, config.GetConfig().Redis.RBACCacheTTL).Err()
}

// GetUserPermissions 获取用户权限缓存
func GetUserPermissions(ctx context.Context, userID uint) ([]models.Permission, error) {
	key := fmt.Sprintf("%s%d", config.GetConfig().Redis.RBACPrefix, userID)
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
func DeleteUserPermissions(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("%s%d", config.GetConfig().Redis.RBACPrefix, userID)
	return RedisClient.Del(ctx, key).Err()
}
