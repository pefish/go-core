// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pefish/go-http (interfaces: IHttp)

// Package mock_go_http is a generated GoMock package.
package mock_go_http

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	go_http "github.com/pefish/go-http"
)

// MockIHttp is a mock of IHttp interface.
type MockIHttp struct {
	ctrl     *gomock.Controller
	recorder *MockIHttpMockRecorder
}

// MockIHttpMockRecorder is the mock recorder for MockIHttp.
type MockIHttpMockRecorder struct {
	mock *MockIHttp
}

// NewMockIHttp creates a new mock instance.
func NewMockIHttp(ctrl *gomock.Controller) *MockIHttp {
	mock := &MockIHttp{ctrl: ctrl}
	mock.recorder = &MockIHttpMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHttp) EXPECT() *MockIHttpMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockIHttp) Get(arg0 go_http.RequestParam) (*http.Response, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockIHttpMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIHttp)(nil).Get), arg0)
}

// GetForStruct mocks base method.
func (m *MockIHttp) GetForStruct(arg0 go_http.RequestParam, arg1 interface{}) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForStruct indicates an expected call of GetForStruct.
func (mr *MockIHttpMockRecorder) GetForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForStruct", reflect.TypeOf((*MockIHttp)(nil).GetForStruct), arg0, arg1)
}

// MustGet mocks base method.
func (m *MockIHttp) MustGet(arg0 go_http.RequestParam) (*http.Response, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustGet", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// MustGet indicates an expected call of MustGet.
func (mr *MockIHttpMockRecorder) MustGet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustGet", reflect.TypeOf((*MockIHttp)(nil).MustGet), arg0)
}

// MustGetForStruct mocks base method.
func (m *MockIHttp) MustGetForStruct(arg0 go_http.RequestParam, arg1 interface{}) *http.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustGetForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	return ret0
}

// MustGetForStruct indicates an expected call of MustGetForStruct.
func (mr *MockIHttpMockRecorder) MustGetForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustGetForStruct", reflect.TypeOf((*MockIHttp)(nil).MustGetForStruct), arg0, arg1)
}

// MustPost mocks base method.
func (m *MockIHttp) MustPost(arg0 go_http.RequestParam) (*http.Response, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustPost", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// MustPost indicates an expected call of MustPost.
func (mr *MockIHttpMockRecorder) MustPost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustPost", reflect.TypeOf((*MockIHttp)(nil).MustPost), arg0)
}

// MustPostForStruct mocks base method.
func (m *MockIHttp) MustPostForStruct(arg0 go_http.RequestParam, arg1 interface{}) *http.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustPostForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	return ret0
}

// MustPostForStruct indicates an expected call of MustPostForStruct.
func (mr *MockIHttpMockRecorder) MustPostForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustPostForStruct", reflect.TypeOf((*MockIHttp)(nil).MustPostForStruct), arg0, arg1)
}

// MustPostMultipart mocks base method.
func (m *MockIHttp) MustPostMultipart(arg0 go_http.PostMultipartParam) (*http.Response, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustPostMultipart", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// MustPostMultipart indicates an expected call of MustPostMultipart.
func (mr *MockIHttpMockRecorder) MustPostMultipart(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustPostMultipart", reflect.TypeOf((*MockIHttp)(nil).MustPostMultipart), arg0)
}

// Post mocks base method.
func (m *MockIHttp) Post(arg0 go_http.RequestParam) (*http.Response, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Post indicates an expected call of Post.
func (mr *MockIHttpMockRecorder) Post(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockIHttp)(nil).Post), arg0)
}

// PostForStruct mocks base method.
func (m *MockIHttp) PostForStruct(arg0 go_http.RequestParam, arg1 interface{}) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostForStruct", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostForStruct indicates an expected call of PostForStruct.
func (mr *MockIHttpMockRecorder) PostForStruct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostForStruct", reflect.TypeOf((*MockIHttp)(nil).PostForStruct), arg0, arg1)
}

// PostMultipart mocks base method.
func (m *MockIHttp) PostMultipart(arg0 go_http.PostMultipartParam) (*http.Response, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostMultipart", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostMultipart indicates an expected call of PostMultipart.
func (mr *MockIHttpMockRecorder) PostMultipart(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostMultipart", reflect.TypeOf((*MockIHttp)(nil).PostMultipart), arg0)
}
