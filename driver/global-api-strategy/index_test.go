package global_api_strategy

import (
	_type "github.com/pefish/go-core/api-session/type"
	go_error "github.com/pefish/go-error"
)

type TestGlobalStrategy struct {
	Test string
}

func (esd *TestGlobalStrategy) Init(param interface{}) {
	esd.Test = "test"
}

func (esd *TestGlobalStrategy) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	return nil
}
func (esd *TestGlobalStrategy) GetName() string {
	return esd.Test
}
func (esd *TestGlobalStrategy) GetDescription() string {
	return "haha"
}
func (esd *TestGlobalStrategy) GetErrorCode() uint64 {
	return 2002
}

var TestGlobalStrategyInstance = TestGlobalStrategy{}