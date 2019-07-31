package api_strategy

import (
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
)

type ServiceBaseInfoStrategyClass struct {
}

var ServiceBaseInfoApiStrategy = ServiceBaseInfoStrategyClass{}

type ServiceBaseInfoParam struct {
	RouteName string
}

func (this *ServiceBaseInfoStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	newParam := param.(ServiceBaseInfoParam)
	out.RouteName = newParam.RouteName

	lang := ctx.GetHeader(`lang`)
	if lang == `` {
		lang = `zh-CN`
	}
	out.Lang = lang

	clientType := ctx.GetHeader(`client_type`)
	if clientType == `` {
		clientType = `web`
	}
	out.ClientType = clientType
}
