package global_api_strategy

import "github.com/pefish/go-core/driver/global-api-strategy/type"

type GlobalStrategyData struct {
	Strategy _type.IGlobalStrategy
	Param    interface{}
	Disable  bool
}

type GlobalApiStrategyDriverClass struct {
	GlobalStrategies []GlobalStrategyData
}

var GlobalApiStrategyDriver = GlobalApiStrategyDriverClass{
	GlobalStrategies: []GlobalStrategyData{},
}

func (this *GlobalApiStrategyDriverClass) Startup() {

}

func (this *GlobalApiStrategyDriverClass) Register(strategyData GlobalStrategyData) bool {
	this.GlobalStrategies = append(this.GlobalStrategies, strategyData)
	return true
}
