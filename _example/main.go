package main

import (
	"contrib.go.opencensus.io/exporter/stackdriver"
	"fmt"
	go_application "github.com/pefish/go-application"
	go_config "github.com/pefish/go-config"
	go_core "github.com/pefish/go-core"
	api_strategy "github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/logger"
	go_logger "github.com/pefish/go-logger"
	"log"
	"os"
	"runtime/debug"
	"test/route"
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

	go_core.Service.SetName(`测试服务api`)
	go_core.Service.AddGlobalStrategy(&api_strategy.OpenCensusStrategy, api_strategy.OpenCensusStrategyParam{
		StackDriverOption: &stackdriver.Options{
			ProjectID:    `pefish`,
			MetricPrefix: `test`,
		}})

	go_logger.Logger.Init(go_core.Service.GetName(), `debug`)
	logger.LoggerDriver.Register(go_logger.Logger)

	//go_mysql.MysqlHelper.ConnectWithMap(go_config.Config.MustGetMap(`mysql`))

	go_core.Service.SetPath(`/api/test`)
	api_strategy.RateLimitApiStrategy.SetErrorCode(2006)
	api_strategy.ParamValidateStrategy.SetErrorCode(2005)
	api_strategy.IpFilterStrategy.SetErrorCode(2007)
	api_strategy.CorsApiStrategy.SetAllowedOrigins([]string{`*`})
	go_core.Service.SetRoutes(route.TestRoute)
	go_core.Service.SetPort(go_config.Config.GetUint64(`port`))



	go_core.Service.Run()
}
