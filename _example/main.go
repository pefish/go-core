package main

import (
	"github.com/pefish/go-core/api"
	_type2 "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-core/api-strategy/type"
	global_api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	api_strategy "github.com/pefish/go-core/pkg/api-strategy"
	global_api_strategy3 "github.com/pefish/go-core/pkg/global-api-strategy"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-error"
	"log"
	"time"
)

func main() {
	service.Service.SetName(`test service`) // set service name
	service.Service.SetPath(`/api/test`)
	global_api_strategy.ParamValidateStrategy.SetErrorCode(2005)
	service.Service.SetRoutes([]*api.Api{
		{
			Description: "this is a test api",
			Path:        "/v1/test_api",
			Method:      `POST`,
			Strategies: []_type.StrategyData{
				{
					Strategy: &api_strategy.IpFilterStrategy,
					Param: api_strategy.IpFilterParam{
						GetValidIp: func(apiSession _type2.IApiSession) []string {
							return []string{`127.0.0.1`}
						},
					},
					Disable: true,
				},
			},
			ParamType:  global_api_strategy.ALL_TYPE,
			Controller: func(apiSession _type2.IApiSession) (i interface{}, info *go_error.ErrorInfo) {
				return "haha", nil
			},
		},
	})
	global_api_strategy3.GlobalRateLimitStrategy.SetErrorCode(10000)
	global_api_strategy2.GlobalApiStrategyDriver.Register(global_api_strategy2.GlobalStrategyData{
		Strategy: &global_api_strategy3.GlobalRateLimitStrategy,
		Param:    global_api_strategy3.GlobalRateLimitStrategyParam{
			FillInterval: 1000 * time.Millisecond,
		},
		Disable:  false,
	})
	service.Service.SetPort(3000)

	err := service.Service.Run()
	if err != nil {
		log.Fatal(err)
	}
}

