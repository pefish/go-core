package global_api_strategy

import (
	"bytes"
	"fmt"
	"io"

	api_session "github.com/pefish/go-core-type/api-session"
	api_strategy "github.com/pefish/go-core-type/api-strategy"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-core/util"
	go_error "github.com/pefish/go-error"
)

type ServiceBaseInfoStrategy struct {
	errorCode uint64
	errorMsg  string
}

var ServiceBaseInfoApiStrategyInstance = ServiceBaseInfoStrategy{}

func (sbis *ServiceBaseInfoStrategy) Name() string {
	return `ServiceBaseInfoStrategy`
}

func (sbis *ServiceBaseInfoStrategy) Description() string {
	return `get service base info`
}

func (sbis *ServiceBaseInfoStrategy) SetErrorCode(code uint64) api_strategy.IApiStrategy {
	sbis.errorCode = code
	return sbis
}

func (sbis *ServiceBaseInfoStrategy) ErrorCode() uint64 {
	if sbis.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return sbis.errorCode
}

func (sbis *ServiceBaseInfoStrategy) SetErrorMsg(msg string) api_strategy.IApiStrategy {
	sbis.errorMsg = msg
	return sbis
}

func (sbis *ServiceBaseInfoStrategy) ErrorMsg() string {
	return sbis.errorMsg
}

func (sbis *ServiceBaseInfoStrategy) Init(param interface{}) api_strategy.IApiStrategy {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init`, sbis.Name())
	defer logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s Init defer`, sbis.Name())
	return sbis
}

func (sbis *ServiceBaseInfoStrategy) Execute(out api_session.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriverInstance.Logger.DebugF(`api-strategy %s trigger`, sbis.Name())
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
