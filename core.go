package go_core

import (
	api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/global-api-strategy"
	"github.com/pefish/go-core/service"
)

// New Service instance
func NewService(name string) *service.ServiceClass {
	svc := &service.ServiceClass{}
	svc.SetName(name)
	api_strategy2.GlobalApiStrategyDriver.Register(api_strategy2.GlobalStrategyData{
		Strategy: &global_api_strategy.ServiceBaseInfoApiStrategy,
	})
	api_strategy2.GlobalApiStrategyDriver.Register(api_strategy2.GlobalStrategyData{
		Strategy: &global_api_strategy.ParamValidateStrategy,
	})
	return svc
}

// Default Service instance
var Service = NewService(`default`)
