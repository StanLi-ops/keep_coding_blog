package config

// Config 配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
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
	}
}
