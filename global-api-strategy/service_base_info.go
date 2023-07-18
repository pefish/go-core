package global_api_strategy

import (
	"bytes"
	"fmt"
	api_session "github.com/pefish/go-core-type/api-session"
	global_api_strategy "github.com/pefish/go-core-type/global-api-strategy"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
	"io"
)

type ServiceBaseInfoStrategy struct {
	errorCode uint64
	errorMsg  string
}

var ServiceBaseInfoApiStrategyInstance = ServiceBaseInfoStrategy{}

func (sbis *ServiceBaseInfoStrategy) GetName() string {
	return `ServiceBaseInfoStrategy`
}

func (sbis *ServiceBaseInfoStrategy) GetDescription() string {
	return `get service base info`
}

func (sbis *ServiceBaseInfoStrategy) SetErrorCode(code uint64) global_api_strategy.IGlobalApiStrategy {
	sbis.errorCode = code
	return sbis
}

func (sbis *ServiceBaseInfoStrategy) GetErrorCode() uint64 {
	if sbis.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return sbis.errorCode
}

func (sbis *ServiceBaseInfoStrategy) SetErrorMsg(msg string) global_api_strategy.IGlobalApiStrategy {
	sbis.errorMsg = msg
	return sbis
}

func (sbis *ServiceBaseInfoStrategy) GetErrorMsg() string {
	return sbis.errorMsg
}

func (sbis *ServiceBaseInfoStrategy) Init(param interface{}) {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init`, sbis.GetName())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init defer`, sbis.GetName())
}

func (sbis *ServiceBaseInfoStrategy) Execute(out api_session.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s trigger`, sbis.GetName())
	apiMsg := fmt.Sprintf(`%s %s %s`, out.RemoteAddress(), out.Path(), out.Method())
	logger.LoggerDriverInstance.Logger.DebugF(`---------------- %s ----------------`, apiMsg)
	util.UpdateSessionErrorMsg(out, `apiMsg`, apiMsg)
	logger.LoggerDriverInstance.Logger.DebugF(`UrlParams: %#v`, out.UrlParams())
	logger.LoggerDriverInstance.Logger.DebugF(`Headers: %#v`, out.Request().Header)

	rawData, _ := io.ReadAll(out.Request().Body)
	out.Request().Body = io.NopCloser(bytes.NewBuffer(rawData)) // 读出来后又新建一个流填进去，使out.request.Body可以被再次读
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
