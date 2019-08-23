package route

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-strategy"
	"test/controllers"
	"time"
)

var TestRoute = map[string]*api_channel_builder.Route{
	`test_api`: {
		Description: "这是测试路由",
		Path:        "/v1/test_api",
		Method:      `POST`,
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.RateLimitApiStrategy,
				Param: api_strategy.RateLimitParam{
					Limit: 1 * time.Second,
				},
			},
		},
		ParamType: api_strategy.ALL_TYPE,
		Controller: controllers.TestController.Test,
	},
}
