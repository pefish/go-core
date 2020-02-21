package global_api_strategy


type StrategyData struct {
	Strategy InterfaceStrategy
	Param    interface{}
	Disable  bool
}

type GlobalApiStrategyDriverClass struct {
	GlobalStrategies []StrategyData
}

var GlobalApiStrategyDriver = GlobalApiStrategyDriverClass{
	GlobalStrategies: []StrategyData{},
}

func (this *GlobalApiStrategyDriverClass) Startup() {

}

func (this *GlobalApiStrategyDriverClass) Register(strategyData StrategyData) bool {
	this.GlobalStrategies = append(this.GlobalStrategies, strategyData)
	return true
}
