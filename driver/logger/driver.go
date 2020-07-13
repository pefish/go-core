package logger

import go_interface_logger "github.com/pefish/go-interface-logger"

type LoggerDriverClass struct {
	Logger go_interface_logger.InterfaceLogger
}

var LoggerDriver = LoggerDriverClass{
	Logger: go_interface_logger.DefaultLogger,
}

func (loggerDriver *LoggerDriverClass) Startup() {

}

func (loggerDriver *LoggerDriverClass) Register(logger go_interface_logger.InterfaceLogger) bool {
	loggerDriver.Logger = logger
	return true
}
