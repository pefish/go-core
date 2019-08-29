package api_channel_builder

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"github.com/pefish/go-logger"
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
	InjectObjects map[string]InjectObject
}

var ApiChannelBuilder = ApiChannelBuilderClass{}

func NewApiChannelBuilder() *ApiChannelBuilderClass {
	return &ApiChannelBuilderClass{
		InjectObjects: map[string]InjectObject{},
	}
}

/**
注入前置处理
*/
func (this *ApiChannelBuilderClass) Inject(key string, injectObject InjectObject) *ApiChannelBuilderClass {
	if this.InjectObjects == nil {
		this.InjectObjects = map[string]InjectObject{}
	}
	this.InjectObjects[key] = injectObject
	return this
}

func (this *ApiChannelBuilderClass) register() {
	this.Hero = hero.New()
	if this.InjectObjects != nil && len(this.InjectObjects) > 0 {
		this.Hero.Register(func(ctx iris.Context) *api_session.ApiSessionClass {
			// api入口
			apiSession := api_session.NewApiSession() // 新建会话
			apiSession.Ctx = ctx
			for _, injectObject := range this.InjectObjects { // 利用闭包实现注入函数的分发
				func() {
					defer go_error.Recover(func(msg string, code uint64, data interface{}, err interface{}) {
						if code == go_error.INTERNAL_ERROR_CODE {
							code = injectObject.This.GetErrorCode()
						}
						go_error.Throw(msg, code)
					})
					injectObject.Func(injectObject.Route, apiSession, injectObject.Param)
				}()
			}
			return apiSession
		})
	} else {
		this.Hero.Register(func(ctx iris.Context) *api_session.ApiSessionClass {
			// api入口
			apiSession := api_session.NewApiSession()
			apiSession.Ctx = ctx
			return apiSession
		})
	}
}

/**
wrap api处理器
*/
func (this *ApiChannelBuilderClass) WrapJson(func_ api_session.ApiHandlerType) func(ctx iris.Context) {
	this.register()
	return this.Hero.Handler(func(apiContext *api_session.ApiSessionClass) {
		result := func_(apiContext)
		if result != nil {
			apiResult := ApiResult{
				Msg:  ``,
				Code: 0,
				Data: result,
			}
			_, err := apiContext.Ctx.JSON(apiResult)
			if err != nil {
				go_logger.Logger.Error(err)
				return
			}
			go_logger.Logger.DebugF(`api return: %#v`, apiResult)
		}
	})
}
