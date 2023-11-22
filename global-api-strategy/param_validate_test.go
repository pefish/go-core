package global_api_strategy

import (
	"github.com/golang/mock/gomock"
	_type "github.com/pefish/go-core-type/api-session"
	"github.com/pefish/go-core/api"
	api_session "github.com/pefish/go-core/api-session"
	mock_type "github.com/pefish/go-core/mock/mock-api-session"
	go_error "github.com/pefish/go-error"
	go_test_ "github.com/pefish/go-test"
	"testing"
)

func TestParamValidateStrategyClass_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	apiSession := mock_type.NewMockIApiSession(ctrl)
	apiSession.EXPECT().Method().Return("GET").AnyTimes()
	apiSession.EXPECT().UrlParams().Return(map[string]string{
		"go_test_": "haha",
	}).AnyTimes()
	apiSession.EXPECT().SetOriginalParams(gomock.Any()).Do(func(originalParams map[string]interface{}) {
		go_test_.Equal(t, "haha", originalParams["go_test_"].(string))
	}).AnyTimes()
	apiSession.EXPECT().SetParams(gomock.Any()).Do(func(originalParams map[string]interface{}) {
		go_test_.Equal(t, "haha", originalParams["go_test_"].(string))
	}).AnyTimes()
	apiSession.EXPECT().Data(gomock.Any()).DoAndReturn(func(p string) interface{} {
		if p == "error_msg" {
			return "this is error"
		}
		return ""
	}).AnyTimes()
	apiSession.EXPECT().SetData(gomock.Any(), gomock.Any()).AnyTimes()
	testApi := &api.Api{
		Description:            "go_test_",
		Path:                   "/",
		IgnoreRootPath:         true,
		IgnoreGlobalStrategies: true,
		Method:                 api_session.ApiMethod_Get,
		Controller: func(apiSession _type.IApiSession) (i interface{}, info *go_error.ErrorInfo) {
			apiSession.WriteText(`this is a get api`)
			return nil, nil
		},
		ParamType: ALL_TYPE,
	}
	apiSession.EXPECT().Api().Return(testApi).AnyTimes()

	ParamValidateStrategyInstance.Init(nil)

	err := ParamValidateStrategyInstance.Execute(apiSession, nil)
	go_test_.Equal(t, (*go_error.ErrorInfo)(nil), err)

	testApi.Params = struct {
		Test string `json:"go_test_" validate:"required,is-mobile"`
	}{}
	err = ParamValidateStrategyInstance.Execute(apiSession, nil)
	go_test_.Equal(t, uint64(1), err.Code)
	go_test_.Equal(t, "go_test_", err.Data.(map[string]interface{})["field"].(string))
	go_test_.Equal(t, `ErrorInfo -> msg: Key: 'Test'; Error:Field validation for 'go_test_' failed on the 'is-mobile' tag; sql-inject-check,required,is-mobile, code: 1, data: map[string]interface {}{"field":"go_test_"}, err: &errors.errorString{s:"Key: 'Test'; Error:Field validation for 'go_test_' failed on the 'is-mobile' tag; sql-inject-check,required,is-mobile"}`, err.String())
}
