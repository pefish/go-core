package go_core

import (
	api_strategy "github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/service"
)

func NewService(name string) *service.ServiceClass {
	svc := &service.ServiceClass{}
	svc.SetName(name)
	svc.AddGlobalStrategy(&api_strategy.CorsApiStrategy, nil)
	svc.AddGlobalStrategy(&api_strategy.ServiceBaseInfoApiStrategy, nil)
	svc.AddGlobalStrategy(&api_strategy.ParamValidateStrategy, nil)
	return svc
}

var Service = NewService(`default`)
