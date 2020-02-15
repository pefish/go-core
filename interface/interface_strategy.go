package _interface

import (
	api_session "github.com/pefish/go-core/api-session"
	go_error "github.com/pefish/go-error"
)

type Route struct {
	Description    string                     // api描述
	Path           string                     // api路径
	IgnoreRootPath bool                       // api路径是否忽略根路径
	Method         string                     // api方法
	Strategies     []StrategyRoute            // api前置处理策略
	Params         interface{}                // api参数
	Return         interface{}                // api返回值
	Redirect       map[string]interface{}     // api重定向
	Debug          bool                       // api是否mock
	Controller     api_session.ApiHandlerType // api业务处理器
	ParamType      string                     // 参数类型。默认 application/json，可选 multipart/form-data，空表示都支持
	ReturnHookFunc ReturnHookFuncType         // 返回前的处理函数
}

type StrategyRoute struct {
	Strategy InterfaceStrategy
	Param    interface{}
	Disable  bool
}

type ReturnHookFuncType func(apiContext *api_session.ApiSessionClass, apiResult *ApiResult) (interface{}, *go_error.ErrorInfo)

type ApiResult struct {
	Msg         string      `json:"msg"`
	InternalMsg string      `json:"internal_msg"`
	Code        uint64      `json:"code"`
	Data        interface{} `json:"data"`
}

type InterfaceStrategy interface {
	Init(param interface{})  // 同步的初始化函数
	InitAsync(param interface{}, onAppTerminated chan interface{})  // 异步的初始化函数，应用推出init才能退出的场景
	Execute(route *Route, out *api_session.ApiSessionClass, param interface{})
	GetName() string
	GetDescription() string
	GetErrorCode() uint64
}