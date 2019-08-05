package api_session

import (
	"encoding/json"
	"github.com/kataras/iris"
)

type ApiHandlerType func(apiContext *ApiSessionClass) interface{}

type ApiSessionClass struct {
	Ctx iris.Context

	JwtHeaderName string
	JwtPayload    map[string]interface{}
	UserId        uint64

	RouteName string

	Lang       string
	ClientType string // web、android、ios

	Options map[string]interface{}

	Params interface{}
}

func NewApiSession() *ApiSessionClass {
	return &ApiSessionClass{
		Options: map[string]interface{}{},
	}
}

func (this *ApiSessionClass) ScanParams(dest interface{}) {
	result, err := json.Marshal(this.Params)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(result, &dest); err != nil {
		panic(err)
	}
}
