package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	once.Do(func() {
		l := logrus.New()
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.DebugLevel)
		l.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		instance = &LogrusLogger{
			Logger: l,
		}
	})
}

func GetLogger() Logger {
	return instance
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}
