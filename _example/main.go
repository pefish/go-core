package main

import (
	"fmt"
	"github.com/pefish/go-config"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-mysql"
	"log"
	"os"
	"runtime/debug"
	"test/src/controllers"
	"test/src/service"
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

	go_config.Config.LoadJsonConfig(go_config.Configuration{})
	fmt.Printf("%#v", go_config.Config.GetAll())

	mysqlConfig := go_config.Config.GetMap(`mysql`)
	go_mysql.MysqlHelper.ConnectWithConfiguration(go_mysql.Configuration{
		Host:     mysqlConfig[`host`].(string),
		Port:     3306,
		Username: mysqlConfig[`username`].(string),
		Password: mysqlConfig[`password`].(string),
		Database: mysqlConfig[`database`].(string),
	})

	service.TestService.Init(map[string]interface{}{
		`apiControllers`: map[string]api_session.ApiHandlerType{
			`test_api`: controllers.TestController.Test,
		},
	}).SetHealthyCheck(nil)
	service.TestService.SetPort(go_config.Config.GetUint64(`port`))
	service.TestService.Run()
}
