package api_channel_builder

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver"
	_interface "github.com/pefish/go-core/interface"
	"github.com/pefish/go-error"
	"github.com/pefish/go-stack"
)

// 必须是一个输入一个输出，输入必须是iris.Context，输出是任意类型，会成为控制器的输入
type InjectFunc func(route *_interface.Route, out *api_session.ApiSessionClass, param interface{})

type InjectObject struct {
	Func      InjectFunc                   // 前置处理器方法
	Param     interface{}                  // 前置处理器的预设参数
	Route     *_interface.Route            // api路由信息
	This      _interface.InterfaceStrategy // 这个策略本身实例
}

type ApiChannelBuilderClass struct { // 负责构建通道以及管理api通道
	Hero          *hero.Hero
	InjectObjects []InjectObject

	ReturnHookFunc _interface.ReturnHookFuncType
}

func DefaultReturnDataFunc(msg string, internalMsg string, code uint64, data interface{}) *_interface.ApiResult {
	if go_application.Application.Debug {
		return &_interface.ApiResult{
			Msg:         msg,
			InternalMsg: internalMsg,
			Code:        code,
			Data:        data,
		}

	} else {
		return &_interface.ApiResult{
			Msg:         msg,
			InternalMsg: ``,
			Code:        code,
			Data:        data,
		}
	}
}

func NewApiChannelBuilder() *ApiChannelBuilderClass {
	return &ApiChannelBuilderClass{
		InjectObjects: []InjectObject{},
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
			driver.Logger.Error(
				"err: " +
					fmt.Sprint(err) +
					"\n" +
					errMsg +
					"\n" +
					apiContext.Ctx.Values().GetString(`error_msg`) +
					"\n" +
					go_stack.Stack.GetStack(go_stack.Option{Skip: 0, Count: 30}))
			apiResult := DefaultReturnDataFunc(msg, internalMsg, code, data)
			if this.ReturnHookFunc != nil {
				hookApiResult, err := this.ReturnHookFunc(apiContext, apiResult)
				if err != nil {
					apiContext.Ctx.JSON(DefaultReturnDataFunc(err.ErrorMessage, err.InternalErrorMessage, err.ErrorCode, err.Data))
					return
				}
				if hookApiResult == nil {
					return
				}
				apiContext.Ctx.JSON(hookApiResult)
			} else {
				apiContext.Ctx.JSON(apiResult)
			}
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
		apiResult := DefaultReturnDataFunc(``, ``, 0, result)
		if this.ReturnHookFunc != nil {
			hookApiResult, err := this.ReturnHookFunc(apiContext, apiResult)
			if err != nil {
				apiContext.Ctx.JSON(DefaultReturnDataFunc(err.ErrorMessage, err.InternalErrorMessage, err.ErrorCode, err.Data))
				return
			}
			if hookApiResult == nil {
				return
			}
			apiContext.Ctx.JSON(hookApiResult)
		} else {
			apiContext.Ctx.JSON(apiResult)
		}
	})
}
