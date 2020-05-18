package logger

import go_interface_logger "github.com/pefish/go-interface-logger"

type LoggerDriverClass struct {
	Logger go_interface_logger.InterfaceLogger
}

var LoggerDriver = LoggerDriverClass{}

func (this *LoggerDriverClass) Startup() {

}

func (this *LoggerDriverClass) Register(logger go_interface_logger.InterfaceLogger) bool {
	this.Logger = logger
	return true
}
