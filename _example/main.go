package main

import (
	"fmt"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/config"
	"github.com/pefish/go-mysql"
	"log"
	"os"
	"runtime/debug"
	"test/src/controllers"
	"test/src/global"
	"test/src/service"
)

func main() {
	defer func() {
		if global.MysqlHelper != nil {
			global.MysqlHelper.Close()
		}

		if err := recover(); err != nil {
			log.Println(err)
			fmt.Println(string(debug.Stack()))
			os.Exit(1)
		}
		os.Exit(0)
	}()

	config.Config.LoadJsonConfig(config.Configuration{})

	mysqlConfig := config.Config.GetMap(`mysql`)
	global.MysqlHelper = &p_mysql.MysqlClass{}
	global.MysqlHelper.ConnectWithConfiguration(p_mysql.Configuration{
		Host: mysqlConfig[`host`].(string),
		Port: 3306,
		Username: mysqlConfig[`username`].(string),
		Password: mysqlConfig[`password`].(string),
		Database: mysqlConfig[`database`].(string),
	})

	service.TestService.Init(map[string]interface{}{
		`apiControllers`: map[string]api_session.ApiHandlerType{
			`test_api`: controllers.TestController.Test,
		},
	}).SetHealthyCheck(nil)
	service.TestService.SetHost(config.Config.GetString(`host`))
	service.TestService.SetPort(config.Config.GetUint64(`port`))
	service.TestService.Run()
}
