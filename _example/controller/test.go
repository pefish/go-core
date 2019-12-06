package controller

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
)

type TestControllerClass struct {
}

var TestController = TestControllerClass{}

type TestParam struct {
	UserId uint64 `json:"user_id" validate:"required"`
	Token  string `json:"token" validate:"required,min=2" desc:"token desc"`
	Haha   uint64 `json:"haha" validate:"omitempty" desc:"haha desc" default:"100"`
	Xixi   string `json:"xixi" validate:"omitempty" desc:"xixi desc" default:"xixi"`
}
type TestReturn struct {
	Test string `json:"test"`
}

func (this *TestControllerClass) PostTest(apiSession *api_session.ApiSessionClass) interface{} {
	var params TestParam
	apiSession.ScanParams(&params)
	//fmt.Println(external_service.DepositAddressService.Test(`1`, `22`))
	//go_error.ThrowWithData(`haha`, 2000, map[string]interface{}{
	//	`haha`: 145,
	//})
	return params
}

type Test1Param struct {
	Haha uint64 `json:"haha" validate:"omitempty" desc:"haha desc" default:"100"`
	Xixi string `json:"xixi" validate:"omitempty" desc:"xixi desc" default:"xixi"`
}
type Test1Return struct {
	Test string `json:"test"`
}

func (this *TestControllerClass) GetTest1(apiSession *api_session.ApiSessionClass) interface{} {
	var params Test1Param
	apiSession.ScanParams(&params)
	//util.DepositAddressService.ValidateAddress(`Eth`, `hfghsfghsh`)
	//go_error.ThrowErrorWithInternalMsg(`haha`, `敏感信息`, 2000, errors.New(`hsgfhsgs`))
	return params
	//apiSession.Ctx.Write([]byte(`xixi`))
	//return nil
}

func (this *TestControllerClass) Test1ReturnHook(apiSession *api_session.ApiSessionClass, apiResult *api_channel_builder.ApiResult) (interface{}, error) {
	//a := data.(Test1Return)
	//a.PostTest = `222`
	apiSession.Ctx.Header(`haha`, `xixi`)
	//apiSession.Ctx.Write([]byte(`hhah`))
	//return nil, errors.New(`haha`)
	//apiResult.InternalMsg = `tywtryt`
	return apiResult, nil
}
