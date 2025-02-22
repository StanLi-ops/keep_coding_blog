package main

import (
	"fmt"
	"keep_coding_blog/config"
	"keep_coding_blog/db"
	"keep_coding_blog/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// setupLogger 设置日志
func setupLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// 根据环境配置日志级别
	if os.Getenv("ENV") == "production" {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}

	return log
}

// main 主函数
func main() {
	// 设置日志
	logger := setupLogger()

	// 加载配置
	cfg := config.GetConfig()

	// 初始化数据库
	if err := db.InitDB(cfg, logger); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化Redis
	if err := db.InitRedis(cfg); err != nil {
		logger.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 创建 Gin 实例
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r, logger, cfg)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Infof("Server starting on %s", serverAddr)
	if cfg.Server.TLS.Enable {
		if err := r.RunTLS(serverAddr,
			cfg.Server.TLS.CertFile,
			cfg.Server.TLS.KeyFile); err != nil {
			logger.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		if err := r.Run(serverAddr); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}
}
