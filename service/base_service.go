package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mitchellh/mapstructure"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/middleware"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-reflect"
	"reflect"
)

type StrategyRoute struct {
	Strategy api_strategy.InterfaceStrategy
	Param    interface{}
	Disable  bool
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
	ParamType      string                     // 参数类型。默认 application/json
}

type BaseServiceClass struct {
	name              string                                      // 服务名
	description       string                                      // 服务描述
	path              string                                      // 服务的基础路径
	host              string                                      // 服务监听host
	port              uint64                                      // 服务监听port
	accessHost        string                                      // 服务访问host，没有设置的话使用监听host
	accessPort        uint64                                      // 服务访问port，没有设置的话使用监听port
	routes            map[string]*Route                           // 服务的所有路由
	Middlewires       map[string]api_channel_builder.InjectObject // 每个api的前置处理器（框架的）
	GlobalMiddlewires map[string]context.Handler                  // 每个api的前置处理器（iris的）
	App               *iris.Application                           // iris实例
	Opts              map[string]interface{}                      // 一些可选参数
}

func (this *BaseServiceClass) SetRoutes(routes map[string]*Route) {
	this.routes = routes
}

func (this *BaseServiceClass) SetPath(path string) {
	this.path = path
}

func (this *BaseServiceClass) SetName(name string) {
	this.name = name
}

func (this *BaseServiceClass) GetHost() string {
	return this.host
}

func (this *BaseServiceClass) SetHost(host string) {
	this.host = host
}

func (this *BaseServiceClass) GetPort() uint64 {
	return this.port
}

func (this *BaseServiceClass) SetPort(port uint64) {
	this.port = port
}

func (this *BaseServiceClass) GetAccessHost() string {
	return this.accessHost
}

func (this *BaseServiceClass) SetAccessHost(accessHost string) {
	this.accessHost = accessHost
}

func (this *BaseServiceClass) GetAccessPort() uint64 {
	return this.accessPort
}

func (this *BaseServiceClass) SetAccessPort(accessPort uint64) {
	this.accessPort = accessPort
}

func (this *BaseServiceClass) SetDescription(desc string) {
	this.description = desc
}

func (this *BaseServiceClass) Init(opts ...interface{}) InterfaceService {
	return this
}

func (this *BaseServiceClass) SetHealthyCheck(func_ func()) InterfaceService {
	this.routes[`healthy_check`] = &Route{
		Description: "健康检查api",
		Path:        "/healthz",
		Method:      "ALL",
		IgnoreRootPath: true,
		Controller: func(apiContext *api_session.ApiSessionClass) interface{} {
			defer func() {
				if err := recover(); err != nil {
					go_logger.Logger.Error(err)
					apiContext.Ctx.StatusCode(iris.StatusInternalServerError)
					apiContext.Ctx.Text(`not ok`)
				}
			}()
			if func_ != nil {
				func_()
			}

			apiContext.Ctx.StatusCode(iris.StatusOK)
			go_logger.Logger.Debug(`I am healthy`)
			apiContext.Ctx.Text(`ok`)
			return nil
		},
	}
	return this
}

func (this *BaseServiceClass) Use(key string, injectObject api_channel_builder.InjectObject) InterfaceService {
	if this.Middlewires == nil {
		this.Middlewires = map[string]api_channel_builder.InjectObject{}
	}
	this.Middlewires[key] = injectObject
	return this
}

func (this *BaseServiceClass) UseGlobal(key string, func_ context.Handler) InterfaceService {
	if this.GlobalMiddlewires == nil {
		this.GlobalMiddlewires = map[string]context.Handler{}
	}
	this.GlobalMiddlewires[key] = func_
	return this
}

func (this *BaseServiceClass) GetName() string {
	return this.name
}

func (this *BaseServiceClass) GetDescription() string {
	return this.description
}

func (this *BaseServiceClass) GetPath() string {
	return this.path
}

func (this *BaseServiceClass) GetRoutes() map[string]*Route {
	return this.routes
}

func (this *BaseServiceClass) ExactOpt(name string) interface{} {
	if name == `jwt_header_name` && (this.Opts == nil || this.Opts[name] == nil) {
		return `Json-Web-Token`
	} else {
		return this.Opts[name]
	}
}

func (this *BaseServiceClass) RequestWithErr(apiName string, args ...interface{}) (interface{}, error) {
	body := this.RequestRawMap(apiName, args...)
	code := body[`code`].(uint64)
	if code != 0 {
		errorMessage := go_error.INTERNAL_ERROR
		if body[`msg`] != nil {
			errorMessage = body[`msg`].(string)
		}
		return body, errors.New(errorMessage)
	}
	return body[`data`], nil
}

func (this *BaseServiceClass) RequestRawMap(apiName string, args ...interface{}) map[string]interface{} {
	var params interface{}
	if len(args) > 0 && args[0] != nil {
		params = args[0]
	} else {
		params = map[string]interface{}{}
	}

	headers := map[string]interface{}{}
	// header内容转发
	if len(args) > 1 && args[1] != nil {
		if apiSession, ok := args[1].(*api_session.ApiSessionClass); ok {
			jwtHeaderName := go_reflect.Reflect.ToString(this.ExactOpt(`jwt_header_name`))
			headers = map[string]interface{}{
				`lang`:            apiSession.Lang,
				`client_type`:     apiSession.ClientType,
				jwtHeaderName:     apiSession.Ctx.GetHeader(jwtHeaderName),
				`X-Forwarded-For`: apiSession.Ctx.GetHeader(`X-Forwarded-For`),
			}
		}
	}
	if this.routes[apiName] == nil {
		go_error.Throw(`api request 404`, go_error.INTERNAL_ERROR_CODE)
	}
	method := this.routes[apiName].Method
	fullUrl := this.GetRequestUrl(apiName)
	body := map[string]interface{}{}
	if method == `GET` {
		body = go_http.Http.GetWithParamsForMap(go_http.RequestParam{
			Url:     fullUrl,
			Params:  params,
			Headers: headers,
		})
	} else if method == `POST` {
		body = go_http.Http.PostForMap(go_http.RequestParam{
			Url:     fullUrl,
			Params:  params,
			Headers: headers,
		})
	} else {
		go_error.Throw(`request not support method`, go_error.INTERNAL_ERROR_CODE)
	}
	return body
}

/**
http请求其他服务。
apiName：请求哪个api
args：args[0]是参数，可以是struct或者map；args[1]是ApiSessionClass，如果存在则会转发一些预制header；
*/
func (this *BaseServiceClass) Request(apiName string, args ...interface{}) (data interface{}) {
	body := this.RequestRawMap(apiName, args...)
	code := body[`code`].(uint64)
	if code != 0 {
		errorMessage := go_error.INTERNAL_ERROR
		if body[`msg`] != nil {
			errorMessage = body[`msg`].(string)
		}
		go_error.ThrowErrorWithData(errorMessage, code, body[`data`], nil)
	}
	data = body[`data`]
	return
}

func (this *BaseServiceClass) RequestForSlice(apiName string, args ...interface{}) []map[string]interface{} {
	requestResult := this.Request(apiName, args...)
	if requestResult == nil {
		return []map[string]interface{}{}
	}
	out := []map[string]interface{}{}
	for _, val := range requestResult.([]interface{}) {
		out = append(out, val.(map[string]interface{}))
	}
	return out
}

func (this *BaseServiceClass) RequestForSliceWithScan(dest interface{}, apiName string, args ...interface{}) {
	requestResult := this.Request(apiName, args...)
	if requestResult == nil {
		dest = nil
	}
	inrec, err := json.Marshal(requestResult)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(inrec, dest)
}

func (this *BaseServiceClass) RequestForMap(apiName string, args ...interface{}) map[string]interface{} {
	requestResult := this.Request(apiName, args...)
	if requestResult == nil {
		return map[string]interface{}{}
	}
	return requestResult.(map[string]interface{})
}

func (this *BaseServiceClass) RequestForMapWithScan(dest interface{}, apiName string, args ...interface{}) {
	requestResult := this.Request(apiName, args...)
	if requestResult == nil {
		dest = nil
	}
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &dest,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(requestResult.(map[string]interface{}))
	if err != nil {
		panic(err)
	}
}

func (this *BaseServiceClass) Run() {
	this.buildRoutes()
	irisConfig := iris.Configuration{}
	irisConfig.RemoteAddrHeaders = map[string]bool{
		`X-Forwarded-For`: true,
	}
	this.printRoutes()
	host := this.host
	if host == `` {
		host = `0.0.0.0`
	}
	err := this.App.Run(iris.Addr(host+`:`+go_reflect.Reflect.ToString(this.port)), iris.WithConfiguration(irisConfig))
	if err != nil {
		panic(err)
	}
}

func (this *BaseServiceClass) printRoutes() {
	for _, route := range this.routes {
		apiPath := this.path+route.Path
		if route.IgnoreRootPath == true {
			apiPath = route.Path
		}
		go_logger.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, route.Method, apiPath, route.Description))
	}
}

func (this *BaseServiceClass) buildRoutes() {
	this.App = iris.New()
	if go_application.Application.Debug {
		this.App.UseGlobal(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{`*`},
			AllowedMethods:   []string{`PUT`, `POST`, `GET`, `DELETE`, `OPTIONS`},
			Debug:            go_application.Application.Debug,
		}))
	}
	this.App.UseGlobal(middleware.ErrorHandle)
	this.App.UseGlobal(middleware.OptionHandle)
	for _, fun := range this.GlobalMiddlewires {
		this.App.UseGlobal(fun)
	}

	for name, route := range this.GetRoutes() {
		var apiChannelBuilder = api_channel_builder.NewApiChannelBuilder()
		// 预定义前置处理器
		apiChannelBuilder.Inject(api_strategy.ServiceBaseInfoApiStrategy.GetName(), api_channel_builder.InjectObject{
			Func: api_strategy.ServiceBaseInfoApiStrategy.Execute,
			Param: api_strategy.ServiceBaseInfoParam{
				RouteName: name,
			},
		})
		apiChannelBuilder.Inject(api_strategy.ParamValidateStrategy.GetName(), api_channel_builder.InjectObject{
			Func: api_strategy.ParamValidateStrategy.Execute,
			Param: api_strategy.ParamValidateParam{
				Param: route.Params,
			},
		})
		for key, injectObject := range this.Middlewires {
			apiChannelBuilder.Inject(key, injectObject)
		}
		if route.Strategies != nil {
			for _, strategyRoute := range route.Strategies {
				if !strategyRoute.Disable {
					apiChannelBuilder.Inject(strategyRoute.Strategy.GetName(), api_channel_builder.InjectObject{
						Func:  strategyRoute.Strategy.Execute,
						Param: strategyRoute.Param,
					})
				}
			}
		}
		apiPath := this.path+route.Path
		if route.IgnoreRootPath == true {
			apiPath = route.Path
		}
		if route.Method == `` {
			route.Method = `ALL`
		}
		if route.Controller == nil {
			if route.Redirect != nil { // 自动转发。不会校验参数
				redirectMap := route.Redirect
				return_ := this.parseReturn(route.Return)
				this.App.AllowMethods(iris.MethodOptions).Handle(route.Method, apiPath, apiChannelBuilder.WrapJson(func(apiContext *api_session.ApiSessionClass) interface{} {
					params := apiContext.Params
					service := redirectMap[`service`].(InterfaceService)
					routeName := redirectMap[`route_name`].(string)
					if service.GetRoutes()[routeName] == nil && return_ != nil { // 目标服务路由不存在，则返回规定的返回值(自动mock)
						return return_
					}
					return service.Request(routeName, params, apiContext) // 自动定位目标api的method
				}))
			} else { // 自动mock
				return_ := this.parseReturn(route.Return)
				if return_ == nil {
					go_error.ThrowInternal(`route config error; route name: ` + name)
				}
				this.App.AllowMethods(iris.MethodOptions).Handle(route.Method, apiPath, apiChannelBuilder.WrapJson(func(apiContext *api_session.ApiSessionClass) interface{} {
					return return_
				}))
			}
		} else {
			this.App.AllowMethods(iris.MethodOptions).Handle(route.Method, apiPath, apiChannelBuilder.WrapJson(route.Controller))
		}
	}

	this.App.AllowMethods(iris.MethodOptions).Handle(``, `/*`, func(ctx context.Context) {
		ctx.StatusCode(iris.StatusNotFound)
		go_logger.Logger.Debug(`api not found`)
		ctx.Text(`Not Found`)
	})
}

func (this *BaseServiceClass) recurStruct(type_ reflect.Type, result map[string]interface{}) {
	for i := 0; i < type_.NumField(); i++ {
		field := type_.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Struct {
			this.recurStruct(fieldType, result)
		} else {
			tagName := field.Tag.Get(`example`)
			if tagName != `` {
				result[field.Tag.Get(`json`)] = tagName
			} else {
				result[field.Tag.Get(`json`)] = nil
			}
		}
	}
}

func (this *BaseServiceClass) parseReturn(return_ interface{}) interface{} {
	var result interface{}
	if return_ == nil {
		return nil
	}
	type_ := reflect.TypeOf(return_)
	kind := type_.Kind()
	if kind == reflect.Map {
		map_ := return_.(map[string]map[string]interface{})
		resultTemp := map[string]interface{}{}
		for key, obj := range map_ {
			resultTemp[key] = obj[`example`]
		}
		result = resultTemp
	} else if kind == reflect.Struct {
		resultTemp := map[string]interface{}{}
		this.recurStruct(type_, resultTemp)
		result = resultTemp
	} else if kind == reflect.Slice {
		resultTemp := []map[string]interface{}{}
		eleType := type_.Elem()
		if eleType.Kind() == reflect.Struct {
			tempMap := map[string]interface{}{}
			this.recurStruct(eleType, tempMap)
			resultTemp = append(resultTemp, tempMap)
		} else if type_.Elem().Kind() == reflect.Map {
			slice_ := return_.([]map[string]map[string]interface{})
			for _, map_ := range slice_ {
				tempMap := map[string]interface{}{}
				for key, obj := range map_ {
					tempMap[key] = obj[`example`]
				}
				resultTemp = append(resultTemp, tempMap)
			}
		} else {
			go_error.ThrowInternal(`route return error`)
		}
		result = resultTemp
	} else {
		go_error.ThrowInternal(`route return error`)
	}
	return result
}

func (this *BaseServiceClass) GetRequestUrl(apiName string) string {
	if this.routes[apiName].Debug == true {
		return this.routes[apiName].Path
	}
	host := this.host
	if this.accessHost != `` {
		host = this.accessHost
	}
	port := this.port
	if this.accessPort != 0 {
		port = this.accessPort
	}
	return fmt.Sprintf(`http://%s:%s%s%s`, host, port, this.path, this.routes[apiName].Path)
}
