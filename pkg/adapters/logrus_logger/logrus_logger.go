package logrus_logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	instance := logrus.New()
	instance.SetFormatter(&logrus.TextFormatter{})
	instance.SetLevel(logrus.DebugLevel)
	instance.SetOutput(os.Stdout)

	return &LogrusLogger{
		logger: instance,
	}
}

func (l *LogrusLogger) Passthrough(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		fmt.Println(msg)
	} else {
		fmt.Printf(msg, keysAndValues...)
	}
}

func (l *LogrusLogger) Debug(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		l.logger.Debug(msg)
	} else {
		l.logger.Debugf(msg, keysAndValues...)
	}
}

func (l *LogrusLogger) Info(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		l.logger.Info(msg)
	} else {
		l.logger.Infof(msg, keysAndValues...)
	}
}

func (l *LogrusLogger) Warning(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		l.logger.Warning(msg)
	} else {
		l.logger.Warningf(msg, keysAndValues...)
	}
}

func (l *LogrusLogger) Error(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		l.logger.Error(msg)
	} else {
		l.logger.Errorf(msg, keysAndValues...)
	}
}

func (l *LogrusLogger) Fatal(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		l.logger.Fatal(msg)
	} else {
		l.logger.Fatalf(msg, keysAndValues...)
	}
}
