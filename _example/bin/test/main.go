package main

import (
	go_core "github.com/pefish/go-core"
	"github.com/pefish/go-core/api"
	api_session "github.com/pefish/go-core/api-session"
	api_strategy "github.com/pefish/go-core/api-strategy"
	global_api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	go_logger "github.com/pefish/go-logger"
	"time"
)

func PostTest(apiSession *api_session.ApiSessionClass) interface{} {
	return `haha, this is return value`
}

func main() {
	go_core.Service.SetName(`test service`) // set service name
	logger.LoggerDriver.Register(go_logger.NewLogger(go_logger.WithIsDebug(true))) // register logger
	go_core.Service.SetPath(`/api/test`)
	global_api_strategy.ParamValidateStrategy.SetErrorCode(2005)
	go_core.Service.SetRoutes([]*api.Api{
		{
			Description: "this is a test api",
			Path:        "/v1/test_api",
			Method:      `POST`,
			Strategies: []api_strategy.StrategyData{
				{
					Strategy: &api_strategy.IpFilterStrategy,
					Param: api_strategy.IpFilterParam{
						GetValidIp: func(apiSession *api_session.ApiSessionClass) []string {
							return []string{`127.0.0.1`}
						},
					},
					Disable: true,
				},
			},
			ParamType:  global_api_strategy.ALL_TYPE,
			Controller: PostTest,
		},
	})
	global_api_strategy.GlobalRateLimitStrategy.SetErrorCode(10000)
	global_api_strategy2.GlobalApiStrategyDriver.Register(global_api_strategy2.GlobalStrategyData{
		Strategy: &global_api_strategy.GlobalRateLimitStrategy,
		Param:    global_api_strategy.GlobalRateLimitStrategyParam{
			FillInterval: 1000 * time.Millisecond,
		},
		Disable:  false,
	})
	go_core.Service.SetPort(8080)

	go_core.Service.Run()
}

