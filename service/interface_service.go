package service

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/kataras/iris/context"
)

type InterfaceService interface {
	// opts[0]: map[string]interface{} (apiControllers、)
	Init(opts ...interface{}) InterfaceService
	// 设置安全检查函数
	SetHealthyCheck(func_ func()) InterfaceService
	// 注入处理器
	Use(key string, func_ api_channel_builder.InjectFuncType) InterfaceService
	// 注入全局处理器
	UseGlobal(key string, func_ context.Handler) InterfaceService
	// 获取所有路由
	GetRoutes() map[string]*api_session.Route
	// 获取端点路径
	GetPath() string
	// 获取服务名
	GetName() string
	// 获取服务描述
	GetDescription() string
	// 运行服务
	Run()
	// 获取请求uri
	GetRequestUrl(apiName string) string
	// 发起请求
	Request(apiName string, args ...interface{}) interface{}
	RequestForMap(apiName string, args ...interface{}) map[string]interface{}
	RequestForMapWithScan(dest interface{}, apiName string, args ...interface{})
	RequestForSlice(apiName string, args ...interface{}) []map[string]interface{}
	RequestForSliceWithScan(dest interface{}, apiName string, args ...interface{})
}
