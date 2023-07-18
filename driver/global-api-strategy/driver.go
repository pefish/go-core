package global_api_strategy

import (
	global_api_strategy "github.com/pefish/go-core-type/global-api-strategy"
	"sync"
)

type GlobalStrategyData struct {
	Strategy global_api_strategy.IGlobalApiStrategy
	Param    interface{}
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
	gasd.Do(func() {
		for _, globalStrategy := range gasd.globalStrategies {
			globalStrategy.Strategy.Init(globalStrategy.Param)
		}
	})
}

func (gasd *GlobalApiStrategyDriver) Register(strategyData GlobalStrategyData) bool {
	gasd.globalStrategies = append(gasd.globalStrategies, strategyData)
	return true
}

func (gasd *GlobalApiStrategyDriver) GlobalStrategies() []GlobalStrategyData {
	return gasd.globalStrategies
}
