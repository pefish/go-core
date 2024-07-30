package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	driver_global_api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"

	api_session "github.com/pefish/go-core/api-session"
	i_core "github.com/pefish/go-interface/i-core"
	t_core "github.com/pefish/go-interface/t-core"
	t_error "github.com/pefish/go-interface/t-error"
)

type StrategyData struct {
	Strategy i_core.IApiStrategy
	Disable  bool
}

type Api struct {
	description              string                // api描述
	path                     string                // api路径
	isIgnoreRootPath         bool                  // api路径是否忽略根路径
	isIgnoreGlobalStrategies bool                  // 是否跳过全局策略
	method                   api_session.ApiMethod // api方法
	strategies               []StrategyData        // api前置处理策略,不包含全局策略
	params                   interface{}           // api参数
	returnValue              interface{}           // api返回值
	controllerFunc           ApiHandlerType        // api业务处理器
	paramType                string                // 参数类型。默认 application/json，可选 multipart/form-data，空表示都支持
	returnHookFunc           ReturnHookFuncType    // 返回前的处理函数
}

type NewApiParamsType struct {
	Description              string
	Path                     string                // api路径
	IsIgnoreRootPath         bool                  // api路径是否忽略根路径
	IsIgnoreGlobalStrategies bool                  // 是否跳过全局策略
	Method                   api_session.ApiMethod // api方法
	Strategies               []StrategyData        // api前置处理策略,不包含全局策略
	Params                   interface{}           // api参数
	ReturnValue              interface{}           // api返回值
	ControllerFunc           ApiHandlerType        // api业务处理器
	ParamType                string                // 参数类型。默认 application/json，可选 multipart/form-data，空表示都支持
	ReturnHookFunc           ReturnHookFuncType    // 返回前的处理函数
}

func NewApi(params *NewApiParamsType) *Api {
	return &Api{
		description:              params.Description,
		path:                     params.Path,
		isIgnoreRootPath:         params.IsIgnoreRootPath,
		isIgnoreGlobalStrategies: params.IsIgnoreGlobalStrategies,
		method:                   params.Method,
		strategies:               params.Strategies,
		params:                   params.Params,
		returnValue:              params.ReturnValue,
		controllerFunc:           params.ControllerFunc,
		paramType:                params.ParamType,
		returnHookFunc:           params.ReturnHookFunc,
	}
}

func New404Api() *Api {
	return &Api{
		description:              "404 not found",
		path:                     "/",
		isIgnoreRootPath:         true,
		isIgnoreGlobalStrategies: true,
		method:                   api_session.ApiMethod_All,
		controllerFunc: func(apiSession i_core.IApiSession) (interface{}, *t_error.ErrorInfo) {
			global_api_strategy.ServiceBaseInfoStrategyInstance.Execute(apiSession)

			apiSession.SetStatusCode(api_session.StatusCode_NotFound)
			logger.LoggerDriverInstance.Logger.DebugF("api not found. request path: %s, request method: %s", apiSession.Path(), apiSession.Method())
			apiSession.WriteText(`Not Found`)
			return nil, nil
		},
		paramType: global_api_strategy.ALL_TYPE,
	}
}

func (api *Api) Strategies() []StrategyData {
	return api.strategies
}

func (api *Api) Path() string {
	return api.path
}

func (api *Api) ControllerFunc() ApiHandlerType {
	return api.controllerFunc
}

func (api *Api) IsIgnoreRootPath() bool {
	return api.isIgnoreRootPath
}

func (api *Api) Description() string {
	return api.description
}

func (api *Api) Method() api_session.ApiMethod {
	return api.method
}

func (api *Api) ParamType() string {
	return api.paramType
}

func (api *Api) Params() interface{} {
	return api.params
}

type ReturnHookFuncType func(apiSession i_core.IApiSession, apiResult *t_core.ApiResult) (interface{}, *t_error.ErrorInfo)

type ApiHandlerType func(apiSession i_core.IApiSession) (interface{}, *t_error.ErrorInfo)

/*
*
wrap api处理器. 一个path一个，方法内分别处理method
*/
func WrapJson(methodController map[string]*Api) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		apiSession := api_session.NewApiSession() // 新建会话
		apiSession.SetPathVars(mux.Vars(request))
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
			logger.LoggerDriverInstance.Logger.DebugF("api found but method not found. request path: %s, request method: %s", apiSession.Path(), apiSession.Method())
			apiSession.WriteText(`Not found`)
			return
		}

		errorHandler := func(errorInfo *t_error.ErrorInfo) {
			apiResult := &t_core.ApiResult{
				Msg:  errorInfo.Err.Error(),
				Code: errorInfo.Code,
				Data: errorInfo.Data,
			}
			if currentApi.returnHookFunc != nil {
				hookApiResult, errorInfo := currentApi.returnHookFunc(apiSession, apiResult)
				if errorInfo != nil {
					apiSession.WriteJson(&t_core.ApiResult{
						Msg:  errorInfo.Err.Error(),
						Code: errorInfo.Code,
						Data: errorInfo.Data,
					})
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

		defer t_error.Recover(func(errInfo *t_error.ErrorInfo) {
			errMsg := fmt.Sprint(errInfo)
			logger.LoggerDriverInstance.Logger.Error(
				apiSession.Data(`error_msg`).(string) +
					"\n" +
					"err: " +
					errMsg)
			errorHandler(t_error.INTERNAL_ERROR)
		})

		if !currentApi.isIgnoreGlobalStrategies {
			for _, strategyData := range driver_global_api_strategy.GlobalApiStrategyDriverInstance.GlobalStrategies() {
				logger.LoggerDriverInstance.Logger.DebugF("global strategy [%s]: %#v", strategyData.Strategy.Name(), strategyData)
				if strategyData.Disable {
					continue
				}
				errInfo := strategyData.Strategy.Execute(apiSession)
				if errInfo != nil {
					errMsg := fmt.Sprint(errInfo)
					logger.LoggerDriverInstance.Logger.Error(
						apiSession.Data(`error_msg`).(string) +
							"\n" +
							"err: " +
							errMsg)
					errorHandler(errInfo)
					return
				}
			}
		}

		for _, strategyData := range currentApi.strategies {
			logger.LoggerDriverInstance.Logger.DebugF("strategy [%s]: %#v", strategyData.Strategy.Name(), strategyData)
			if strategyData.Disable {
				continue
			}
			errInfo := strategyData.Strategy.Execute(apiSession)
			if errInfo != nil {
				errMsg := fmt.Sprint(errInfo)
				logger.LoggerDriverInstance.Logger.Error(
					apiSession.Data(`error_msg`).(string) +
						"\n" +
						"err: " +
						errMsg)
				errorHandler(errInfo)
				return
			}
		}

		defer func() {
			for _, defer_ := range apiSession.Defers() {
				defer_()
			}
		}()

		result, errInfo := currentApi.controllerFunc(apiSession)
		if result == nil && errInfo == nil {
			return
		}
		if errInfo != nil {
			errMsg := fmt.Sprint(errInfo)
			logger.LoggerDriverInstance.Logger.Error(
				apiSession.Data(`error_msg`).(string) +
					"\n" +
					"err: " +
					errMsg)
			errorHandler(errInfo)
			return
		}
		apiResult := &t_core.ApiResult{
			Msg:  ``,
			Code: 0,
			Data: result,
		}
		if currentApi.returnHookFunc != nil {
			hookApiResult, errorInfo := currentApi.returnHookFunc(apiSession, apiResult)
			if errorInfo != nil {
				apiSession.WriteJson(&t_core.ApiResult{
					Msg:  errorInfo.Err.Error(),
					Code: errorInfo.Code,
					Data: errorInfo.Data,
				})
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
