package db

import (
	"context"
	"fmt"
	"keep_coding_blog/config"
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
