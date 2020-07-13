package api_strategy

import (
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-error"
)

type IpFilterStrategyClass struct {
	errorCode uint64
}

var IpFilterStrategy = IpFilterStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

type IpFilterParam struct {
	GetValidIp func(apiSession api_session.InterfaceApiSession) []string
}

func (ipFilter *IpFilterStrategyClass) GetName() string {
	return `ipFilter`
}

func (ipFilter *IpFilterStrategyClass) GetDescription() string {
	return `filter ip`
}

func (ipFilter *IpFilterStrategyClass) SetErrorCode(code uint64) {
	ipFilter.errorCode = code
}

func (ipFilter *IpFilterStrategyClass) GetErrorCode() uint64 {
	return ipFilter.errorCode
}

func (ipFilter *IpFilterStrategyClass) Execute(out api_session.InterfaceApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, ipFilter.GetName())
	if param == nil {
		return &go_error.ErrorInfo{
			InternalErrorMessage: `strategy need param`,
			ErrorMessage: `strategy need param`,
			ErrorCode: ipFilter.errorCode,
		}
	}
	newParam := param.(IpFilterParam)
	if newParam.GetValidIp == nil {
		return nil
	}
	clientIp := out.RemoteAddress()
	allowedIps := newParam.GetValidIp(out)
	for _, ip := range allowedIps {
		if ip == clientIp {
			return nil
		}
	}
	return &go_error.ErrorInfo{
		InternalErrorMessage: `ip is baned`,
		ErrorMessage: `ip is baned`,
		ErrorCode: ipFilter.errorCode,
	}
}
