package logger

import (
	go_logger "github.com/pefish/go-logger"
)

func ExampleLoggerDriver_Register() {
	LoggerDriverInstance.Register(go_logger.Logger)

	LoggerDriverInstance.Startup()

	LoggerDriverInstance.Logger.Info("haha")

	// Output:
	// [INFO] haha
}
