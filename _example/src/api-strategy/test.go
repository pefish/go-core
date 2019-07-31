package api_strategy

import (
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
)

type TestStrategyClass struct {

}

var TestApiStrategy = TestStrategyClass{}

func (this *TestStrategyClass) GetName() string {
	return `test`
}

func (this *TestStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	p_error.ThrowInternal(`12345test`)
}
