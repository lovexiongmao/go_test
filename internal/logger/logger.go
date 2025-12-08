package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go_test/internal/config"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(cfg *config.Config) *Logger {
	log := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// 设置日志格式
	if cfg.Log.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置输出
	output := strings.ToLower(cfg.Log.Output)
	writers := []io.Writer{}

	// 始终允许 stdout
	if output == "stdout" || output == "both" || output == "" {
		writers = append(writers, os.Stdout)
	}

	// 可选的文件输出
	if output == "file" || output == "both" {
		// 确保目录存在
		if cfg.Log.File == "" {
			cfg.Log.File = "logs/app.log"
		}
		if err := os.MkdirAll(filepath.Dir(cfg.Log.File), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "创建日志目录失败: %v\n", err)
		} else {
			f, ferr := os.OpenFile(cfg.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
			if ferr != nil {
				fmt.Fprintf(os.Stderr, "打开日志文件失败，改为stdout: %v\n", ferr)
				writers = append(writers, os.Stdout)
			} else {
				writers = append(writers, f)
			}
		}
	}

	// 如果 writers 为空，兜底用 stdout
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	log.SetOutput(io.MultiWriter(writers...))

	return &Logger{Logger: log}
}

// 提供便捷方法
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}
