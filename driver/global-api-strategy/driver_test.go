package global_api_strategy

import (
	go_test_ "github.com/pefish/go-test"
	"testing"
)

func TestExternalServiceDriver_Register(t *testing.T) {
	GlobalApiStrategyDriverInstance.Register(GlobalStrategyData{
		Strategy: &TestGlobalStrategyInstance,
		Param:    nil,
		Disable:  false,
	})

	GlobalApiStrategyDriverInstance.Startup()

	results := GlobalApiStrategyDriverInstance.GlobalStrategies()
	go_test_.Equal(t, 1, len(results))

	go_test_.Equal(t, "go_test_", results[0].Strategy.GetName())
}
