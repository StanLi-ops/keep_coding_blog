package config

import (
	"net/http"
	"os"
	"time"
)

// GetConfig 获取配置
func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080", // 端口
			TLS: TLSConfig{
				Enable:   false, // 是否启用TLS
				CertFile: "",    // 证书文件
				KeyFile:  "",    // 密钥文件
			},
		},
		Database: DatabaseConfig{
			Host:         "localhost", // 主机
			Port:         "5432",      // 端口
			User:         "postgres",  // 用户
			Password:     "1",         // 密码
			DBName:       "postgres",  // 数据库名称
			MaxOpenConns: 25,          // 最大打开连接数
			MaxIdleConns: 5,           // 最大空闲连接数
			SSLMode:      "disable",   // SSL模式
		},
		Redis: RedisConfig{
			Host:         "192.168.1.88",      // 主机
			Port:         "6379",              // 端口
			Password:     "1",                 // 密码
			DB:           0,                   // 数据库
			RBACPrefix:   "user_permissions:", // RBAC前缀
			RBACCacheTTL: 30 * time.Minute,    // RBAC缓存过期时间
		},
		JWT: JWTConfig{
			AccessTokenSecret:  getEnvOrDefault("JWT_ACCESS_SECRET", "your-access-secret-key"),   // 访问令牌密钥
			RefreshTokenSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key"), // 刷新令牌密钥
			AccessTokenTTL:     15 * time.Minute,                                                 // 访问令牌15分钟过期
			RefreshTokenTTL:    7 * 24 * time.Hour,                                               // 刷新令牌7天过期
		},
		RateLimit: RateLimitConfig{
			PublicAPILimit:  100,         // 100次/分钟
			PrivateAPILimit: 60,          // 60次/分钟
			AuthAPILimit:    5,           // 5次/分钟
			Duration:        time.Minute, // 1分钟时间窗口
		},
		CORS: CORSConfig{
			AllowOrigins: []string{
				"http://localhost:8080",              // 开发环境
				"https://your-production-domain.com", // 生产环境
			},
			AllowMethods: []string{
				http.MethodGet,     // GET方法
				http.MethodPost,    // POST方法
				http.MethodPut,     // PUT方法
				http.MethodDelete,  // DELETE方法
				http.MethodOptions, // OPTIONS方法
			},
			AllowHeaders: []string{
				"Origin",          // 来源
				"Content-Type",    // 内容类型
				"Content-Length",  // 内容长度
				"Accept-Encoding", // 接受编码
				"X-CSRF-Token",    // CSRF令牌
				"Authorization",   // 授权
			},
			ExposeHeaders: []string{
				"Content-Length",               // 内容长度
				"Access-Control-Allow-Origin",  // 允许来源
				"Access-Control-Allow-Headers", // 允许头
			},
			AllowCredentials: true,  // 允许凭证
			MaxAge:           86400, // 24小时
		},
		Security: SecurityConfig{
			EncryptionKey: getEnvOrDefault("ENCRYPTION_KEY", "12345678901234567890123456789012"), // 加密密钥
		},
	}
}

// Config
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
	Security  SecurityConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	TLS  TLSConfig
}

type TLSConfig struct {
	Enable   bool
	CertFile string
	KeyFile  string
}

// DatabaseConfig 数据库配置
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

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	RBACPrefix   string
	RBACCacheTTL time.Duration
	RatePrefix   string
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	PublicAPILimit  int
	PrivateAPILimit int
	AuthAPILimit    int
	Duration        time.Duration
}

// JWTConfig JWT配置
type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EncryptionKey string
}

// 获取环境变量，如果没有则使用默认值。
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
