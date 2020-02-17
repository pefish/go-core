package global_api_strategy


type StrategyData struct {
	Strategy InterfaceStrategy
	Param    interface{}
	Disable  bool
}

type GlobalApiStrategyDriverClass struct {
	GlobalStrategies map[string]StrategyData
}

var GlobalApiStrategyDriver = GlobalApiStrategyDriverClass{
	GlobalStrategies: map[string]StrategyData{},
}

func (this *GlobalApiStrategyDriverClass) Startup() {

}

func (this *GlobalApiStrategyDriverClass) Register(strategyData StrategyData) bool {
	this.GlobalStrategies[strategyData.Strategy.GetName()] = strategyData
	return true
}
