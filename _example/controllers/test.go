package controllers

import (
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
)

type TestControllerClass struct {
}

var TestController = TestControllerClass{}

type TestParams struct {
	UserId uint64 `json:"user_id" validate:"required"`
	Token  string `json:"token" validate:"required,min=2"`
}

func (this *TestControllerClass) Test(apiSession *api_session.ApiSessionClass) interface{} {

	go_error.ThrowWithData(`haha`, 2000, map[string]interface{}{
		`haha`: 145,
	})
	return apiSession.Params
}
