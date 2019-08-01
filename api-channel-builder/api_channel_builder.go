package api_channel_builder

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-logger"
)

type ApiResult struct {
	ErrorMessage *string     `json:"msg"`
	ErrorCode    uint64      `json:"code"`
	Data         interface{} `json:"data"`
}

// 必须是一个输入一个输出，输入必须是iris.Context，输出是任意类型，会成为控制器的输入
type InjectFunc func(ctx iris.Context, out *api_session.ApiSessionClass, param interface{})

type InjectObject struct {
	Func  InjectFunc
	Param interface{}
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
			apiSession := api_session.NewApiSession() // 新建回话
			apiSession.Ctx = ctx
			for _, injectObject := range this.InjectObjects { // 利用闭包实现注入函数的分发
				injectObject.Func(ctx, apiSession, injectObject.Param)
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
			_, err := apiContext.Ctx.JSON(ApiResult{
				ErrorMessage: nil,
				ErrorCode:    0,
				Data:         result,
			})
			if err != nil {
				go_logger.Logger.Error(err)
			}
		}
	})
}

func (this *ApiChannelBuilderClass) Wrap(func_ api_session.ApiHandlerType) func(ctx iris.Context) {
	this.register()
	return this.Hero.Handler(func(apiContext *api_session.ApiSessionClass) {
		func_(apiContext)
	})
}
