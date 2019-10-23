package service

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-core/service-driver"
	"github.com/pefish/go-reflect"
	"io/ioutil"
	"reflect"
)

type BaseServiceClass struct {
	name             string                                // 服务名
	description      string                                // 服务描述
	path             string                                // 服务的基础路径
	host             string                                // 服务监听host
	port             uint64                                // 服务监听port
	accessHost       string                                // 服务访问host，没有设置的话使用监听host
	accessPort       uint64                                // 服务访问port，没有设置的话使用监听port
	routes           map[string]*api_channel_builder.Route // 服务的所有路由
	globalStrategies []GlobalStrategyStruct                // 全局的也就是每个api的前置处理器
	App              *iris.Application                     // iris实例
	healthyCheckFunc func()                                // 健康检查函数
}

type GlobalStrategyStruct struct {
	Strategy api_channel_builder.InterfaceStrategy
	Param    interface{}
}

func (this *BaseServiceClass) SetRoutes(routes ...map[string]*api_channel_builder.Route) {
	this.routes = map[string]*api_channel_builder.Route{}
	for n, route := range routes {
		for k, v := range route {
			this.routes[go_reflect.Reflect.ToString(n)+`_`+k] = v
		}
	}
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

func (this *BaseServiceClass) SetHealthyCheckFunc(func_ func()) InterfaceService {
	this.healthyCheckFunc = func_
	return this
}

func (this *BaseServiceClass) AddGlobalStrategy(strategy api_channel_builder.InterfaceStrategy, param interface{}) InterfaceService {
	if this.globalStrategies == nil {
		this.globalStrategies = []GlobalStrategyStruct{}
	}
	this.globalStrategies = append(this.globalStrategies, GlobalStrategyStruct{
		Strategy: strategy,
		Param:    param,
	})
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

func (this *BaseServiceClass) GetRoutes() map[string]*api_channel_builder.Route {
	return this.routes
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
	service_driver.ServiceDriver.Init() // 初始化外接服务驱动
	err := this.App.Run(iris.Addr(host+`:`+go_reflect.Reflect.ToString(this.port)), iris.WithConfiguration(irisConfig))
	if err != nil {
		panic(err)
	}
}

func (this *BaseServiceClass) printRoutes() {
	for _, route := range this.routes {
		apiPath := this.path + route.Path
		if route.IgnoreRootPath == true {
			apiPath = route.Path
		}
		logger.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, route.Method, apiPath, route.Description))
	}
}

func (this *BaseServiceClass) buildRoutes() {
	this.App = iris.New()

	this.routes[`healthy_check`] = &api_channel_builder.Route{
		Description:    "健康检查api",
		Path:           "/healthz",
		Method:         "ALL",
		IgnoreRootPath: true,
		Controller: func(apiContext *api_session.ApiSessionClass) interface{} {
			defer func() {
				if err := recover(); err != nil {
					logger.Logger.Error(err)
					apiContext.Ctx.StatusCode(iris.StatusInternalServerError)
					apiContext.Ctx.Text(`not ok`)
				}
			}()
			if this.healthyCheckFunc != nil {
				this.healthyCheckFunc()
			}

			apiContext.Ctx.StatusCode(iris.StatusOK)
			logger.Logger.Debug(`I am healthy`)
			apiContext.Ctx.Text(`ok`)
			return nil
		},
	}

	for name, route := range this.GetRoutes() {
		var apiChannelBuilder = api_channel_builder.NewApiChannelBuilder()
		apiChannelBuilder.Inject(api_strategy.CorsApiStrategy.GetName(), api_channel_builder.InjectObject{
			Func: api_strategy.CorsApiStrategy.Execute,
			This: &api_strategy.CorsApiStrategy,
		})
		apiChannelBuilder.Inject(api_strategy.ServiceBaseInfoApiStrategy.GetName(), api_channel_builder.InjectObject{
			Func: api_strategy.ServiceBaseInfoApiStrategy.Execute,
			Param: api_strategy.ServiceBaseInfoParam{
				RouteName: name,
			},
			This: &api_strategy.ServiceBaseInfoApiStrategy,
		})
		apiChannelBuilder.Inject(api_strategy.ParamValidateStrategy.GetName(), api_channel_builder.InjectObject{
			Func: api_strategy.ParamValidateStrategy.Execute,
			Param: api_strategy.ParamValidateParam{
				Param: route.Params,
			},
			Route: route,
			This:  &api_strategy.ParamValidateStrategy,
		})
		for _, globalStrategy := range this.globalStrategies {
			apiChannelBuilder.Inject(globalStrategy.Strategy.GetName(), api_channel_builder.InjectObject{
				Func:  globalStrategy.Strategy.Execute,
				This:  globalStrategy.Strategy,
				Param: globalStrategy.Param,
				Route: route,
			})
		}
		if route.Strategies != nil {
			for _, strategyRoute := range route.Strategies {
				if strategyRoute.Disable {
					continue
				}
				apiChannelBuilder.Inject(strategyRoute.Strategy.GetName(), api_channel_builder.InjectObject{
					Func:  strategyRoute.Strategy.Execute,
					Param: strategyRoute.Param,
					Route: route,
					This:  strategyRoute.Strategy,
				})
			}
		}
		apiPath := this.path + route.Path
		if route.IgnoreRootPath == true {
			apiPath = route.Path
		}
		if route.Method == `` {
			route.Method = `ALL`
		}
		if route.Controller != nil {
			this.App.AllowMethods(iris.MethodOptions).Handle(route.Method, apiPath, apiChannelBuilder.WrapJson(route.Controller))
		}
	}

	// 处理未知路由
	var apiChannelBuilder = api_channel_builder.NewApiChannelBuilder()
	apiChannelBuilder.Inject(api_strategy.CorsApiStrategy.GetName(), api_channel_builder.InjectObject{
		Func: api_strategy.CorsApiStrategy.Execute,
		This: &api_strategy.CorsApiStrategy,
	})
	apiChannelBuilder.Inject(api_strategy.ServiceBaseInfoApiStrategy.GetName(), api_channel_builder.InjectObject{
		Func: api_strategy.ServiceBaseInfoApiStrategy.Execute,
		Param: api_strategy.ServiceBaseInfoParam{
			RouteName: `*`,
		},
		This: &api_strategy.ServiceBaseInfoApiStrategy,
	})
	this.App.AllowMethods(iris.MethodOptions).Handle(``, `/*`, apiChannelBuilder.WrapJson(func(apiContext *api_session.ApiSessionClass) interface{} {
		rawData, _ := ioutil.ReadAll(apiContext.Ctx.Request().Body)
		logger.Logger.DebugF(`Body: %s`, string(rawData))
		apiContext.Ctx.StatusCode(iris.StatusNotFound)
		logger.Logger.Debug(`api not found`)
		apiContext.Ctx.Text(`Not Found`)
		return nil
	}))
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
