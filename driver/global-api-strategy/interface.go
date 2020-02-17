package global_api_strategy

import api_session "github.com/pefish/go-core/api-session"

type InterfaceStrategy interface {
	Init(param interface{})  // 同步的初始化函数
	InitAsync(param interface{}, onAppTerminated chan interface{})  // 异步的初始化函数，应用推出init才能退出的场景
	Execute(out *api_session.ApiSessionClass, param interface{})
	GetName() string
	GetDescription() string
	GetErrorCode() uint64
}
