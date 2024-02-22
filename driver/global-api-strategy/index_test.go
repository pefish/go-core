package global_api_strategy

import (
	_type "github.com/pefish/go-core-type/api-session"
	api_strategy "github.com/pefish/go-core-type/api-strategy"
	go_error "github.com/pefish/go-error"
)

type TestGlobalStrategy struct {
	Test string
}

func (esd *TestGlobalStrategy) SetErrorCode(code uint64) api_strategy.IApiStrategy {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) SetErrorMsg(msg string) api_strategy.IApiStrategy {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) ErrorMsg() string {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) Init(param interface{}) api_strategy.IApiStrategy {
	esd.Test = "test"
	return esd
}

func (esd *TestGlobalStrategy) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	return nil
}
func (esd *TestGlobalStrategy) Name() string {
	return esd.Test
}
func (esd *TestGlobalStrategy) Description() string {
	return "haha"
}
func (esd *TestGlobalStrategy) ErrorCode() uint64 {
	return 2002
}

var TestGlobalStrategyInstance = TestGlobalStrategy{}
