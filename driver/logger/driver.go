package logger

import (
	go_logger "github.com/pefish/go-logger"
)

type LoggerDriver struct {
	Logger go_logger.InterfaceLogger
}

var LoggerDriverInstance = LoggerDriver{
	Logger: go_logger.Logger,
}

func (loggerDriver *LoggerDriver) Startup() {

}

func (loggerDriver *LoggerDriver) Register(logger go_logger.InterfaceLogger) bool {
	loggerDriver.Logger = logger
	return true
}
