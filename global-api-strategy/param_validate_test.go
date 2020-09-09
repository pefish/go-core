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
	apiSession.EXPECT().Method().Return("GET")
	apiSession.EXPECT().UrlParams().Return(map[string]string{
		"test": "haha",
	})
	apiSession.EXPECT().SetOriginalParams(gomock.Any()).Do(func(originalParams map[string]interface{}) {
		test.Equal(t, "haha", originalParams["test"].(string))
	})
	apiSession.EXPECT().SetParams(gomock.Any()).Do(func(originalParams map[string]interface{}) {
		test.Equal(t, "haha", originalParams["test"].(string))
	})
	apiSession.EXPECT().Data(gomock.Any()).DoAndReturn(func(p string) interface{} {
		if p == "error_msg" {
			return "this is error"
		}
		return ""
	})
	apiSession.EXPECT().SetData(gomock.Any(), gomock.Any())
	apiSession.EXPECT().Api().Return(&api.Api{
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
	})

	ParamValidateStrategyInstance.Init(nil)

	err := ParamValidateStrategyInstance.Execute(apiSession, nil)
	test.Equal(t, (*go_error.ErrorInfo)(nil), err)
}