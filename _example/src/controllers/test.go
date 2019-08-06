package controllers

import (
	"github.com/pefish/go-core/api-session"
)

type TestControllerClass struct {
}

var TestController = TestControllerClass{}

type TestParams struct {
	UserId uint64 `json:"user_id" validate:"required"`
	Token  string `json:"token" validate:"required,min=2"`
}

func (this *TestControllerClass) Test(apiSession *api_session.ApiSessionClass) interface{} {
	testParams := TestParams{}
	apiSession.ScanParams(&testParams)
	return testParams
}
