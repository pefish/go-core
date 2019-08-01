package api_strategy

import (
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
)


type InterfaceStrategy interface {
	Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{})
	GetName() string
}
