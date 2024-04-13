package app

import (
	"ghostel/pkg/adapters/logrus_logger"
	"ghostel/pkg/definitions"
)

var globalLoggerInstance definitions.ILogger

func GetGlobalLogger() definitions.ILogger {
	if globalLoggerInstance == nil {
		globalLoggerInstance = logrus_logger.NewLogrusLogger()
	}
	return globalLoggerInstance
}
