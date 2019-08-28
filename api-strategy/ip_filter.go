package api_strategy

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
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
	return `ip_filter`
}

func (this *IpFilterStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *IpFilterStrategyClass) GetErrorCode() uint64 {
	return this.errorCode
}

func (this *IpFilterStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	newParam := param.(IpFilterParam)
	if newParam.GetValidIp == nil {
		return
	}
	clientIp := out.Ctx.RemoteAddr()
	allowedIps := newParam.GetValidIp(out)
	for _, ip := range allowedIps {
		if ip == clientIp {
			return
		}
	}
	go_error.ThrowInternal(`ip is baned`)
}
