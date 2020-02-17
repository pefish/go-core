package logger

type LoggerDriverClass struct {
	Logger InterfaceLogger
}

var LoggerDriver = LoggerDriverClass{}

func (this *LoggerDriverClass) Startup() {

}

func (this *LoggerDriverClass) Register(logger InterfaceLogger) bool {
	this.Logger = logger
	return true
}
