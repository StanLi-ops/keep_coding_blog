package logger

import (
	"fmt"
	"io"
	"keep_learning_blog/config"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log 全局日志实例
	Log *logrus.Logger
)

// InitLogger 初始化日志
func InitLogger(config *config.LogConfig) error {
	Log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("parse log level error: %v", err)
	}
	Log.SetLevel(level)

	// 设置日志格式
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	// 启用调用者信息
	Log.SetReportCaller(true)

	// 设置输出
	var writers []io.Writer

	// 如果需要输出到文件
	if config.FilePath != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("create log directory error: %v", err)
		}

		// 配置日志轮转
		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,    // MB
			MaxBackups: config.MaxBackups, // 保留的旧文件个数
			MaxAge:     config.MaxAge,     // 天
			Compress:   config.Compress,   // 是否压缩
		}
		writers = append(writers, fileWriter)
	}

	// 如果需要输出到控制台
	if config.ConsoleOutput {
		writers = append(writers, os.Stdout)
	}

	// 设置多输出
	if len(writers) > 0 {
		Log.SetOutput(io.MultiWriter(writers...))
	}

	return nil
}

// Fields 创建日志字段
func Fields(fields map[string]interface{}) logrus.Fields {
	return logrus.Fields(fields)
}
