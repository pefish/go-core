package main

import (
	"context"
	api_strategy "github.com/pefish/go-core-strategy/api-strategy"
	global_api_strategy3 "github.com/pefish/go-core-strategy/global-api-strategy"
	_type2 "github.com/pefish/go-core-type/api-session"
	"github.com/pefish/go-core/api"
	"github.com/pefish/go-core/api-strategy/type"
	global_api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-error"
	go_logger "github.com/pefish/go-logger"
	task_driver "github.com/pefish/go-task-driver"
	"time"
)

func main() {

	service.Service.SetName(`test service`) // set service name
	service.Service.SetPath(`/api/test`)
	logger.LoggerDriverInstance.Register(go_logger.Logger.CloneWithLevel("debug"))
	global_api_strategy.ParamValidateStrategyInstance.SetErrorCode(2005)

	type Params1 struct {
		Test     string `json:"test" validate:"is-mobile"`
		TestNum  uint64 `json:"test_num" validate:"required,lte=100"`
		TestNum1 uint64 `json:"test_num1" validate:"required,gte=1"`
	}

	type Params2 struct {
		TestNum uint64 `json:"test_num" validate:"gte=1" default:"1"`
	}

	service.Service.SetRoutes([]*api.Api{
		{
			Description: "this is a test api",
			Path:        "/v1/test_api/{token_id:[0-9]*}.json",
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
				var params Params1
				apiSession.MustScanParams(&params)
				tokenId := apiSession.PathVars()["token_id"]
				//return nil, go_error.Wrap(errors.New("haha"))
				return map[string]interface{}{
					"params":  params,
					"tokenId": tokenId,
				}, nil
			},
			Params: Params1{},
		},
		{
			Description: "this is a test api",
			Path:        "/v1/test_api/{token_id:[0-9]*}.json",
			Method:      `GET`,
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
				var params Params2
				apiSession.MustScanParams(&params)

				tokenId := apiSession.PathVars()["token_id"]
				//return nil, go_error.Wrap(errors.New("haha"))
				return map[string]interface{}{
					"params":  params,
					"tokenId": tokenId,
				}, nil
			},
			Params: Params2{},
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

	taskDriver := task_driver.NewTaskDriver()
	taskDriver.SetLogger(go_logger.Logger)
	taskDriver.Register(service.Service)

	taskDriver.RunWait(context.Background())
}

// curl --location --request POST 'http://0.0.0.0:8080/api/test/v1/test_api/1234.json' \
// --header 'Content-Type: application/json' \
// --data-raw '{
//     "test": "16265445433",
//     "test_num": 34
// }'
