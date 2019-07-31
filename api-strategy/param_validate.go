package api_strategy

import (
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/validator"
)

type ParamValidateStrategyClass struct {

}

var ParamValidateApiStrategy = ParamValidateStrategyClass{}


func (this *ParamValidateStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	myValidator := validator.ValidatorClass{}
	myValidator.Init()
	out.Validator = myValidator.Validator
}
