package global_api_strategy

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-core/util"
	i_core "github.com/pefish/go-interface/i-core"
	t_error "github.com/pefish/go-interface/t-error"
)

type ServiceBaseInfoStrategy struct {
	errorCode uint64
	errorMsg  string
}

var ServiceBaseInfoStrategyInstance = NewServiceBaseInfoStrategy()

func NewServiceBaseInfoStrategy() *ServiceBaseInfoStrategy {
	return &ServiceBaseInfoStrategy{}
}

func (sbis *ServiceBaseInfoStrategy) Name() string {
	return `ServiceBaseInfoStrategy`
}

func (sbis *ServiceBaseInfoStrategy) Description() string {
	return `get service base info`
}

func (sbis *ServiceBaseInfoStrategy) SetErrorCode(code uint64) i_core.IApiStrategy {
	sbis.errorCode = code
	return sbis
}

func (sbis *ServiceBaseInfoStrategy) ErrorCode() uint64 {
	if sbis.errorCode == 0 {
		return t_error.INTERNAL_ERROR_CODE
	}
	return sbis.errorCode
}

func (sbis *ServiceBaseInfoStrategy) SetErrorMsg(msg string) i_core.IApiStrategy {
	sbis.errorMsg = msg
	return sbis
}

func (sbis *ServiceBaseInfoStrategy) ErrorMsg() string {
	return sbis.errorMsg
}

func (sbis *ServiceBaseInfoStrategy) Execute(out i_core.IApiSession) *t_error.ErrorInfo {
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
