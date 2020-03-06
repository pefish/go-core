package service

import (
	"fmt"
	go_application "github.com/pefish/go-application"
	"github.com/pefish/go-core/api"
	api_session "github.com/pefish/go-core/api-session"
	external_service "github.com/pefish/go-core/driver/external-service"
	api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
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
	healthyCheckFunc func()            // 健康检查函数

	Mux *http.ServeMux
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

	// 执行各个全局策略的初始化函数
	for _, globalStrategy := range api_strategy.GlobalApiStrategyDriver.GlobalStrategies {
		if !globalStrategy.Disable {
			globalStrategy.Strategy.Init(globalStrategy.Param)
			go globalStrategy.Strategy.InitAsync(globalStrategy.Param, go_application.OnTerminated)
		}
	}

	this.buildRoutes()
	host := this.host
	if host == `` {
		host = `0.0.0.0`
	}

	addr := host + `:` + go_reflect.Reflect.MustToString(this.port)
	logger.LoggerDriver.Logger.InfoF(`server started!!! http://%s`, addr)
	err := http.ListenAndServe(addr, this.Mux)
	if err != nil {
		panic(err)
	}
}

func (this *ServiceClass) buildRoutes() {
	this.Mux = http.NewServeMux()
	for _, apiObject := range this.GetApis() {
		// 得到apiPath
		apiPath := this.path + apiObject.Path
		if apiObject.IgnoreRootPath == true {
			apiPath = apiObject.Path
		}
		// 方法为空字符串就是All
		method := apiObject.Method
		// 挂载处理器
		if apiObject.Controller != nil {
			this.Mux.HandleFunc(apiPath, apiObject.WrapJson(method, apiObject.Controller))
			logger.LoggerDriver.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, method, apiPath, apiObject.Description))
		}
	}

	// healthz
	var healthApiObject = api.NewApi()
	this.Mux.HandleFunc(`/healthz`, healthApiObject.WrapJson(api_session.ApiMethod_All, func(apiSession *api_session.ApiSessionClass) interface{} {
		defer func() {
			if err := recover(); err != nil {
				logger.LoggerDriver.Logger.Error(err)
				apiSession.SetStatusCode(api_session.StatusCode_InternalServerError)
				apiSession.WriteText(`not ok`)
			}
		}()
		if this.healthyCheckFunc != nil {
			this.healthyCheckFunc()
		}

		apiSession.SetStatusCode(api_session.StatusCode_OK)
		apiSession.WriteText(`ok`)
		return nil
	}))
	logger.LoggerDriver.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, api_session.ApiMethod_All, `/healthz`, `健康检查api`))
	// 处理未知路由
	var apiObject = api.NewApi()
	this.Mux.HandleFunc(`/*`, apiObject.WrapJson(api_session.ApiMethod_All, func(apiSession *api_session.ApiSessionClass) interface{} {
		rawData, _ := ioutil.ReadAll(apiSession.Request.Body)
		logger.LoggerDriver.Logger.DebugF(`Body: %s`, string(rawData))
		apiSession.SetStatusCode(api_session.StatusCode_NotFound)
		logger.LoggerDriver.Logger.Debug(`api not found`)
		apiSession.WriteText(`Not Found`)
		return nil
	}))
	logger.LoggerDriver.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, api_session.ApiMethod_All, `/*`, `404 not found`))
}
