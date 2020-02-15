package main

import (
	"contrib.go.opencensus.io/exporter/stackdriver"
	"fmt"
	"go.opencensus.io/trace"
	"log"
	"os"
	"runtime/debug"
	"test/route"
	"time"

	go_application "github.com/pefish/go-application"
	go_config "github.com/pefish/go-config"
	api_strategy "github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-core/service"
	go_logger "github.com/pefish/go-logger"
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



	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "pefish",
		MetricPrefix: "test",
		ReportingInterval: 60 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	}
	defer sd.Flush()
	trace.RegisterExporter(sd)



	service.Service.SetName(`测试服务api`)

	go_logger.Logger.Init(service.Service.GetName(), `debug`)
	logger.LoggerDriver.Register(go_logger.Logger)

	//go_mysql.MysqlHelper.ConnectWithMap(go_config.Config.MustGetMap(`mysql`))

	service.Service.SetPath(`/api/test`)
	api_strategy.RateLimitApiStrategy.SetErrorCode(2006)
	api_strategy.ParamValidateStrategy.SetErrorCode(2005)
	api_strategy.IpFilterStrategy.SetErrorCode(2007)
	api_strategy.CorsApiStrategy.SetAllowedOrigins([]string{`*`})
	service.Service.SetRoutes(route.TestRoute)
	service.Service.SetPort(go_config.Config.GetUint64(`port`))

	service.Service.Run()
}
