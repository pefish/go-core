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
	errorCode uint64
}

var ServiceBaseInfoApiStrategy = ServiceBaseInfoStrategyClass{

}

func (this *ServiceBaseInfoStrategyClass) GetName() string {
	return `serviceBaseInfo`
}

func (this *ServiceBaseInfoStrategyClass) GetDescription() string {
	return `get service base info`
}

func (this *ServiceBaseInfoStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *ServiceBaseInfoStrategyClass) GetErrorCode() uint64 {
	if this.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return go_error.INTERNAL_ERROR_CODE
}

func (this *ServiceBaseInfoStrategyClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())
}

func (this *ServiceBaseInfoStrategyClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, this.GetName())
	apiMsg := fmt.Sprintf(`%s %s %s`, out.GetRemoteAddress(), out.GetPath(), out.GetMethod())
	logger.LoggerDriver.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateSessionErrorMsg(out, `apiMsg`, apiMsg)
	logger.LoggerDriver.Logger.DebugF(`UrlParams: %#v`, out.GetUrlParams())
	logger.LoggerDriver.Logger.DebugF(`Headers: %#v`, out.Request.Header)

	rawData, _ := ioutil.ReadAll(out.Request.Body)
	out.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData)) // 读出来后又新建一个流填进去，使out.Request.Body可以被再次读
	logger.LoggerDriver.Logger.DebugF(`Body: %s`, string(rawData))

	lang := out.GetHeader(`lang`)
	if lang == `` {
		lang = `zh-CN`
	}
	out.Lang = lang

	clientType := out.GetHeader(`client_type`)
	if clientType == `` {
		clientType = `web`
	}
	out.ClientType = clientType
}
