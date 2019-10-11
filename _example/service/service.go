package service

import (
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/service"
	"test/service/route"
)

type TestServiceClass struct {
	service.BaseServiceClass
}

var TestService = TestServiceClass{}

func (this *TestServiceClass) Init(opts ...interface{}) service.InterfaceService {
	this.SetName(`测试服务api`)
	this.SetPath(`/api/test`)
	api_strategy.RateLimitApiStrategy.SetErrorCode(2006)
	api_strategy.ParamValidateStrategy.SetErrorCode(2005)
	api_strategy.IpFilterStrategy.SetErrorCode(2007)
	api_strategy.CorsApiStrategy.SetAllowedOrigins([]string{`*`})
	//this.AddGlobalStrategy(&api_strategy2.TestApiStrategy, api_strategy2.ApikeyAuthParam{
	//	AllowedType: `hsfgh`,
	//})
	this.SetRoutes(route.TestRoute)
	return this
}
