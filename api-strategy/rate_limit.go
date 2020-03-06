package api_strategy

import (
	"fmt"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-error"
	"time"
)

type RateLimitStrategyClass struct {
	errorCode uint64
	db        *map[string]time.Time // 外部传入的存储api访问频率限制的信息，应当是全局变量
}

var defaultDb = map[string]time.Time{}

var RateLimitApiStrategy = RateLimitStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
	db:        &defaultDb,
}

type RateLimitParam struct {
	Limit time.Duration // 限制多少s只能访问一次
}

func (this *RateLimitStrategyClass) GetName() string {
	return `rateLimit`
}

func (this *RateLimitStrategyClass) GetDescription() string {
	return `rate limit`
}

func (this *RateLimitStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *RateLimitStrategyClass) GetErrorCode() uint64 {
	return this.errorCode
}

func (this *RateLimitStrategyClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, this.GetName())
	if param == nil {
		go_error.Throw(`strategy need param`, this.errorCode)
	}
	newParam := param.(RateLimitParam)
	methodPath := fmt.Sprintf(`%s_%s`, out.GetMethod(), out.GetPath())
	key := fmt.Sprintf(`%s_%s`, out.GetRemoteAddress(), methodPath)
	if !(*this.db)[key].IsZero() && time.Now().Sub((*this.db)[key]) < newParam.Limit {
		go_error.Throw(`api ratelimit`, this.errorCode)
	}

	(*this.db)[key] = time.Now()
}
