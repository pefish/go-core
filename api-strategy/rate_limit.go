package api_strategy

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"time"
)

type RateLimitStrategyClass struct {
	errorCode uint64
}

var RateLimitApiStrategy = RateLimitStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

type RateLimitParam struct {
	Db    *map[string]time.Time // 外部传入的存储api访问频率限制的信息，应当是全局变量
	Limit time.Duration         // 限制多少s只能访问一次
}

func (this *RateLimitStrategyClass) GetName() string {
	return `rateLimit`
}

func (this *RateLimitStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *RateLimitStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	newParam := param.(RateLimitParam)
	methodPath := fmt.Sprintf(`%s_%s`, ctx.Method(), ctx.Path())
	key := fmt.Sprintf(`%s_%s`, ctx.RemoteAddr(), methodPath)
	if newParam.Db == nil {
		newParam.Db = &map[string]time.Time{}
	}
	if !(*newParam.Db)[key].IsZero() && time.Now().Sub((*newParam.Db)[key]) < newParam.Limit {
		go_error.Throw(`api ratelimit`, this.errorCode)
	}

	(*newParam.Db)[key] = time.Now()
}
