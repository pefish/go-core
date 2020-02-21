package service

import (
	"fmt"
	"github.com/kataras/iris"
	go_application "github.com/pefish/go-application"
	"github.com/pefish/go-core/api"
	api_session "github.com/pefish/go-core/api-session"
	api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	external_service "github.com/pefish/go-core/driver/external-service"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-reflect"
	"io/ioutil"
	"net/http"
	"runtime"
)

type ServiceClass struct {
	name             string            // 服务名
	description      string            // 服务描述
	path             string            // 服务的基础路径
	host             string            // 服务监听host
	port             uint64            // 服务监听port
	accessHost       string            // 服务访问host，没有设置的话使用监听host
	accessPort       uint64            // 服务访问port，没有设置的话使用监听port
	apis             []*api.Api        // 服务的所有路由
	App              *iris.Application // iris实例
	healthyCheckFunc func()            // 健康检查函数

}

func (this *ServiceClass) SetRoutes(routes ...[]*api.Api) {
	this.apis = []*api.Api{}
	for _, route := range routes {
		this.apis = append(this.apis, route...)
	}
}

func (this *ServiceClass) SetPath(path string) {
	this.path = path
}

func (this *ServiceClass) SetName(name string) {
	this.name = name
}

func (this *ServiceClass) GetHost() string {
	return this.host
}

func (this *ServiceClass) SetHost(host string) {
	this.host = host
}

func (this *ServiceClass) GetPort() uint64 {
	return this.port
}

func (this *ServiceClass) SetPort(port uint64) {
	this.port = port
}

func (this *ServiceClass) GetAccessHost() string {
	return this.accessHost
}

func (this *ServiceClass) SetAccessHost(accessHost string) {
	this.accessHost = accessHost
}

func (this *ServiceClass) GetAccessPort() uint64 {
	return this.accessPort
}

func (this *ServiceClass) SetAccessPort(accessPort uint64) {
	this.accessPort = accessPort
}

func (this *ServiceClass) SetDescription(desc string) {
	this.description = desc
}

func (this *ServiceClass) SetHealthyCheckFunc(func_ func()) *ServiceClass {
	this.healthyCheckFunc = func_
	return this
}

func (this *ServiceClass) GetName() string {
	return this.name
}

func (this *ServiceClass) GetDescription() string {
	return this.description
}

func (this *ServiceClass) GetPath() string {
	return this.path
}

func (this *ServiceClass) GetApis() []*api.Api {
	return this.apis
}

func (this *ServiceClass) Run() {
	defer func() {
		close(go_application.OnTerminated) // 关闭通道。实现广播让所有订阅此通道都得到消息
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	external_service.ExternalServiceDriver.Startup() // 启动外接服务驱动
	logger.LoggerDriver.Startup()                    // 启动日志驱动
	api_strategy.GlobalApiStrategyDriver.Startup()   // 启动外接全局前置处理器驱动

	// 执行各个策略的初始化函数
	for _, globalStrategy := range api_strategy.GlobalApiStrategyDriver.GlobalStrategies {
		if !globalStrategy.Disable {
			globalStrategy.Strategy.Init(globalStrategy.Param)
			go globalStrategy.Strategy.InitAsync(globalStrategy.Param, go_application.OnTerminated)
		}
	}

	this.buildRoutes()
	irisConfig := iris.Configuration{}
	irisConfig.RemoteAddrHeaders = map[string]bool{
		`X-Forwarded-For`: true,
	}
	irisConfig.DisableBodyConsumptionOnUnmarshal = true // 使ReadJson后Body内容可以反复读
	host := this.host
	if host == `` {
		host = `0.0.0.0`
	}

	addr := host + `:` + go_reflect.Reflect.MustToString(this.port)
	err := this.App.Run(func(application *iris.Application) error {
		return this.App.NewHost(&http.Server{
			Addr: addr}).Configure().ListenAndServe()
	}, iris.WithConfiguration(irisConfig))

	if err != nil && err.Error() != `http: Server closed` {
		panic(err)
	}
}

func (this *ServiceClass) buildRoutes() {
	this.App = iris.New()
	this.apis = append(this.apis, &api.Api{
		Description:    "健康检查api",
		Path:           "/healthz",
		Method:         "ALL",
		IgnoreRootPath: true,
		Controller: func(apiContext *api_session.ApiSessionClass) interface{} {
			defer func() {
				if err := recover(); err != nil {
					logger.LoggerDriver.Logger.Error(err)
					apiContext.Ctx.StatusCode(iris.StatusInternalServerError)
					apiContext.Ctx.Text(`not ok`)
				}
			}()
			if this.healthyCheckFunc != nil {
				this.healthyCheckFunc()
			}

			apiContext.Ctx.StatusCode(iris.StatusOK)
			apiContext.Ctx.Text(`ok`)
			return nil
		},
	})

	for _, apiObject := range this.GetApis() {
		// 注入全局前置处理器
		apiObject.Strategies = append(api_strategy.GlobalApiStrategyDriver.GlobalStrategies, apiObject.Strategies...)
		// 得到apiPath
		apiPath := this.path + apiObject.Path
		if apiObject.IgnoreRootPath == true {
			apiPath = apiObject.Path
		}
		// 方法为空字符串就是All
		method := apiObject.Method
		// 挂载处理器
		if apiObject.Controller != nil {
			this.App.AllowMethods(iris.MethodOptions).Handle(method, apiPath, apiObject.WrapJson(apiObject.Controller))
			methodText := method
			if method == `` {
				methodText = `ALL`
			}
			logger.LoggerDriver.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, methodText, apiPath, apiObject.Description))
		}
	}

	// 处理未知路由
	var apiObject = api.NewApi()
	apiObject.Strategies = append(api_strategy.GlobalApiStrategyDriver.GlobalStrategies, apiObject.Strategies...)
	this.App.AllowMethods(iris.MethodOptions).Handle(``, `/*`, apiObject.WrapJson(func(apiContext *api_session.ApiSessionClass) interface{} {
		rawData, _ := ioutil.ReadAll(apiContext.Ctx.Request().Body)
		logger.LoggerDriver.Logger.DebugF(`Body: %s`, string(rawData))
		apiContext.Ctx.StatusCode(iris.StatusNotFound)
		logger.LoggerDriver.Logger.Debug(`api not found`)
		apiContext.Ctx.Text(`Not Found`)
		return nil
	}))
	logger.LoggerDriver.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, `ALL`, `/*`, `404 not found`))
}
