package go_core

import (
	api_strategy "github.com/pefish/go-core/api-strategy"
	api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/service"
)

func NewService(name string) *service.ServiceClass {
	svc := &service.ServiceClass{}
	svc.SetName(name)
	api_strategy2.GlobalApiStrategyDriver.Register(api_strategy2.StrategyData{
		Strategy: &api_strategy.CorsApiStrategy,
	})
	api_strategy2.GlobalApiStrategyDriver.Register(api_strategy2.StrategyData{
		Strategy: &api_strategy.ServiceBaseInfoApiStrategy,
	})
	api_strategy2.GlobalApiStrategyDriver.Register(api_strategy2.StrategyData{
		Strategy: &api_strategy.ParamValidateStrategy,
	})
	return svc
}

var Service = NewService(`default`)
