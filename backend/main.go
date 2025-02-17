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

	// 创建 Gin 实例
	r := gin.Default()

	// 配置 CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(r, logger)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Infof("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
