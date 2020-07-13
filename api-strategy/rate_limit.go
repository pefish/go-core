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

func (rateLimit *RateLimitStrategyClass) GetName() string {
	return `rateLimit`
}

func (rateLimit *RateLimitStrategyClass) GetDescription() string {
	return `rate limit`
}

func (rateLimit *RateLimitStrategyClass) SetErrorCode(code uint64) {
	rateLimit.errorCode = code
}

func (rateLimit *RateLimitStrategyClass) GetErrorCode() uint64 {
	return rateLimit.errorCode
}

func (rateLimit *RateLimitStrategyClass) Execute(out api_session.InterfaceApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, rateLimit.GetName())
	if param == nil {
		return &go_error.ErrorInfo{
			InternalErrorMessage: `strategy need param`,
			ErrorMessage: `strategy need param`,
			ErrorCode: rateLimit.errorCode,
		}
	}
	newParam := param.(RateLimitParam)
	methodPath := fmt.Sprintf(`%s_%s`, out.Method(), out.Path())
	key := fmt.Sprintf(`%s_%s`, out.RemoteAddress(), methodPath)
	if !(*rateLimit.db)[key].IsZero() && time.Now().Sub((*rateLimit.db)[key]) < newParam.Limit {
		return &go_error.ErrorInfo{
			InternalErrorMessage: `api ratelimit`,
			ErrorMessage: `api ratelimit`,
			ErrorCode: rateLimit.errorCode,
		}
	}

	(*rateLimit.db)[key] = time.Now()
	return nil
}
