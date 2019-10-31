package logger

type LoggerDriverClass struct {
	logger InterfaceLogger
}

var LoggerDriver = LoggerDriverClass{}

func (this *LoggerDriverClass) Startup() {

}

func (this *LoggerDriverClass) Register(logger InterfaceLogger) bool {
	this.logger = logger
	return true
}

func (this *LoggerDriverClass) Debug(args ...interface{}) {
	this.logger.Debug(args...)
}
func (this *LoggerDriverClass) DebugF(format string, args ...interface{}) {
	this.logger.DebugF(format, args...)
}
func (this *LoggerDriverClass) Info(args ...interface{}) {
	this.logger.Info(args...)
}
func (this *LoggerDriverClass) InfoF(format string, args ...interface{}) {
	this.logger.InfoF(format, args...)
}
func (this *LoggerDriverClass) Warn(args ...interface{}) {
	this.logger.Warn(args...)
}
func (this *LoggerDriverClass) WarnF(format string, args ...interface{}) {
	this.logger.WarnF(format, args...)
}
func (this *LoggerDriverClass) Error(args ...interface{}) {
	this.logger.Error(args...)
}
func (this *LoggerDriverClass) ErrorF(format string, args ...interface{}) {
	this.logger.ErrorF(format, args...)
}
