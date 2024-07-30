package global_api_strategy

import (
	"sync"

	i_core "github.com/pefish/go-interface/i-core"
)

type GlobalStrategyData struct {
	Strategy i_core.IApiStrategy
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
