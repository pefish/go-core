package main

import (
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-config"
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-logger"
	"log"
	"os"
	"runtime/debug"
	"test/route"
	//"contrib.go.opencensus.io/exporter/stackdriver"
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

	//sd, err := stackdriver.NewExporter(stackdriver.Options{
	//	ProjectID: "demo-project-id",
	//	// MetricPrefix helps uniquely identify your metrics.
	//	MetricPrefix: "demo-prefix",
	//	// ReportingInterval sets the frequency of reporting metrics
	//	// to stackdriver backend.
	//	ReportingInterval: 60 * time.Second,
	//})
	//if err != nil {
	//	log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	//}
	//// It is imperative to invoke flush before your main function exits
	//defer sd.Flush()
	//
	//// Register it as a trace exporter
	//trace.RegisterExporter(sd)

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
