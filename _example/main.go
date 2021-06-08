package main

import (
	api_strategy "github.com/pefish/go-core-strategy/api-strategy"
	global_api_strategy3 "github.com/pefish/go-core-strategy/global-api-strategy"
	"github.com/pefish/go-core/api"
	_type2 "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-core/api-strategy/type"
	global_api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-error"
	"log"
	"time"
)

func main() {
	service.Service.SetName(`test service`) // set service name
	service.Service.SetPath(`/api/test`)
	global_api_strategy.ParamValidateStrategyInstance.SetErrorCode(2005)
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
			ParamType: global_api_strategy.ALL_TYPE,
			Controller: func(apiSession _type2.IApiSession) (i interface{}, info *go_error.ErrorInfo) {
				var params struct {
					Test string `json:"test" validate:"is-mobile"`
				}
				apiSession.ScanParams(&params)
				//return nil, go_error.WrapWithAll(errors.New("haha"), 2000, map[string]interface{}{
				//	"haha": "u7ytu7",
				//})
				return params.Test, nil
			},
			Params: struct {
				Test string `json:"test" validate:"required,is-mobile"`
			}{},
		},
	})
	global_api_strategy3.GlobalRateLimitStrategy.SetErrorCode(10000)
	global_api_strategy2.GlobalApiStrategyDriverInstance.Register(global_api_strategy2.GlobalStrategyData{
		Strategy: &global_api_strategy3.GlobalRateLimitStrategy,
		Param: global_api_strategy3.GlobalRateLimitStrategyParam{
			FillInterval: 1000 * time.Millisecond,
		},
		Disable: false,
	})
	service.Service.SetPort(8080)

	err := service.Service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
