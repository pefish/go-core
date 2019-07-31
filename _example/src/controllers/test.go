package controllers

import (
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"test/src/dao"
)

type TestControllerClass struct {
}

var TestController = TestControllerClass{}

type TestParams struct {
	UserId string `json:"user_id" validate:"required"`
	Token  string `json:"token"`
}

func (this *TestControllerClass) Test(apiSession *api_session.ApiSessionClass) interface{} {
	testParams := TestParams{}
	apiSession.ScanParams(&testParams)

	p_error.ThrowInternal(`haha`)
	return dao.TestDao.GetNameByUid(testParams.UserId)
}
