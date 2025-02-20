package config

import (
	"os"
	"time"
)

// Config 配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Port string
}

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
	SSLMode      string
}

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	RBACPrefix   string
	RBACCacheTTL time.Duration
}

// JWTConfig JWT配置结构体
type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

// GetConfig 获取配置
func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         "5432",
			User:         "postgres",
			Password:     "1",
			DBName:       "postgres",
			MaxOpenConns: 25,
			MaxIdleConns: 5,
			SSLMode:      "disable",
		},
		Redis: RedisConfig{
			Host:         "192.168.1.88",
			Port:         "6379",
			Password:     "1",
			DB:           0,
			RBACPrefix:   "user_permissions:",
			RBACCacheTTL: 30 * time.Minute,
		},
		JWT: JWTConfig{
			AccessTokenSecret:  getEnvOrDefault("JWT_ACCESS_SECRET", "your-access-secret-key"),
			RefreshTokenSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			AccessTokenTTL:     15 * time.Minute,   // 访问令牌15分钟过期
			RefreshTokenTTL:    7 * 24 * time.Hour, // 刷新令牌7天过期
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
