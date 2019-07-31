package controllers

import (
	"fmt"
	"github.com/pefish/go-core/api-session"
)

type TestControllerClass struct {
}

var TestController = TestControllerClass{}

type TestParams struct {
	UserId int64 `json:"user_id" validate:"required"`
	Token  string `json:"token" validate:"required,min=2"`
}

func (this *TestControllerClass) Test(apiSession *api_session.ApiSessionClass) interface{} {
	testParams := TestParams{}
	fmt.Println(1, apiSession.Params)
	apiSession.ScanParams(&testParams)
	return testParams
}
