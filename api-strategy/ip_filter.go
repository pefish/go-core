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
	GetValidIp func(apiSession *api_session.ApiSessionClass) []string
}

func (this *IpFilterStrategyClass) GetName() string {
	return `ipFilter`
}

func (this *IpFilterStrategyClass) GetDescription() string {
	return `filter ip`
}

func (this *IpFilterStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *IpFilterStrategyClass) GetErrorCode() uint64 {
	return this.errorCode
}

func (this *IpFilterStrategyClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, this.GetName())
	if param == nil {
		go_error.Throw(`strategy need param`, this.errorCode)
	}
	newParam := param.(IpFilterParam)
	if newParam.GetValidIp == nil {
		return
	}
	clientIp := out.GetRemoteAddress()
	allowedIps := newParam.GetValidIp(out)
	for _, ip := range allowedIps {
		if ip == clientIp {
			return
		}
	}
	go_error.ThrowInternal(`ip is baned`)
}
