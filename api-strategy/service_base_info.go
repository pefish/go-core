package api_strategy

import (
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
)

type ServiceBaseInfoStrategyClass struct {
	errorCode uint64
}

var ServiceBaseInfoApiStrategy = ServiceBaseInfoStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

type ServiceBaseInfoParam struct {
	RouteName string
}

func (this *ServiceBaseInfoStrategyClass) GetName() string {
	return `serviceBaseInfo`
}

func (this *ServiceBaseInfoStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
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
