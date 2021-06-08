package global_api_strategy

import (
	"github.com/pefish/go-test-assert"
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
	test.Equal(t, 1, len(results))

	test.Equal(t, "test", results[0].Strategy.GetName())
}