package driver

import "github.com/pefish/go-interface-logger"

type LoggerDriverClass struct {

}

var Logger go_interface_logger.InterfaceLogger

var LoggerDriver = LoggerDriverClass{}

func (this *LoggerDriverClass) Startup() {

}

func (this *LoggerDriverClass) Register(logger go_interface_logger.InterfaceLogger) bool {
	Logger = logger
	return true
}
