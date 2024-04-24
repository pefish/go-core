package global_api_strategy

import (
	"sync"

	api_strategy "github.com/pefish/go-core-type/api-strategy"
)

type GlobalStrategyData struct {
	Strategy api_strategy.IApiStrategy
	Disable  bool
}

type GlobalApiStrategyDriver struct {
	globalStrategies []GlobalStrategyData
	sync.Once
}

var GlobalApiStrategyDriverInstance = GlobalApiStrategyDriver{
	globalStrategies: []GlobalStrategyData{},
}

func (gasd *GlobalApiStrategyDriver) Startup() {

}

func (gasd *GlobalApiStrategyDriver) Register(strategyData GlobalStrategyData) bool {
	gasd.globalStrategies = append(gasd.globalStrategies, strategyData)
	return true
}

func (gasd *GlobalApiStrategyDriver) GlobalStrategies() []GlobalStrategyData {
	return gasd.globalStrategies
}
