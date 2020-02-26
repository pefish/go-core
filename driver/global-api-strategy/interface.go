package global_api_strategy

import (
	api_strategy "github.com/pefish/go-core/api-strategy"
)

type InterfaceGlobalStrategy interface {
	Init(param interface{})  // 同步的初始化函数
	InitAsync(param interface{}, onAppTerminated chan interface{})  // 异步的初始化函数，应用推出init才能退出的场景
	api_strategy.InterfaceStrategy
}
