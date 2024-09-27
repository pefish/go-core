package service

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/pefish/go-core/api"
	external_service "github.com/pefish/go-core/driver/external-service"
	api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	go_format "github.com/pefish/go-format"
	i_logger "github.com/pefish/go-interface/i-logger"
	"golang.org/x/net/http2"
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

	Mux *mux.Router
}

// New Service instance
func NewService(name string) *ServiceClass {
	svc := &ServiceClass{
		registeredApi: make(map[string]map[string]*api.Api),
		apis:          make([]*api.Api, 0, 20),
	}
	svc.SetName(name)
	api_strategy.GlobalApiStrategyDriverInstance.Register(api_strategy.GlobalStrategyData{
		Strategy: global_api_strategy.ServiceBaseInfoStrategyInstance,
	})
	api_strategy.GlobalApiStrategyDriverInstance.Register(api_strategy.GlobalStrategyData{
		Strategy: global_api_strategy.ParamValidateStrategyInstance,
	})
	return svc
}

// Default Service instance
var Service = NewService(`default`)

func (serviceInstance *ServiceClass) Interval() time.Duration {
	return 0
}

func (serviceInstance *ServiceClass) SetRoutes(routes ...[]*api.Api) {
	for _, apis := range routes {
		serviceInstance.AddApis(apis...)
	}
}

func (serviceInstance *ServiceClass) AddApis(apis ...*api.Api) {
	serviceInstance.apis = append(serviceInstance.apis, apis...)
}

func (serviceInstance *ServiceClass) SetPath(path string) {
	serviceInstance.path = path
}

func (serviceInstance *ServiceClass) SetName(name string) {
	serviceInstance.name = name
}

func (serviceInstance *ServiceClass) Host() string {
	return serviceInstance.host
}

func (serviceInstance *ServiceClass) SetHost(host string) {
	serviceInstance.host = host
}

func (serviceInstance *ServiceClass) Port() uint64 {
	return serviceInstance.port
}

func (serviceInstance *ServiceClass) SetPort(port uint64) {
	serviceInstance.port = port
}

func (serviceInstance *ServiceClass) AccessHost() string {
	return serviceInstance.accessHost
}

func (serviceInstance *ServiceClass) SetAccessHost(accessHost string) {
	serviceInstance.accessHost = accessHost
}

func (serviceInstance *ServiceClass) AccessPort() uint64 {
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

func (serviceInstance *ServiceClass) Name() string {
	return serviceInstance.name
}

func (serviceInstance *ServiceClass) Description() string {
	return serviceInstance.description
}

func (serviceInstance *ServiceClass) Path() string {
	return serviceInstance.path
}

func (serviceInstance *ServiceClass) Apis() []*api.Api {
	return serviceInstance.apis
}

func (serviceInstance *ServiceClass) Stop() error {
	return nil
}

func (serviceInstance *ServiceClass) Logger() i_logger.ILogger {
	return logger.LoggerDriverInstance.Logger
}

func (serviceInstance *ServiceClass) Init(ctx context.Context) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	external_service.ExternalServiceDriverInstance.Startup() // 启动外接服务驱动
	logger.LoggerDriverInstance.Startup()                    // 启动日志驱动
	api_strategy.GlobalApiStrategyDriverInstance.Startup()   // 启动外接全局前置处理器驱动

	serviceInstance.buildRoutes()

	return nil
}

func (serviceInstance *ServiceClass) Run(ctx context.Context) error {
	host := serviceInstance.host
	if host == `` {
		host = `0.0.0.0`
	}

	addr := host + `:` + go_format.ToString(serviceInstance.port)
	logger.LoggerDriverInstance.Logger.InfoF(`Server started!!! http://%s`, addr)

	for apiPath, map_ := range serviceInstance.registeredApi {
		serviceInstance.Mux.HandleFunc(apiPath, api.WrapJson(map_))
	}
	s := &http.Server{
		Addr:    addr,
		Handler: serviceInstance.Mux,
	}
	err := http2.ConfigureServer(s, &http2.Server{}) // 可以使用 http2 协议
	if err != nil {
		return err
	}

	exited := make(chan bool)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.LoggerDriverInstance.Logger.Error(err)
			exited <- true
		}
	}()

	select {
	case <-ctx.Done():
		s.Shutdown(ctx)
	case <-exited:
	}
	return nil
}

func (serviceInstance *ServiceClass) buildRoutes() {
	serviceInstance.AddApis(api.New404Api()) // 添加缺省api

	serviceInstance.Mux = mux.NewRouter()
	for _, apiObject := range serviceInstance.Apis() {
		// 得到apiPath
		apiPath := serviceInstance.path + apiObject.Path()
		if apiObject.IsIgnoreRootPath() == true {
			apiPath = apiObject.Path()
		}
		method := apiObject.Method()

		// 挂载处理器
		if apiObject.ControllerFunc() != nil {
			if serviceInstance.registeredApi[apiPath] == nil {
				serviceInstance.registeredApi[apiPath] = map[string]*api.Api{
					string(method): apiObject,
				}
			} else {
				if serviceInstance.registeredApi[apiPath][string(method)] == nil {
					serviceInstance.registeredApi[apiPath][string(method)] = apiObject
				}
			}

			logger.LoggerDriverInstance.Logger.Info(fmt.Sprintf(`--- %s %s %s ---`, method, apiPath, apiObject.Description()))

		}
	}
}
