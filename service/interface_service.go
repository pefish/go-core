package service

import (
	"github.com/pefish/go-core/api-channel-builder"
)

type InterfaceService interface {
	// opts[0]: map[string]interface{}
	Init(opts ...interface{}) InterfaceService
	// 设置安全检查函数
	SetHealthyCheck(func_ func()) InterfaceService
	// 注入处理器
	Use(key string, func_ api_channel_builder.InjectObject) InterfaceService
	// 获取所有路由
	GetRoutes() map[string]*api_channel_builder.Route
	SetRoutes(routes ...map[string]*api_channel_builder.Route)
	// 获取端点路径
	GetPath() string
	SetPath(path string)
	// 获取服务名
	GetName() string
	SetName(name string)
	// 获取Host
	GetHost() string
	SetHost(host string)
	// 获取Port
	GetPort() uint64
	SetPort(port uint64)

	GetAccessHost() string
	SetAccessHost(accessHost string)

	GetAccessPort() uint64
	SetAccessPort(accessPort uint64)
	// 获取服务描述
	GetDescription() string
	SetDescription(desc string)
	// 运行服务
	Run()
	// 获取请求uri
	GetRequestUrl(apiName string) string
	// 发起请求
	Request(apiName string, args ...interface{}) interface{}
	RequestWithErr(apiName string, args ...interface{}) (interface{}, error)
	RequestRawMap(apiName string, args ...interface{}) map[string]interface{}
	RequestForMap(apiName string, args ...interface{}) map[string]interface{}
	RequestForMapWithScan(dest interface{}, apiName string, args ...interface{})
	RequestForSlice(apiName string, args ...interface{}) []map[string]interface{}
	RequestForSliceWithScan(dest interface{}, apiName string, args ...interface{})
}
