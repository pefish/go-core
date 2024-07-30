package logger

import (
	i_logger "github.com/pefish/go-interface/i-logger"
)

type LoggerDriver struct {
	Logger i_logger.ILogger
}

var LoggerDriverInstance = LoggerDriver{
	Logger: &i_logger.DefaultLogger,
}

func (loggerDriver *LoggerDriver) Startup() {

}

func (loggerDriver *LoggerDriver) Register(logger i_logger.ILogger) bool {
	loggerDriver.Logger = logger
	return true
}
