package api_channel_builder

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
	"github.com/pefish/go-stack"
)

type ApiResult struct {
	Msg  string      `json:"msg"`
	Code uint64      `json:"code"`
	Data interface{} `json:"data"`
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
		defer CatchError(apiContext.Ctx)
		apiMsg := fmt.Sprintf(`%s %s %s`, apiContext.Ctx.RemoteAddr(), apiContext.Ctx.Path(), apiContext.Ctx.Method())
		logger.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
		util.UpdateCtxValuesErrorMsg(apiContext.Ctx, `apiMsg`, apiMsg)
		logger.Logger.Debug(apiContext.Ctx.Request().Header)

		if apiContext.Ctx.Method() == `OPTIONS` {
			apiContext.Ctx.StatusCode(200)
			return
		}

		for _, injectObject := range this.InjectObjects {
			func() {
				defer go_error.Recover(func(msg string, code uint64, data interface{}, err interface{}) {
					if code == go_error.INTERNAL_ERROR_CODE {
						code = injectObject.This.GetErrorCode()
					}
					go_error.Throw(msg, code)
				})
				injectObject.Func(injectObject.Route, apiContext, injectObject.Param)
			}()
		}
		result := func_(apiContext)
		if result != nil {
			apiResult := ApiResult{
				Msg:  ``,
				Code: 0,
				Data: result,
			}
			_, err := apiContext.Ctx.JSON(apiResult)
			if err != nil {
				logger.Logger.Error(err)
				return
			}
			logger.Logger.DebugF(`api return: %#v`, apiResult)
		}
	})
}


func CatchError(ctx iris.Context) {
	if err := recover(); err != nil {
		var apiResult ApiResult
		if _, ok := err.(go_error.ErrorInfo); !ok {
			errorMessage := ``
			if _, ok := err.(error); !ok {
				errorMessage = err.(string)
			} else {
				errorMessage = err.(error).Error()
			}
			logger.Logger.Error(`system_error: ` + errorMessage+"\n"+ctx.Values().Get(`error_msg`).(string)+"\n"+go_stack.Stack.GetStack(go_stack.Option{Skip: 0, Count: 7}))
			ctx.StatusCode(iris.StatusOK)
			if go_application.Application.Debug {
				apiResult = ApiResult{
					Msg:  errorMessage,
					Code: 1,
					Data: nil,
				}
			} else {
				apiResult = ApiResult{
					Msg:  ``,
					Code: 1,
					Data: nil,
				}
			}
			ctx.JSON(apiResult)
		} else {
			ctx.StatusCode(iris.StatusOK)
			errorInfoStruct := err.(go_error.ErrorInfo)
			errMsg := `error: ` + errorInfoStruct.ErrorMessage
			if errorInfoStruct.Err != nil {
				errMsg += "\nsystem_error: " + errorInfoStruct.Err.Error()
			}
			logger.Logger.Error(errMsg +"\n"+ctx.Values().GetString(`error_msg`)+"\n"+go_stack.Stack.GetStack(go_stack.Option{Skip: 0, Count: 7}))
			if go_application.Application.Debug {
				apiResult = ApiResult{
					Msg:  errorInfoStruct.ErrorMessage,
					Code: errorInfoStruct.ErrorCode,
					Data: errorInfoStruct.Data,
				}
			} else {
				apiResult = ApiResult{
					Msg:  ``,
					Code: errorInfoStruct.ErrorCode,
					Data: errorInfoStruct.Data,
				}
			}
			ctx.JSON(apiResult)
		}
	}
}

