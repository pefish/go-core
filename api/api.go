package api

import (
	"fmt"
	_type2 "github.com/pefish/go-core/api-session/type"
	_type "github.com/pefish/go-core/api-strategy/type"
	global_api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	go_stack "github.com/pefish/go-stack"
	"net/http"

	"github.com/pefish/go-application"
	api_session "github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
)

type Api struct {
	Description            string                // api描述
	Path                   string                // api路径
	IgnoreRootPath         bool                  // api路径是否忽略根路径
	IgnoreGlobalStrategies bool                  // 是否跳过全局策略
	Method                 api_session.ApiMethod // api方法
	Strategies             []_type.StrategyData  // api前置处理策略,不包含全局策略
	Params                 interface{}           // api参数
	Return                 interface{}           // api返回值
	Controller             ApiHandlerType        // api业务处理器
	ParamType              string                // 参数类型。默认 application/json，可选 multipart/form-data，空表示都支持
	ReturnHookFunc         ReturnHookFuncType    // 返回前的处理函数
}

func (api *Api) GetDescription() string {
	return api.Description
}

func (api *Api) GetParamType() string {
	return api.ParamType
}

func (api *Api) GetParams() interface{} {
	return api.Params
}

type ReturnHookFuncType func(apiContext _type2.IApiSession, apiResult *ApiResult) (interface{}, *go_error.ErrorInfo)

type ApiResult struct {
	Msg         string      `json:"msg"`
	InternalMsg string      `json:"internal_msg"`
	Code        uint64      `json:"code"`
	Data        interface{} `json:"data"`
}

type ApiHandlerType func(apiSession _type2.IApiSession) (interface{}, *go_error.ErrorInfo)

func DefaultReturnDataFunc(msg string, code uint64, data interface{}) *ApiResult {
	if go_application.Application.Debug {
		return &ApiResult{
			Msg:        msg,
			InternalMsg: msg,
			Code:        code,
			Data:        data,
		}

	} else {
		return &ApiResult{
			Msg:         msg,
			InternalMsg: ``,
			Code:        code,
			Data:       data,
		}
	}
}

/**
wrap api处理器. 一个path一个，方法内分别处理method
*/
func WrapJson(methodController map[string]*Api) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		apiSession := api_session.NewApiSession() // 新建会话
		apiSession.SetResponseWriter(response)
		apiSession.SetRequest(request)
		apiSession.SetStatusCode(api_session.StatusCode_OK)
		// 应用层直接允许跨域。推荐接口层做跨域处理
		apiSession.SetHeader("Vary", "Origin, Access-Control-request-Method, Access-Control-request-Headers")
		apiSession.SetHeader("Access-Control-Allow-Origin", apiSession.Header("Origin"))
		apiSession.SetHeader("Access-Control-Allow-Methods", apiSession.Method())
		apiSession.SetHeader("Access-Control-Allow-Headers", "*")
		apiSession.SetHeader("Access-Control-Allow-Credentials", "true")
		requestMethod := apiSession.Method()
		if requestMethod == string(api_session.ApiMethod_Option) {
			apiSession.WriteText(`ok`)
			return
		}
		var currentApi *Api
		if methodController[requestMethod] != nil { // 优先使用具体方法注册的控制器
			currentApi = methodController[requestMethod]
			apiSession.SetApi(currentApi)
		} else if methodController[string(api_session.ApiMethod_All)] != nil {
			currentApi = methodController[string(api_session.ApiMethod_All)]
			apiSession.SetApi(currentApi)
		} else {
			logger.LoggerDriver.Logger.DebugF("api found but method not found. request path: %s, request method: %s", apiSession.Path(), apiSession.Method())
			apiSession.WriteText(`Not found`)
			return
		}

		errorHandler := func(errorInfo *go_error.ErrorInfo) {
			errMsg := fmt.Sprint(errorInfo)
			logger.LoggerDriver.Logger.Error(
				"err: " +
					errMsg +
					"\n" +
					apiSession.Data(`error_msg`).(string) +
					"\n" +
					go_stack.Stack.GetStack(go_stack.Option{Skip: 0, Count: 30}))
			apiResult := DefaultReturnDataFunc(errorInfo.Err.Error(), errorInfo.Code, errorInfo.Data)
			if currentApi.ReturnHookFunc != nil {
				hookApiResult, err := currentApi.ReturnHookFunc(apiSession, apiResult)
				if err != nil {
					apiSession.WriteJson(DefaultReturnDataFunc(err.Err.Error(), err.Code, err.Data))
					return
				}
				if hookApiResult == nil {
					return
				}
				apiSession.WriteJson(hookApiResult)
			} else {
				apiSession.WriteJson(apiResult)
			}
		}

		defer go_error.Recover(func(errorInfo *go_error.ErrorInfo) {
			errorHandler(errorInfo)
		})

		if !currentApi.IgnoreGlobalStrategies {
			for _, strategyData := range global_api_strategy.GlobalApiStrategyDriver.GlobalStrategies {
				if strategyData.Disable {
					continue
				}
				err := strategyData.Strategy.Execute(apiSession, strategyData.Param)
				if err != nil {
					errorHandler(err)
					return
				}
			}
		}

		for _, strategyData := range currentApi.Strategies {
			if strategyData.Disable {
				continue
			}
			err := strategyData.Strategy.Execute(apiSession, strategyData.Param)
			if err != nil {
				errorHandler(err)
				return
			}
		}

		defer func() {
			for _, defer_ := range apiSession.Defers() {
				defer_()
			}
		}()

		result, errInfo := currentApi.Controller(apiSession)
		if result == nil && errInfo == nil {
			return
		}
		if errInfo != nil {
			errorHandler(errInfo)
			return
		}
		apiResult := DefaultReturnDataFunc(``, 0, result)
		if currentApi.ReturnHookFunc != nil {
			hookApiResult, err := currentApi.ReturnHookFunc(apiSession, apiResult)
			if err != nil {
				apiSession.WriteJson(DefaultReturnDataFunc(err.Err.Error(), err.Code, err.Data))
				return
			}
			if hookApiResult == nil {
				return
			}
			apiSession.WriteJson(hookApiResult)
		} else {
			apiSession.WriteJson(apiResult)
		}
	}
}
