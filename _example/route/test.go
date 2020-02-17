package route

import (
	"github.com/pefish/go-core/api"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/api-strategy"
	api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	"test/controller"
	"time"
)

var TestRoute = []*api.Api{
	{
		Description: "这是测试路由",
		Path:        "/v1/test_api",
		Method:      `POST`,
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.RateLimitApiStrategy,
				Param: api_strategy.RateLimitParam{
					Limit: 1 * time.Second,
				},
			},
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
		ParamType:  api_strategy.ALL_TYPE,
		Controller: controller.TestController.PostTest,
		Params: controller.TestParam{
			UserId: 122,
			Token:  "fghsfghs",
		},
		Return: api.ApiResult{
			Data: controller.TestReturn{
				Test: `hha`,
			},
		},
	},
	{
		Description: "这是测试路由",
		Path:        "/v1/test_api1",
		Method:      `GET`,
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.RateLimitApiStrategy,
				Param: api_strategy.RateLimitParam{
					Limit: 1 * time.Second,
				},
			},
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
		ParamType:  api_strategy.ALL_TYPE,
		Controller: controller.TestController.GetTest1,
		ReturnHookFunc: controller.TestController.Test1ReturnHook,
		Params: controller.Test1Param{
			Haha: 122,
		},
		Return: api.ApiResult{
			Data: []controller.Test1Return{
				{
					Test: `111`,
				},
			},
		},
	},
}
