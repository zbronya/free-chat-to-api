package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	once     sync.Once
	instance *LogrusLogger
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type LogrusLogger struct {
	*logrus.Logger
}
