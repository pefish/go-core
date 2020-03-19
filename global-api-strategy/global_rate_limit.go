// 全局限流器（令牌桶）
package global_api_strategy

import (
	go_application "github.com/pefish/go-application"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-error"
	"time"
)

type GlobalRateLimitStrategyClass struct {
	tokenBucket chan struct{}
	errorCode uint64
}

var GlobalRateLimitStrategy = GlobalRateLimitStrategyClass{
	tokenBucket: make(chan struct{}, 200),
}

func (this *GlobalRateLimitStrategyClass) GetName() string {
	return `GlobalRateLimit`
}

func (this *GlobalRateLimitStrategyClass) GetDescription() string {
	return `global rate limit for all api`
}

func (this *GlobalRateLimitStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *GlobalRateLimitStrategyClass) GetErrorCode() uint64 {
	if this.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return this.errorCode
}


func (this *GlobalRateLimitStrategyClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())

	go func() {
		params := param.(GlobalRateLimitStrategyParam)
		ticker := time.NewTicker(params.FillInterval)
		for {
			select {
			case <-ticker.C:
				select {
				case this.tokenBucket <- struct{}{}:
				default:
				}
			case <- go_application.Application.OnFinished():
				return
			}
		}
	}()
}

type GlobalRateLimitStrategyParam struct {
	FillInterval time.Duration
}

func (this *GlobalRateLimitStrategyClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, this.GetName())

	succ := this.takeAvailable(false)
	if !succ {
		go_error.ThrowInternal(`global rate limit`)
	}
}

func (this *GlobalRateLimitStrategyClass) takeAvailable(block bool) bool{
	var takenResult bool
	if block {
		select {
		case <-this.tokenBucket:
			takenResult = true
		}
	} else {
		select {
		case <-this.tokenBucket:
			takenResult = true
		default:
			takenResult = false
		}
	}
	logger.LoggerDriver.Logger.DebugF("current global rate limit token count: %d", len(this.tokenBucket))
	return takenResult
}
