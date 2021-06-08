package global_api_strategy

import (
	"github.com/golang/mock/gomock"
	"github.com/pefish/go-core/api"
	api_session "github.com/pefish/go-core/api-session"
	_type "github.com/pefish/go-core/api-session/type"
	mock_type "github.com/pefish/go-core/mock/mock-api-session"
	go_error "github.com/pefish/go-error"
	"github.com/pefish/go-test-assert"
	"testing"
)

func TestParamValidateStrategyClass_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	apiSession := mock_type.NewMockIApiSession(ctrl)
	apiSession.EXPECT().Method().Return("GET").AnyTimes()
	apiSession.EXPECT().UrlParams().Return(map[string]string{
		"test": "haha",
	}).AnyTimes()
	apiSession.EXPECT().SetOriginalParams(gomock.Any()).Do(func(originalParams map[string]interface{}) {
		test.Equal(t, "haha", originalParams["test"].(string))
	}).AnyTimes()
	apiSession.EXPECT().SetParams(gomock.Any()).Do(func(originalParams map[string]interface{}) {
		test.Equal(t, "haha", originalParams["test"].(string))
	}).AnyTimes()
	apiSession.EXPECT().Data(gomock.Any()).DoAndReturn(func(p string) interface{} {
		if p == "error_msg" {
			return "this is error"
		}
		return ""
	}).AnyTimes()
	apiSession.EXPECT().SetData(gomock.Any(), gomock.Any()).AnyTimes()
	testApi := &api.Api{
		Description:            "test",
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
	test.Equal(t, (*go_error.ErrorInfo)(nil), err)

	testApi.Params = struct {
		Test string `json:"test" validate:"required,is-mobile"`
	}{}
	err = ParamValidateStrategyInstance.Execute(apiSession, nil)
	test.Equal(t, uint64(1), err.Code)
	test.Equal(t, "test", err.Data.(map[string]interface{})["field"].(string))
	test.Equal(t, `ErrorInfo -> msg: Key: 'Test'; Error:Field validation for 'test' failed on the 'is-mobile' tag; sql-inject-check,required,is-mobile, code: 1, data: map[string]interface {}{"field":"test"}, err: &errors.errorString{s:"Key: 'Test'; Error:Field validation for 'test' failed on the 'is-mobile' tag; sql-inject-check,required,is-mobile"}`, err.String())
}