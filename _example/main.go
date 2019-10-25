package main

import (
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-config"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-mysql"
	"log"
	"os"
	"runtime/debug"
	"test/service"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			fmt.Println(string(debug.Stack()))
			os.Exit(1)
		}
		os.Exit(0)
	}()


	go_application.Application.Debug = true
	go_config.Config.MustLoadYamlConfig(go_config.Configuration{
		ConfigEnvName: `GO_CONFIG`,
		SecretEnvName: `GO_SECRET`,
	})

	go_logger.Logger.Init(service.TestService.GetName(), `debug`)

	coreLogger := go_logger.LoggerClass{}
	coreLogger.Init(`core`, `debug`)
	logger.Logger = &coreLogger
	fmt.Printf(`%#v`, coreLogger)

	go_mysql.MysqlHelper.ConnectWithMap(go_config.Config.MustGetMap(`mysql`))

	service.TestService.Init()
	service.TestService.SetPort(go_config.Config.GetUint64(`port`))
	service.TestService.Run()
}
