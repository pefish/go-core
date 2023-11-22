package global_api_strategy

import (
	_type "github.com/pefish/go-core-type/api-session"
	global_api_strategy "github.com/pefish/go-core-type/global-api-strategy"
	go_error "github.com/pefish/go-error"
)

type TestGlobalStrategy struct {
	Test string
}

func (esd *TestGlobalStrategy) SetErrorCode(code uint64) global_api_strategy.IGlobalApiStrategy {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) SetErrorMsg(msg string) global_api_strategy.IGlobalApiStrategy {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) GetErrorMsg() string {
	//TODO implement me
	panic("implement me")
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
