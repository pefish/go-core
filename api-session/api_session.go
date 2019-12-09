package api_session

import (
	"github.com/kataras/iris"
	"github.com/mitchellh/mapstructure"
)

type ApiHandlerType func(apiContext *ApiSessionClass) interface{}

type ApiSessionClass struct {
	Ctx iris.Context

	JwtHeaderName string
	JwtBody       map[string]interface{}
	UserId        uint64

	RouteName string

	Lang       string
	ClientType string // web、android、ios

	Datas map[string]interface{}

	OriginalParams map[string]interface{} // 客户端传过来的原始参数
	Params         map[string]interface{} // 经过前置处理器修饰过的参数
}

func NewApiSession() *ApiSessionClass {
	return &ApiSessionClass{
		Datas: map[string]interface{}{},
	}
}

func (this *ApiSessionClass) ScanParams(dest interface{}) {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &dest,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(this.Params)
	if err != nil {
		panic(err)
	}
}
