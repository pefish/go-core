package api_channel_builder

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-error"
	"github.com/pefish/go-stack"
)

type ApiResult struct {
	Msg         string      `json:"msg"`
	InternalMsg string      `json:"internal_msg"`
	Code        uint64      `json:"code"`
	Data        interface{} `json:"data"`
}

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
	ReturnDataFunc ReturnDataFuncType         // 每个接口支持自定义返回格式
	ReturnHookFunc ReturnHookFuncType         // 返回前的处理函数
}

type StrategyRoute struct {
	Strategy InterfaceStrategy
	Param    interface{}
	Disable  bool
}

type InterfaceStrategy interface {
	Execute(route *Route, out *api_session.ApiSessionClass, param interface{})
	GetName() string
	GetDescription() string
	GetErrorCode() uint64
}

// 必须是一个输入一个输出，输入必须是iris.Context，输出是任意类型，会成为控制器的输入
type InjectFunc func(route *Route, out *api_session.ApiSessionClass, param interface{})

type InjectObject struct {
	Func  InjectFunc  // 前置处理器方法
	Param interface{} // 前置处理器的预设参数
	Route *Route      // api路由信息
	This  InterfaceStrategy
}

type ApiChannelBuilderClass struct { // 负责构建通道以及管理api通道
	Hero          *hero.Hero
	InjectObjects []InjectObject

	ReturnHookFunc ReturnHookFuncType
	ReturnDataFunc ReturnDataFuncType
}

type ReturnDataFuncType func(msg string, internalMsg string, code uint64, data interface{}, err interface{}) interface{}
type ReturnHookFuncType func(data interface{}, apiContext *api_session.ApiSessionClass) interface{}

func DefaultReturnDataFunc(msg string, internalMsg string, code uint64, data interface{}, err interface{}) interface{} {
	if go_application.Application.Debug {
		return ApiResult{
			Msg:         msg,
			InternalMsg: internalMsg,
			Code:        code,
			Data:        data,
		}

	} else {
		return ApiResult{
			Msg:         msg,
			InternalMsg: ``,
			Code:        code,
			Data:        data,
		}
	}
}

func NewApiChannelBuilder() *ApiChannelBuilderClass {
	return &ApiChannelBuilderClass{
		InjectObjects:  []InjectObject{},
		ReturnDataFunc: DefaultReturnDataFunc,
	}
}

/**
注入前置处理
*/
func (this *ApiChannelBuilderClass) Inject(key string, injectObject InjectObject) *ApiChannelBuilderClass {
	this.InjectObjects = append(this.InjectObjects, injectObject)
	return this
}

/**
wrap api处理器
*/
func (this *ApiChannelBuilderClass) WrapJson(func_ api_session.ApiHandlerType) func(ctx iris.Context) {
	this.Hero = hero.New()
	this.Hero.Register(func(ctx iris.Context) *api_session.ApiSessionClass {
		// api入口
		apiSession := api_session.NewApiSession() // 新建会话
		apiSession.Ctx = ctx
		return apiSession
	})
	return this.Hero.Handler(func(apiContext *api_session.ApiSessionClass) {
		defer go_error.Recover(func(msg string, internalMsg string, code uint64, data interface{}, err interface{}) {
			apiContext.Ctx.StatusCode(iris.StatusOK)
			errMsg := fmt.Sprintf("msg: %s\ninternal_msg: %s", msg, internalMsg)
			logger.LoggerDriver.Error(
				"err: " +
					fmt.Sprint(err) +
					"\n" +
					errMsg +
					"\n" +
					apiContext.Ctx.Values().GetString(`error_msg`) +
					"\n" +
					go_stack.Stack.GetStack(go_stack.Option{Skip: 0, Count: 30}))

			apiContext.Ctx.JSON(this.ReturnDataFunc(msg, internalMsg, code, data, err))
		})

		if apiContext.Ctx.Method() == `OPTIONS` {
			apiContext.Ctx.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
			apiContext.Ctx.Header("Access-Control-Allow-Origin", apiContext.Ctx.GetHeader("Origin"))
			apiContext.Ctx.Header("Access-Control-Allow-Methods", apiContext.Ctx.Method())
			apiContext.Ctx.Header("Access-Control-Allow-Headers", "*")
			apiContext.Ctx.Header("Access-Control-Allow-Credentials", "true")
			apiContext.Ctx.StatusCode(200)
			apiContext.Ctx.Write([]byte(`ok`))
			return
		}

		for _, injectObject := range this.InjectObjects {
			func() {
				defer go_error.Recover(func(msg string, internalMsg string, code uint64, data interface{}, err interface{}) {
					if code == go_error.INTERNAL_ERROR_CODE {
						code = injectObject.This.GetErrorCode()
					}
					go_error.ThrowErrorWithDataInternalMsg(msg, internalMsg, code, data, err)
				})
				injectObject.Func(injectObject.Route, apiContext, injectObject.Param)
			}()
		}
		result := func_(apiContext)
		if result == nil {
			return
		}
		if this.ReturnHookFunc != nil {
			result = this.ReturnHookFunc(result, apiContext)
			if result == nil {
				return
			}
		}
		apiResult := this.ReturnDataFunc(``, ``, 0, result, nil)
		_, err := apiContext.Ctx.JSON(apiResult)
		if err != nil {
			logger.LoggerDriver.Error(err)
			return
		}
		logger.LoggerDriver.DebugF(`api return: %#v`, apiResult)
	})
}
