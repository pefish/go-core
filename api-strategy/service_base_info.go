package api_strategy

import (
	"fmt"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/logger"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
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

func (this *ServiceBaseInfoStrategyClass) GetDescription() string {
	return `get service base info`
}

func (this *ServiceBaseInfoStrategyClass) GetErrorCode() uint64 {
	return go_error.INTERNAL_ERROR_CODE
}

func (this *ServiceBaseInfoStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	apiMsg := fmt.Sprintf(`%s %s %s`, out.Ctx.RemoteAddr(), out.Ctx.Path(), out.Ctx.Method())
	logger.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateCtxValuesErrorMsg(out.Ctx, `apiMsg`, apiMsg)
	logger.Logger.DebugF(`UrlParams: %#v`, out.Ctx.URLParams())
	logger.Logger.DebugF(`Headers: %#v`, out.Ctx.Request().Header)
	//rawData, _ := ioutil.ReadAll(out.Ctx.Request().Body) // body只能被拿一次
	//logger.Logger.DebugF(`Body: %#v`, rawData)

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
