package global_api_strategy

import (
	"bytes"
	"github.com/golang/mock/gomock"
	mock_type "github.com/pefish/go-core/mock/mock-api-session"
	go_error "github.com/pefish/go-error"
	"github.com/pefish/go-test-assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServiceBaseInfoStrategyClass_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	apiSession := mock_type.NewMockIApiSession(ctrl)
	apiSession.EXPECT().RemoteAddress().Return("0.0.0.0")
	apiSession.EXPECT().Path().Return("/test")
	apiSession.EXPECT().Method().Return("GET")
	apiSession.EXPECT().Data(gomock.Any()).DoAndReturn(func(p string) interface{} {
		if p == "error_msg" {
			return "this is error"
		}
		return ""
	}).AnyTimes()
	apiSession.EXPECT().SetData(gomock.Any(), gomock.Any()).AnyTimes()
	apiSession.EXPECT().UrlParams().Return(map[string]string{
		"test": "haha",
	}).AnyTimes()
	apiSession.EXPECT().Request().Return(&http.Request{
		Method: "GET",
		Header: http.Header{
			"testHeader": []string{"hhah"},
		},
		RemoteAddr: "124.56.66.7",
		Body: ioutil.NopCloser(bytes.NewReader([]byte("body"))),
	}).AnyTimes()
	apiSession.EXPECT().Header(gomock.Any()).DoAndReturn(func(headName string) string {
		headers := map[string]string{
			"lang": "zh-CN",
		}
		return headers[headName]
	}).AnyTimes()
	apiSession.EXPECT().SetLang(gomock.Any()).Do(func(lang string) {
		test.Equal(t, "zh-CN", lang)
	})
	apiSession.EXPECT().SetClientType(gomock.Any()).Do(func(clientType string) {
		test.Equal(t, "web", clientType)
	})

	ServiceBaseInfoApiStrategyInstance.Init(nil)
	err := ServiceBaseInfoApiStrategyInstance.Execute(apiSession, nil)
	test.Equal(t, (*go_error.ErrorInfo)(nil), err)
}