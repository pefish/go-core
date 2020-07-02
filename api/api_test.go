package api

import (
	"github.com/golang/mock/gomock"
	api_session "github.com/pefish/go-core/api-session"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	mock_http "github.com/pefish/go-core/mock/mock-http"
	"github.com/pefish/go-test-assert"
	"net/http"
	"testing"
)

func TestWrapJson(t *testing.T) {
	var result string
	ctrl := gomock.NewController(t)
	httpResponseWriter := mock_http.NewMockResponseWriter(ctrl)
	httpResponseWriter.EXPECT().Write(gomock.AssignableToTypeOf([]byte{})).DoAndReturn(func(args []byte) (int, error) {
		result = string(args)
		return len(args), nil
	}).AnyTimes()
	httpResponseWriter.EXPECT().WriteHeader(gomock.AssignableToTypeOf(1)).AnyTimes()
	httpResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()


	handler := WrapJson(map[string]*Api{
		string(api_session.ApiMethod_Get): &Api{
			Description:            "test",
			Path:                   "/",
			IgnoreRootPath:         true,
			IgnoreGlobalStrategies: true,
			Method:                 api_session.ApiMethod_Get,
			Controller: func(apiSession *api_session.ApiSessionClass) interface{} {
				apiSession.WriteText(`this is a get api`)
				return nil
			},
			ParamType: global_api_strategy.ALL_TYPE,
		},
		string(api_session.ApiMethod_Post): &Api{
			Description:            "test1",
			Path:                   "/",
			IgnoreRootPath:         true,
			IgnoreGlobalStrategies: true,
			Method:                 api_session.ApiMethod_Post,
			Controller: func(apiSession *api_session.ApiSessionClass) interface{} {
				apiSession.WriteText(`this is a post api`)
				return nil
			},
			ParamType: global_api_strategy.ALL_TYPE,
		},
	})
	handler(httpResponseWriter, &http.Request{
		Method: "GET",
		Header: http.Header{},
	})
	test.Equal(t, "this is a get api", result)

	handler(httpResponseWriter, &http.Request{
		Method: "POST",
		Header: http.Header{},
	})
	test.Equal(t, "this is a post api", result)
}
