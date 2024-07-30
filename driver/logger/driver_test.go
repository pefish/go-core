package logger

import (
	i_logger "github.com/pefish/go-interface/i-logger"
)

func ExampleLoggerDriver_Register() {
	LoggerDriverInstance.Register(&i_logger.DefaultLogger)

	LoggerDriverInstance.Startup()

	LoggerDriverInstance.Logger.Info("haha")

	// Output:
	// [INFO] haha
}
