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

type ServiceBaseInfoStrategy struct {
	errorCode uint64
}

var ServiceBaseInfoApiStrategyInstance = ServiceBaseInfoStrategy{

}

func (serviceBaseInfo *ServiceBaseInfoStrategy) GetName() string {
	return `serviceBaseInfo`
}

func (serviceBaseInfo *ServiceBaseInfoStrategy) GetDescription() string {
	return `get service base info`
}

func (serviceBaseInfo *ServiceBaseInfoStrategy) SetErrorCode(code uint64) {
	serviceBaseInfo.errorCode = code
}

func (serviceBaseInfo *ServiceBaseInfoStrategy) GetErrorCode() uint64 {
	if serviceBaseInfo.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return go_error.INTERNAL_ERROR_CODE
}

func (serviceBaseInfo *ServiceBaseInfoStrategy) Init(param interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init`, serviceBaseInfo.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init defer`, serviceBaseInfo.GetName())
}

func (serviceBaseInfo *ServiceBaseInfoStrategy) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s trigger`, serviceBaseInfo.GetName())
	apiMsg := fmt.Sprintf(`%s %s %s`, out.RemoteAddress(), out.Path(), out.Method())
	logger.LoggerDriverInstance.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateSessionErrorMsg(out, `apiMsg`, apiMsg)
	logger.LoggerDriverInstance.Logger.DebugF(`UrlParams: %#v`, out.UrlParams())
	logger.LoggerDriverInstance.Logger.DebugF(`Headers: %#v`, out.Request().Header)

	rawData, _ := ioutil.ReadAll(out.Request().Body)
	out.Request().Body = ioutil.NopCloser(bytes.NewBuffer(rawData)) // 读出来后又新建一个流填进去，使out.request.Body可以被再次读
	logger.LoggerDriverInstance.Logger.DebugF(`Body: %s`, string(rawData))

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
