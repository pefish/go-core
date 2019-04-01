package controllers

import (
	"github.com/pefish/go-core/api-session"
	"test/src/dao"
)

type TestControllerClass struct {

}

var TestController = TestControllerClass{}

type TestParams struct {
	UserId string `json:"user_id" validate:"required"`
}
func (this *TestControllerClass) Test(apiSession *api_session.ApiSessionClass) interface{} {
	testParams := TestParams{}
	apiSession.ScanParams(&testParams)
	return dao.TestDao.GetNameByUid(testParams.UserId)
}
