package logger

import (
	"keep_learning_blog/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log      *logrus.Logger // 系统日志
	AuditLog *logrus.Logger // 审计日志
)

// InitLogger 初始化日志
func InitLogger(cfg *config.Config) error {
	// 初始化系统日志
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(&lumberjack.Logger{
		Filename:   cfg.SystemLog.FilePath,
		MaxSize:    cfg.SystemLog.MaxSize,
		MaxBackups: cfg.SystemLog.MaxBackups,
		MaxAge:     cfg.SystemLog.MaxAge,
		Compress:   cfg.SystemLog.Compress,
	})

	// 初始化审计日志
	AuditLog = logrus.New()
	AuditLog.SetFormatter(&logrus.JSONFormatter{})
	AuditLog.SetOutput(&lumberjack.Logger{
		Filename:   cfg.AuditLog.Filename,
		MaxSize:    cfg.AuditLog.MaxSize,
		MaxBackups: cfg.AuditLog.MaxBackups,
		MaxAge:     cfg.AuditLog.MaxAge,
		Compress:   cfg.AuditLog.Compress,
	})

	return nil
}

// Fields 创建日志字段
func Fields(fields map[string]interface{}) logrus.Fields {
	return logrus.Fields(fields)
}
