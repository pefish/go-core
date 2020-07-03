package api_session

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	mock_http "github.com/pefish/go-core/mock/mock-http"
	"github.com/pefish/go-test-assert"
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
		Test string `json:"test"`
		Haha uint64 `json:"haha"`
		Xixi Test `json:"xixi"`
	}

	apiSession := NewApiSession()
	apiSession.Params = map[string]interface{}{
		"test": "test",
		"haha": 4,
		"xixi": Test{a: "a"},
	}
	var result Result
	apiSession.ScanParams(&result)

	test.Equal(t, "test", result.Test)
	test.Equal(t, uint64(4), result.Haha)
	test.Equal(t, "a", result.Xixi.a)
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
	apiSession.ResponseWriter = httpResponseWriter
	apiSession.SetStatusCode(400)
	err := apiSession.WriteJson(map[string]interface{}{
		"haha": "xixi",
	})
	test.Equal(t, nil, err)
	test.Equal(t, `{"haha":"xixi"}`, result)
	test.Equal(t, 400, statusCode)
	test.Equal(t, string(ContentTypeValue_JSON), resHeaders.Get(string(HeaderName_ContentType)))

	type Test struct {
		A string `json:"a"`
	}
	err = apiSession.WriteJson(&Test{A:"a"})
	test.Equal(t, nil, err)
	test.Equal(t, `{"a":"a"}`, result)
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
	apiSession.ResponseWriter = httpResponseWriter
	apiSession.SetStatusCode(400)
	err := apiSession.WriteText("hgfhdfghd")
	test.Equal(t, nil, err)
	test.Equal(t, `hgfhdfghd`, result)
	test.Equal(t, 400, statusCode)
}

func TestApiSessionClass_GetRemoteAddress(t *testing.T) {
	apiSession := NewApiSession()
	apiSession.Request = &http.Request{
		Method: "GET",
		Header: http.Header{},
		RemoteAddr: "124.56.66.7",
	}
	test.Equal(t, "124.56.66.7", apiSession.GetRemoteAddress())

	apiSession.Request = &http.Request{
		Method: "GET",
		Header: http.Header{
			"X-Forwarded-For": []string{"24.56.11.23"},
		},
		RemoteAddr: "124.56.66.7",
	}
	test.Equal(t, "24.56.11.23", apiSession.GetRemoteAddress())
}

func TestApiSessionClass_GetUrlParams(t *testing.T) {
	apiSession := NewApiSession()
	apiSession.Request = &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL: &url.URL{
			RawQuery:   "test=qq&aa=aa&bb=bb",
		},
	}
	test.Equal(t, "map[aa:aa bb:bb test:qq]", fmt.Sprint(apiSession.GetUrlParams()))
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

	apiSession.Request = req
	result, err := apiSession.GetFormValues()
	test.Equal(t, nil, err)
	test.Equal(t, "[value1 value1_value1]", fmt.Sprint(result["field1"]))
	test.Equal(t, "[value2]", fmt.Sprint(result["field2"]))
}

func TestApiSessionClass_ReadJSON(t *testing.T) {
	apiSession := NewApiSession()

	req := &http.Request{
		Method: "POST",
		Header: http.Header{},
		Body:   ioutil.NopCloser(bytes.NewReader([]byte(`{"test":"test","a":"a"}`))),
	}

	apiSession.Request = req
	type Test struct {
		Test string `json:"test"`
		A string `json:"a"`
	}
	var test_ Test
	err := apiSession.ReadJSON(&test_)
	test.Equal(t, nil, err)
	test.Equal(t, "test", test_.Test)
	test.Equal(t, "a", test_.A)
}