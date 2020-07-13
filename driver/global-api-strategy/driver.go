package global_api_strategy


type GlobalStrategyData struct {
	Strategy IGlobalStrategy
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
