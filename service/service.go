package service

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	go_application "github.com/pefish/go-application"
	_type "github.com/pefish/go-core-type/api-session"
	"github.com/pefish/go-core/api"
	api_session "github.com/pefish/go-core/api-session"
	external_service "github.com/pefish/go-core/driver/external-service"
	api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	go_error "github.com/pefish/go-error"
	"github.com/pefish/go-reflect"
	"golang.org/x/net/http2"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type ServiceClass struct {
	name             string                         // 服务名
	description      string                         // 服务描述
	path             string                         // 服务的基础路径
	host             string                         // 服务监听host
	port             uint64                         // 服务监听port
	accessHost       string                         // 服务访问host，没有设置的话使用监听host
	accessPort       uint64                         // 服务访问port，没有设置的话使用监听port
	apis             []*api.Api                     // 服务的所有路由
	healthyCheckFunc func()                         // 健康检查函数
	registeredApi    map[string]map[string]*api.Api // 所有注册了的api。path->method->api

	Mux    *mux.Router
	stopWg sync.WaitGroup
}

// New Service instance
func NewService(name string) *ServiceClass {
	svc := &ServiceClass{
		registeredApi: make(map[string]map[string]*api.Api),
		apis:          make([]*api.Api, 0, 20),
	}
	svc.SetName(name)
	api_strategy.GlobalApiStrategyDriverInstance.Register(api_strategy.GlobalStrategyData{
		Strategy: &global_api_strategy.ServiceBaseInfoApiStrategyInstance,
	})
	api_strategy.GlobalApiStrategyDriverInstance.Register(api_strategy.GlobalStrategyData{
		Strategy: &global_api_strategy.ParamValidateStrategyInstance,
	})
	return svc
}

// Default Service instance
var Service = NewService(`default`)

func (serviceInstance *ServiceClass) GetInterval() time.Duration {
	return 0
}

func (serviceInstance *ServiceClass) SetRoutes(routes ...[]*api.Api) {
	for _, route := range routes {
		serviceInstance.apis = append(serviceInstance.apis, route...)
	}
}

func (serviceInstance *ServiceClass) AddRoute(routes ...*api.Api) {
	serviceInstance.apis = append(serviceInstance.apis, routes...)
}

func (serviceInstance *ServiceClass) SetPath(path string) {
	serviceInstance.path = path
}

func (serviceInstance *ServiceClass) SetName(name string) {
	serviceInstance.name = name
}

func (serviceInstance *ServiceClass) GetHost() string {
	return serviceInstance.host
}

func (serviceInstance *ServiceClass) SetHost(host string) {
	serviceInstance.host = host
}

func (serviceInstance *ServiceClass) GetPort() uint64 {
	return serviceInstance.port
}

func (serviceInstance *ServiceClass) SetPort(port uint64) {
	serviceInstance.port = port
}

func (serviceInstance *ServiceClass) GetAccessHost() string {
	return serviceInstance.accessHost
}

func (serviceInstance *ServiceClass) SetAccessHost(accessHost string) {
	serviceInstance.accessHost = accessHost
}

func (serviceInstance *ServiceClass) GetAccessPort() uint64 {
	return serviceInstance.accessPort
}

func (serviceInstance *ServiceClass) SetAccessPort(accessPort uint64) {
	serviceInstance.accessPort = accessPort
}

func (serviceInstance *ServiceClass) SetDescription(desc string) {
	serviceInstance.description = desc
}

func (serviceInstance *ServiceClass) SetHealthyCheckFunc(func_ func()) *ServiceClass {
	serviceInstance.healthyCheckFunc = func_
	return serviceInstance
}

func (serviceInstance *ServiceClass) GetName() string {
	return serviceInstance.name
}

func (serviceInstance *ServiceClass) GetDescription() string {
	return serviceInstance.description
}

func (serviceInstance *ServiceClass) GetPath() string {
	return serviceInstance.path
}

func (serviceInstance *ServiceClass) GetApis() []*api.Api {
	return serviceInstance.apis
}

func (serviceInstance *ServiceClass) Stop() error {
	serviceInstance.stopWg.Wait()
	return nil
}

func (serviceInstance *ServiceClass) Run(ctx context.Context) error {
	defer func() {
		go_application.Application.Exit()
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	external_service.ExternalServiceDriverInstance.Startup() // 启动外接服务驱动
	logger.LoggerDriverInstance.Startup()                    // 启动日志驱动
	api_strategy.GlobalApiStrategyDriverInstance.Startup()   // 启动外接全局前置处理器驱动

	serviceInstance.buildRoutes()
	host := serviceInstance.host
	if host == `` {
		host = `0.0.0.0`
	}

	addr := host + `:` + go_reflect.Reflect.ToString(serviceInstance.port)
	logger.LoggerDriverInstance.Logger.InfoF(`server started!!! http://%s`, addr)

	for apiPath, map_ := range serviceInstance.registeredApi {
		serviceInstance.Mux.HandleFunc(apiPath, api.WrapJson(map_))
		for method, api_ := range map_ {
			logger.LoggerDriverInstance.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, method, apiPath, api_.Description))
		}
	}
	s := &http.Server{
		Addr:    addr,
		Handler: serviceInstance.Mux,
	}
	err := http2.ConfigureServer(s, &http2.Server{}) // 可以使用http2协议
	if err != nil {
		panic(err)
	}
	go func() {
		serviceInstance.stopWg.Add(1)
		defer serviceInstance.stopWg.Done()
		err := s.ListenAndServe()
		if err != nil {
			logger.LoggerDriverInstance.Logger.Error(err)
		}
	}()
	select {
	case <-ctx.Done():
		s.Shutdown(context.Background())
	}
	return nil
}

func (serviceInstance *ServiceClass) buildRoutes() {
	// healthz
	var healthApi = &api.Api{
		Description:            "健康检查api",
		Path:                   "/healthz",
		IgnoreRootPath:         true,
		IgnoreGlobalStrategies: true,
		Method:                 api_session.ApiMethod_All,
		Controller: func(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
			defer func() {
				if err := recover(); err != nil {
					logger.LoggerDriverInstance.Logger.Error(err)
					apiSession.SetStatusCode(api_session.StatusCode_InternalServerError)
					apiSession.WriteText(`not ok`)
				}
			}()
			global_api_strategy.ServiceBaseInfoApiStrategyInstance.Execute(apiSession, nil)

			if serviceInstance.healthyCheckFunc != nil {
				serviceInstance.healthyCheckFunc()
			}

			apiSession.SetStatusCode(api_session.StatusCode_OK)
			apiSession.WriteText(`ok`)
			return nil, nil
		},
		ParamType: global_api_strategy.ALL_TYPE,
	}

	// 处理未知路由
	var unknownPathApi = &api.Api{
		Description:            "404 not found",
		Path:                   "/",
		IgnoreRootPath:         true,
		IgnoreGlobalStrategies: true,
		Method:                 api_session.ApiMethod_All,
		Controller: func(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
			global_api_strategy.ServiceBaseInfoApiStrategyInstance.Execute(apiSession, nil)

			apiSession.SetStatusCode(api_session.StatusCode_NotFound)
			logger.LoggerDriverInstance.Logger.DebugF("api not found. request path: %s, request method: %s", apiSession.Path(), apiSession.Method())
			apiSession.WriteText(`Not Found`)
			return nil, nil
		},
		ParamType: global_api_strategy.ALL_TYPE,
	}

	serviceInstance.AddRoute(healthApi, unknownPathApi) // 添加缺省api

	serviceInstance.Mux = mux.NewRouter()
	for _, apiObject := range serviceInstance.GetApis() {
		// 得到apiPath
		apiPath := serviceInstance.path + apiObject.Path
		if apiObject.IgnoreRootPath == true {
			apiPath = apiObject.Path
		}
		method := apiObject.Method

		// 挂载处理器
		if apiObject.Controller != nil {
			if serviceInstance.registeredApi[apiPath] == nil {
				serviceInstance.registeredApi[apiPath] = map[string]*api.Api{
					string(method): apiObject,
				}
			} else {
				if serviceInstance.registeredApi[apiPath][string(method)] == nil {
					serviceInstance.registeredApi[apiPath][string(method)] = apiObject
				}
			}
		}
	}
}
