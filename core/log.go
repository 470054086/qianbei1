package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// 日志类
type Log struct {
	logger       *logrus.Logger
	entry        *logrus.Entry
	loggerStatic map[string]*logrus.Entry
}

// 日志格式
type LogFormat string

const (
	JsonFormat   LogFormat = "json"
	NormalFormat LogFormat = "normal"
)
var g_log *logrus.Entry
// 创建类
func init()  {
	l := &Log{}
	dir := fmt.Sprintf("%slog.log", Config().Log.Dir)
	file, err := l.createLogFile(dir)
	if err != nil {
		panic(fmt.Sprint("create is error %v", err))
	}
	logger := logrus.New()
	// 判断需要的格式
	format := Config().Log.Format
	if LogFormat(format) == JsonFormat {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	if os.Getenv("mode") == "local" {
		logger.SetOutput(file)
	} else {
		logger.SetOutput(os.Stdout)
	}

	// todo 这里可以尝试加上requestId等其他的公共的属性
	fields := logger.WithFields(logrus.Fields{
		"ip":     "127.0.0.1",
		"author": "xx@qq.com",
	})
	l.entry = fields
	l.logger = logger
	l.loggerStatic = make(map[string]*logrus.Entry)
	g_log = fields
}
func QLog() *logrus.Entry {
	return g_log
}


func (l *Log) createLogFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 066)
	if err != nil {
		return nil, err
	}
	return file, err
}

// 使用不同的路径生成日志
// 单例模式
// 感觉没什么用呢
func (l *Log) Channel(channel string) *logrus.Entry {
	if _, ok := l.loggerStatic[channel]; !ok {
		logger := logrus.New()
		dir := fmt.Sprintf("%s%s.log", Config().Log.Dir, channel)
		file, _ := l.createLogFile(dir)
		format := Config().Log.Format
		if LogFormat(format) == JsonFormat {
			logger.SetFormatter(&logrus.JSONFormatter{})
		}
		logger.SetOutput(file)
		// 这里可以尝试加上requestId等其他的公共的属性
		fields := logger.WithFields(logrus.Fields{
			"ip":     "127.0.0.1",
			"author": "xx@qq.com",
		})
		l.loggerStatic[channel] = fields
	}
	return l.loggerStatic[channel]
}
