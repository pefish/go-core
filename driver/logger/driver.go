package logger

import go_interface_logger "github.com/pefish/go-interface-logger"

type LoggerDriver struct {
	Logger go_interface_logger.InterfaceLogger
}

var LoggerDriverInstance = LoggerDriver{
	Logger: go_interface_logger.DefaultLogger,
}

func (loggerDriver *LoggerDriver) Startup() {

}

func (loggerDriver *LoggerDriver) Register(logger go_interface_logger.InterfaceLogger) bool {
	loggerDriver.Logger = logger
	return true
}
