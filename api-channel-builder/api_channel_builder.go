package api_channel_builder

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"github.com/pefish/go-jwt"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-reflect"
	"reflect"
	"runtime/debug"
	"time"
)

// 必须是一个输入一个输出，输入必须是iris.Context，输出是任意类型，会成为控制器的输入
type InjectFuncType func(ctx iris.Context, out *api_session.ApiSessionClass)

type ApiChannelBuilderClass struct { // 负责构建通道以及管理api通道
	Hero        *hero.Hero
	InjectFuncs map[string]InjectFuncType
	Shares      map[string]interface{} // 所有session的共享数据
}

var ApiChannelBuilder = ApiChannelBuilderClass{}

func NewApiChannelBuilder() *ApiChannelBuilderClass {
	return &ApiChannelBuilderClass{
		InjectFuncs: map[string]InjectFuncType{},
		Shares:      map[string]interface{}{},
	}
}

func (this *ApiChannelBuilderClass) JwtAuth(jwtHeaderName string, jwtSecret interface{}) *ApiChannelBuilderClass {
	if this.InjectFuncs == nil {
		this.InjectFuncs = map[string]InjectFuncType{}
	}
	this.InjectFuncs[`jwtAuth`] = func(ctx iris.Context, out *api_session.ApiSessionClass) {
		out.JwtHeaderName = jwtHeaderName
		verifyResult := p_jwt.Jwt.VerifyJwt(p_reflect.Reflect.ToString(jwtSecret), ctx.GetHeader(jwtHeaderName))
		if !verifyResult {
			p_error.ThrowInternal(`jwt verify error`)
		}
		out.JwtPayload = p_jwt.Jwt.DecodePayloadOfJwtBody(ctx.GetHeader(jwtHeaderName))
		if out.JwtPayload[`user_id`] == nil {
			p_error.ThrowInternal(`jwt verify error`)
		}

		//needRightInt := p_format.Format.StringToInt64(needRight)
		//if needRightInt&p_reflect.Reflect.ToInt64(out.JwtPayload[`right`]) != needRightInt {
		//	p_error.Throw(`right not enough`, p_error_codes.ERROR_JWT_RIGHT_NOT_ENOUGH)
		//}
		userId := p_reflect.Reflect.ToUint64(out.JwtPayload[`user_id`])
		out.UserId = &userId
	}
	return this
}

func (this *ApiChannelBuilderClass) RateLimit(duration time.Duration) *ApiChannelBuilderClass {
	if this.InjectFuncs == nil {
		this.InjectFuncs = map[string]InjectFuncType{}
	}
	this.InjectFuncs[`rateLimit`] = func(ctx iris.Context, out *api_session.ApiSessionClass) {
		methodPath := fmt.Sprintf(`%s_%s`, ctx.Method(), ctx.Path())
		key := fmt.Sprintf(`%s_%s`, ctx.RemoteAddr(), methodPath)
		rateLimitRequestsTemp := this.Shares[`rateLimitRequests`]
		if rateLimitRequestsTemp == nil {
			this.Shares[`rateLimitRequests`] = map[string]time.Time{}
		}

		rateLimitRequests := this.Shares[`rateLimitRequests`].(map[string]time.Time)
		if !rateLimitRequests[key].IsZero() && time.Now().Sub(rateLimitRequests[key]) < duration {
			p_error.ThrowInternal(`api ratelimit`)
		}

		rateLimitRequests[key] = time.Now()
		this.Shares[`rateLimitRequests`] = rateLimitRequests
	}
	return this
}

func (this *ApiChannelBuilderClass) RouteBaseInfo(routeName string, route *api_session.Route) *ApiChannelBuilderClass {
	if this.InjectFuncs == nil {
		this.InjectFuncs = map[string]InjectFuncType{}
	}
	this.InjectFuncs[`routeBaseInfo`] = func(ctx iris.Context, out *api_session.ApiSessionClass) {
		out.RouteName = routeName
		out.Route = route
	}
	return this
}

func (this *ApiChannelBuilderClass) ParamValidate(validator *validator.Validate) *ApiChannelBuilderClass {
	if this.InjectFuncs == nil {
		this.InjectFuncs = map[string]InjectFuncType{}
	}
	this.InjectFuncs[`paramValidate`] = func(ctx iris.Context, out *api_session.ApiSessionClass) {
		out.Validator = validator
	}
	return this
}

/**
注入前置处理
*/
func (this *ApiChannelBuilderClass) Inject(key string, func_ func(ctx iris.Context, out *api_session.ApiSessionClass)) *ApiChannelBuilderClass {
	if this.InjectFuncs == nil {
		this.InjectFuncs = map[string]InjectFuncType{}
	}
	this.InjectFuncs[key] = func_
	return this
}

func (this *ApiChannelBuilderClass) register() {
	this.Hero = hero.New()
	if this.InjectFuncs != nil && len(this.InjectFuncs) > 0 {
		this.Hero.Register(func(ctx iris.Context) *api_session.ApiSessionClass {
			// api入口
			apiSession := api_session.NewApiSession() // 新建回话
			apiSession.Ctx = ctx
			for _, fun := range this.InjectFuncs { // 利用闭包实现注入函数的分发
				fun(ctx, apiSession)
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
		map_ := map[string]interface{}{}
		map_[`succeed`] = true
		result := func_(apiContext)
		if result != nil {
			if reflect.TypeOf(result).Kind() == reflect.Slice && reflect.ValueOf(result).Len() == 0 {
				map_[`data`] = []interface{}{}
			} else {
				map_[`data`] = result
			}
			apiContext.Ctx.JSON(map_)
		}
	})
}

func (this *ApiChannelBuilderClass) WrapJsonWithTodo(func_ api_session.ApiHandlerType) func(ctx iris.Context) {
	this.register()
	return this.Hero.Handler(func(apiContext *api_session.ApiSessionClass) {
		map_ := map[string]interface{}{}
		map_[`succeed`] = true
		result := func_(apiContext)
		if result != nil {
			map_[`data`] = result
			map_[`todo`] = true
			apiContext.Ctx.JSON(map_)
		}
	})
}

func (this *ApiChannelBuilderClass) Wrap(func_ api_session.ApiHandlerType) func(ctx iris.Context) {
	this.register()
	return this.Hero.Handler(func(apiContext *api_session.ApiSessionClass) {
		func_(apiContext)
	})
}

type ApiResult struct {
	ErrorMessage string      `json:"error_message"`
	ErrorCode    uint64       `json:"error_code"`
	Data         interface{} `json:"data"`
	Succeed      bool        `json:"succeed"`
}

func (this *ApiChannelBuilderClass) CatchError(ctx iris.Context) {
	if err := recover(); err != nil {
		lang := ctx.GetHeader(`lang`)
		if lang == `` {
			lang = `zh`
		}
		var apiResult ApiResult
		if _, ok := err.(p_error.ErrorInfo); !ok {
			p_logger.Logger.Error(fmt.Sprintf(`ERROR: %v`, err))
			errorMessage := ``
			if _, ok := err.(error); !ok {
				errorMessage = err.(string)
			} else {
				errorMessage = err.(error).Error()
			}
			ctx.StatusCode(iris.StatusOK)
			if p_application.Application.Debug {
				p_logger.Logger.Error(string(debug.Stack()))
				apiResult = ApiResult{
					Succeed:      false,
					ErrorMessage: errorMessage,
					ErrorCode:    1000,
					Data:         nil,
				}
			} else {
				apiResult = ApiResult{
					Succeed:   false,
					ErrorCode: 1000,
					Data:      nil,
				}
			}
			ctx.JSON(apiResult)
		} else {
			ctx.StatusCode(iris.StatusOK)
			errorInfoStruct := err.(p_error.ErrorInfo)
			p_logger.Logger.Error(fmt.Sprintf(`ERROR: %v`, errorInfoStruct.ErrorMessage))
			if p_application.Application.Debug {
				p_logger.Logger.Error(string(debug.Stack()))
				apiResult = ApiResult{
					Succeed:      false,
					ErrorMessage: errorInfoStruct.ErrorMessage,
					ErrorCode:    errorInfoStruct.ErrorCode,
					Data:         errorInfoStruct.Data,
				}
			} else {
				apiResult = ApiResult{
					Succeed:   false,
					ErrorCode: errorInfoStruct.ErrorCode,
					Data:      errorInfoStruct.Data,
				}
			}
			ctx.JSON(apiResult)
		}
	}
}
