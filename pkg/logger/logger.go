package logger

import (
	"github.com/sirupsen/logrus"
)

// Log 是全局日志记录器
var Log *logrus.Logger

// Init 初始化日志配置，可以设置日志级别等
func Init(level string) {
	Log = logrus.New()
	// 设置日志格式：文本格式，包含完整时间戳
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	// 解析并设置日志级别，默认使用 Info 级别
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	Log.SetLevel(lvl)
}
