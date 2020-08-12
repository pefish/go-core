package api_strategy

import (
	"errors"
	_type "github.com/pefish/go-core/api-session/type"
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
	GetValidIp func(apiSession _type.IApiSession) []string
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

func (ipFilter *IpFilterStrategyClass) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, ipFilter.GetName())
	if param == nil {
		return go_error.WrapWithAll(errors.New(`strategy need param`), ipFilter.errorCode, nil)
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
	return go_error.WrapWithAll(errors.New(`ip is baned`), ipFilter.errorCode, nil)
}
