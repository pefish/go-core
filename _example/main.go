package main

import (
	"fmt"
	"github.com/pefish/go-config"
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

	go_config.Config.LoadYamlConfig(go_config.Configuration{})

	loggerInstance := go_logger.Log4goClass{}
	go_logger.Logger.Init(&loggerInstance, service.TestService.GetName(), `debug`)

	go_mysql.MysqlHelper.ConnectWithMap(go_config.Config.GetMap(`mysql`))

	service.TestService.Init().SetHealthyCheck(nil)
	service.TestService.SetPort(go_config.Config.GetUint64(`port`))
	service.TestService.Run()
}
