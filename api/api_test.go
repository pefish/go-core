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
	handler := WrapJson(map[string]*Api{
		string(api_session.ApiMethod_All): &Api{
			Description:            "test",
			Path:                   "/",
			IgnoreRootPath:         true,
			IgnoreGlobalStrategies: true,
			Method:                 api_session.ApiMethod_All,
			Controller: func(apiSession *api_session.ApiSessionClass) interface{} {
				apiSession.WriteText(`hahahahahahahaha`)
				return nil
			},
			ParamType: global_api_strategy.ALL_TYPE,
		},
	})

	var result string
	ctrl := gomock.NewController(t)
	httpResponseWriter := mock_http.NewMockResponseWriter(ctrl)
	httpResponseWriter.EXPECT().Write(gomock.AssignableToTypeOf([]byte{})).DoAndReturn(func(args []byte) (int, error) {
		result = string(args)
		return len(args), nil
	})
	httpResponseWriter.EXPECT().WriteHeader(gomock.AssignableToTypeOf(1))
	httpResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
	handler(httpResponseWriter, &http.Request{
		Method: "POST",
		Header: http.Header{},
	})

	test.Equal(t, "hahahahahahahaha", result)
}
