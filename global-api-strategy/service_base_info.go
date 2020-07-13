package global_api_strategy

import (
	"bytes"
	"fmt"
	_type "github.com/pefish/go-core/api-session/type"
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

func (serviceBaseInfo *ServiceBaseInfoStrategyClass) GetName() string {
	return `serviceBaseInfo`
}

func (serviceBaseInfo *ServiceBaseInfoStrategyClass) GetDescription() string {
	return `get service base info`
}

func (serviceBaseInfo *ServiceBaseInfoStrategyClass) SetErrorCode(code uint64) {
	serviceBaseInfo.errorCode = code
}

func (serviceBaseInfo *ServiceBaseInfoStrategyClass) GetErrorCode() uint64 {
	if serviceBaseInfo.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return go_error.INTERNAL_ERROR_CODE
}

func (serviceBaseInfo *ServiceBaseInfoStrategyClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, serviceBaseInfo.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, serviceBaseInfo.GetName())
}

func (serviceBaseInfo *ServiceBaseInfoStrategyClass) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, serviceBaseInfo.GetName())
	apiMsg := fmt.Sprintf(`%s %s %s`, out.RemoteAddress(), out.Path(), out.Method())
	logger.LoggerDriver.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateSessionErrorMsg(out, `apiMsg`, apiMsg)
	logger.LoggerDriver.Logger.DebugF(`UrlParams: %#v`, out.UrlParams())
	logger.LoggerDriver.Logger.DebugF(`Headers: %#v`, out.Request().Header)

	rawData, _ := ioutil.ReadAll(out.Request().Body)
	out.Request().Body = ioutil.NopCloser(bytes.NewBuffer(rawData)) // 读出来后又新建一个流填进去，使out.request.Body可以被再次读
	logger.LoggerDriver.Logger.DebugF(`Body: %s`, string(rawData))

	lang := out.Header(`lang`)
	if lang == `` {
		lang = `zh-CN`
	}
	out.SetLang(lang)

	clientType := out.Header(`client_type`)
	if clientType == `` {
		clientType = `web`
	}
	out.SetClientType(clientType)

	return nil
}
