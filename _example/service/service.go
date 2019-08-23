package service

import (
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/service"
	"test/service/route"
	"time"
)

var rateLimitDb = map[string]time.Time{}

type TestServiceClass struct {
	service.BaseServiceClass
}

var TestService = TestServiceClass{}

func (this *TestServiceClass) Init(opts ...interface{}) service.InterfaceService {
	this.SetName(`测试服务api`)
	this.SetPath(`/api/test`)
	api_strategy.ParamValidateStrategy.SetErrorCode(2005)
	this.SetRoutes(route.TestRoute)
	return this
}
