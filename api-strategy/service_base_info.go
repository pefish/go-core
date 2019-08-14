package api_strategy

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
)

type ServiceBaseInfoStrategyClass struct {

}

var ServiceBaseInfoApiStrategy = ServiceBaseInfoStrategyClass{

}

type ServiceBaseInfoParam struct {
	RouteName string
}

func (this *ServiceBaseInfoStrategyClass) GetName() string {
	return `serviceBaseInfo`
}

func (this *ServiceBaseInfoStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	newParam := param.(ServiceBaseInfoParam)
	out.RouteName = newParam.RouteName

	lang := out.Ctx.GetHeader(`lang`)
	if lang == `` {
		lang = `zh-CN`
	}
	out.Lang = lang

	clientType := out.Ctx.GetHeader(`client_type`)
	if clientType == `` {
		clientType = `web`
	}
	out.ClientType = clientType
}
