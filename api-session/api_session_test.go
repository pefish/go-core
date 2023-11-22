package api_session

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	mock_http "github.com/pefish/go-core/mock/mock-http"
	go_test_ "github.com/pefish/go-test"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestApiSessionClass_ScanParams(t *testing.T) {
	type Test struct {
		a string
	}

	type Result struct {
		Test string `json:"go_test_"`
		Haha uint64 `json:"haha"`
		Xixi Test   `json:"xixi"`
	}

	apiSession := NewApiSession()
	apiSession.params = map[string]interface{}{
		"go_test_": "go_test_",
		"haha":     4,
		"xixi":     Test{a: "a"},
	}
	var result Result
	apiSession.ScanParams(&result)

	go_test_.Equal(t, "go_test_", result.Test)
	go_test_.Equal(t, uint64(4), result.Haha)
	go_test_.Equal(t, "a", result.Xixi.a)
}

func TestApiSessionClass_WriteJson(t *testing.T) {
	var result string
	var statusCode int
	var resHeaders = make(http.Header)

	ctrl := gomock.NewController(t)
	httpResponseWriter := mock_http.NewMockResponseWriter(ctrl)
	httpResponseWriter.EXPECT().Write(gomock.AssignableToTypeOf([]byte{})).DoAndReturn(func(args []byte) (int, error) {
		result = string(args)
		return len(args), nil
	}).AnyTimes()
	httpResponseWriter.EXPECT().WriteHeader(gomock.AssignableToTypeOf(1)).AnyTimes().Do(func(code int) {
		statusCode = code
	})
	httpResponseWriter.EXPECT().Header().Return(resHeaders).AnyTimes()

	apiSession := NewApiSession()
	apiSession.responseWriter = httpResponseWriter
	apiSession.SetStatusCode(400)
	err := apiSession.WriteJson(map[string]interface{}{
		"haha": "xixi",
	})
	go_test_.Equal(t, nil, err)
	go_test_.Equal(t, `{"haha":"xixi"}`, result)
	go_test_.Equal(t, 400, statusCode)
	go_test_.Equal(t, string(ContentTypeValue_JSON), resHeaders.Get(string(HeaderName_ContentType)))

	type Test struct {
		A string `json:"a"`
	}
	err = apiSession.WriteJson(&Test{A: "a"})
	go_test_.Equal(t, nil, err)
	go_test_.Equal(t, `{"a":"a"}`, result)
}

func TestApiSessionClass_WriteText(t *testing.T) {
	var result string
	var statusCode int

	ctrl := gomock.NewController(t)
	httpResponseWriter := mock_http.NewMockResponseWriter(ctrl)
	httpResponseWriter.EXPECT().Write(gomock.AssignableToTypeOf([]byte{})).DoAndReturn(func(args []byte) (int, error) {
		result = string(args)
		return len(args), nil
	}).AnyTimes()
	httpResponseWriter.EXPECT().WriteHeader(gomock.AssignableToTypeOf(1)).AnyTimes().Do(func(code int) {
		statusCode = code
	})
	httpResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()

	apiSession := NewApiSession()
	apiSession.responseWriter = httpResponseWriter
	apiSession.SetStatusCode(400)
	err := apiSession.WriteText("hgfhdfghd")
	go_test_.Equal(t, nil, err)
	go_test_.Equal(t, `hgfhdfghd`, result)
	go_test_.Equal(t, 400, statusCode)
}

func TestApiSessionClass_GetRemoteAddress(t *testing.T) {
	apiSession := NewApiSession()
	apiSession.request = &http.Request{
		Method:     "GET",
		Header:     http.Header{},
		RemoteAddr: "124.56.66.7",
	}
	go_test_.Equal(t, "124.56.66.7", apiSession.RemoteAddress())

	apiSession.request = &http.Request{
		Method: "GET",
		Header: http.Header{
			"X-Forwarded-For": []string{"24.56.11.23"},
		},
		RemoteAddr: "124.56.66.7",
	}
	go_test_.Equal(t, "24.56.11.23", apiSession.RemoteAddress())
}

func TestApiSessionClass_GetUrlParams(t *testing.T) {
	apiSession := NewApiSession()
	apiSession.request = &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL: &url.URL{
			RawQuery: "go_test_=qq&aa=aa&bb=bb",
		},
	}
	go_test_.Equal(t, "map[aa:aa bb:bb go_test_:qq]", fmt.Sprint(apiSession.UrlParams()))
}

func TestApiSessionClass_GetFormValues(t *testing.T) {
	apiSession := NewApiSession()

	postData :=
		`--xxx
Content-Disposition: form-data; name="field1"

value1
--xxx
Content-Disposition: form-data; name="field1"

value1_value1
--xxx
Content-Disposition: form-data; name="field2"

value2
--xxx
`
	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
		Body:   ioutil.NopCloser(bytes.NewReader([]byte(postData))),
	}

	apiSession.request = req
	result, err := apiSession.FormValues()
	go_test_.Equal(t, nil, err)
	go_test_.Equal(t, "[value1 value1_value1]", fmt.Sprint(result["field1"]))
	go_test_.Equal(t, "[value2]", fmt.Sprint(result["field2"]))
}

func TestApiSessionClass_ReadJSON(t *testing.T) {
	apiSession := NewApiSession()

	req := &http.Request{
		Method: "POST",
		Header: http.Header{},
		Body:   ioutil.NopCloser(bytes.NewReader([]byte(`{"go_test_":"go_test_","a":"a"}`))),
	}

	apiSession.request = req
	type Test struct {
		Test string `json:"go_test_"`
		A    string `json:"a"`
	}
	var test_ Test
	err := apiSession.ReadJSON(&test_)
	go_test_.Equal(t, nil, err)
	go_test_.Equal(t, "go_test_", test_.Test)
	go_test_.Equal(t, "a", test_.A)
}
