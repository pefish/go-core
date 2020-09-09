package logger

import (
	go_interface_logger "github.com/pefish/go-interface-logger"
)

func ExampleLoggerDriver_Register() {
	LoggerDriverInstance.Register(go_interface_logger.DefaultLogger)

	LoggerDriverInstance.Startup()

	LoggerDriverInstance.Logger.Info("haha")

	// Output:
	// [INFO] haha
}