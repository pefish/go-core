package service

import (
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/service"
	api_strategy2 "test/src/api-strategy"
	"test/src/controllers"
	"time"
)

var rateLimitDb = map[string]time.Time{}

type TestServiceClass struct {
	service.BaseServiceClass
}

var TestService = TestServiceClass{}

func (this *TestServiceClass) Init(opts ...interface{}) service.InterfaceService {
	this.Name = `测试服务api`
	this.Path = `/api/test`
	params := map[string]interface{}{}
	apiControllers := map[string]api_session.ApiHandlerType{}
	if len(opts) > 0 && opts[0] != nil {
		params = opts[0].(map[string]interface{})
		apiControllers = params[`apiControllers`].(map[string]api_session.ApiHandlerType)
	}
	this.Routes = map[string]*service.Route{
		`test_api`: {
			Description: "这是测试路由",
			Path:        "/v1/test_api",
			Method:      "GET",
			Strategies: []service.StrategyRoute{
				{
					Strategy: &api_strategy.ParamValidateApiStrategy,
				},
				{
					Strategy: &api_strategy.RateLimitApiStrategy,
					Param: api_strategy.RateLimitParam{
						Db: &rateLimitDb,
						Limit: 1 * time.Second,
					},
				},
				{
					Strategy: &api_strategy2.TestStrategyClass{},
				},
			},
			Params:     &controllers.TestParams{},
			Controller: apiControllers[`test_api`],
		},
	}
	return this
}
