package main

import (
	"fmt"
	"keep_learning_blog/config"
	"keep_learning_blog/db"
	"keep_learning_blog/routes"
	"os"

	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
)

// main 主函数
func main() {
	// 获取配置
	cfg := config.GetConfig()

	// 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		fmt.Printf("Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	// 使用全局日志实例
	log := logger.Log

	// 初始化数据库
	if err := db.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化Redis
	if err := db.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 创建 Gin 实例
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r, cfg)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Infof("Server starting on %s", serverAddr)
	if cfg.Server.TLS.Enable {
		if err := r.RunTLS(serverAddr,
			cfg.Server.TLS.CertFile,
			cfg.Server.TLS.KeyFile); err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		if err := r.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
