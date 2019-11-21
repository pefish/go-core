package controller

import (
	"errors"
	"fmt"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"test/external-service"
)

type TestControllerClass struct {
}

var TestController = TestControllerClass{}

type TestParam struct {
	UserId uint64 `json:"user_id" validate:"required"`
	Token  string `json:"token" validate:"required,min=2" desc:"token desc"`
}
type TestReturn struct {
	Test string `json:"test"`
}
func (this *TestControllerClass) Test(apiSession *api_session.ApiSessionClass) interface{} {
	fmt.Println(external_service.DepositAddressService.Test(`1`, `22`))
	go_error.ThrowWithData(`haha`, 2000, map[string]interface{}{
		`haha`: 145,
	})
	return TestReturn{
		Test: `111`,
	}
}


type Test1Param struct {
	Haha uint64 `json:"haha" validate:"omitempty" desc:"haha desc"`
}
type Test1Return struct {
	Test string `json:"test"`
}
func (this *TestControllerClass) Test1(apiSession *api_session.ApiSessionClass) interface{} {
	//util.DepositAddressService.ValidateAddress(`Eth`, `hfghsfghsh`)
	//go_error.ThrowErrorWithInternalMsg(`haha`, `敏感信息`, 2000, errors.New(`hsgfhsgs`))
	return Test1Return{
		Test: `111`,
	}
	//apiSession.Ctx.Write([]byte(`xixi`))
	//return nil
}

func (this *TestControllerClass) Test1ReturnHook(apiSession *api_session.ApiSessionClass, apiResult *api_channel_builder.ApiResult) (interface{}, error) {
	//a := data.(Test1Return)
	//a.Test = `222`
	apiSession.Ctx.Header(`haha`, `xixi`)
	//apiSession.Ctx.Write([]byte(`hhah`))
	return nil, errors.New(`haha`)
	//apiResult.InternalMsg = `tywtryt`
	//return apiResult, nil
}
