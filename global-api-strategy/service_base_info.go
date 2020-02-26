package global_api_strategy

import (
	"bytes"
	"fmt"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
	"io/ioutil"
)

type ServiceBaseInfoStrategyClass struct {

}

var ServiceBaseInfoApiStrategy = ServiceBaseInfoStrategyClass{

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

func (this *ServiceBaseInfoStrategyClass) InitAsync(param interface{}, onAppTerminated chan interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s InitAsync`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s InitAsync defer`, this.GetName())
}

func (this *ServiceBaseInfoStrategyClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())
}

func (this *ServiceBaseInfoStrategyClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, this.GetName())
	apiMsg := fmt.Sprintf(`%s %s %s`, out.Ctx.RemoteAddr(), out.Ctx.Path(), out.Ctx.Method())
	logger.LoggerDriver.Logger.Debug(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateCtxValuesErrorMsg(out.Ctx, `apiMsg`, apiMsg)
	logger.LoggerDriver.Logger.DebugF(`UrlParams: %#v`, out.Ctx.URLParams())
	logger.LoggerDriver.Logger.DebugF(`Headers: %#v`, out.Ctx.Request().Header)

	rawData, _ := ioutil.ReadAll(out.Ctx.Request().Body)
	if out.Ctx.Application().ConfigurationReadOnly().GetDisableBodyConsumptionOnUnmarshal() {
		out.Ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
	}
	logger.LoggerDriver.Logger.DebugF(`Body: %s`, string(rawData))

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
