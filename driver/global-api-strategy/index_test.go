package global_api_strategy

import (
	i_core "github.com/pefish/go-interface/i-core"
	t_error "github.com/pefish/go-interface/t-error"
)

type TestGlobalStrategy struct {
	Test string
}

func (esd *TestGlobalStrategy) SetErrorCode(code uint64) i_core.IApiStrategy {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) SetErrorMsg(msg string) i_core.IApiStrategy {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) ErrorMsg() string {
	//TODO implement me
	panic("implement me")
}

func (esd *TestGlobalStrategy) Execute(out i_core.IApiSession) *t_error.ErrorInfo {
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
