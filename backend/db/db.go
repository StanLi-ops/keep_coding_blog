package db

import (
	"fmt"
	"keep_coding_blog/config"
	"keep_coding_blog/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(cfg *config.Config, logger *logrus.Logger) error {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Error("Failed to connect to database")
		return err
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Comment{})
	if err != nil {
		logger.WithError(err).Error("Failed to migrate database")
		return err
	}

	DB = db
	logger.Info("Database connected and migrated successfully")
	return nil
}

// GetDB 返回数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
