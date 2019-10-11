package api_strategy

import (
	"fmt"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
)

type TestStrategyClass struct {

}

var TestApiStrategy = TestStrategyClass{}

func (this *TestStrategyClass) GetName() string {
	return `test`
}

func (this *TestStrategyClass) GetDescription() string {
	return `对apikey以及签名进行校验`
}

func (this *TestStrategyClass) GetErrorCode() uint64 {
	return 1000
}

type ApikeyAuthParam struct {
	AllowedType string
}

func (this *TestStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	var p ApikeyAuthParam
	p = param.(ApikeyAuthParam)
	fmt.Printf(`%#v`, p)

	go_error.ThrowInternal(`12345test`)
}
