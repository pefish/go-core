package api_strategy

import (
	"github.com/kataras/iris"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
	"github.com/pefish/go-jwt"
	"github.com/pefish/go-reflect"
)

type JwtAuthStrategyClass struct {
	errorCode uint64
}

var JwtAuthApiStrategy = JwtAuthStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

type JwtAuthParam struct {
	JwtHeaderName string
	PubKey        string
}

func (this *JwtAuthStrategyClass) GetName() string {
	return `jwtAuth`
}

func (this *JwtAuthStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *JwtAuthStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	newParam := param.(JwtAuthParam)
	defer func() {
		if err := recover(); err != nil {
			go_error.Throw(`jwt verify error`, this.errorCode)
		}
	}()
	jwtHeaderName := newParam.JwtHeaderName
	out.JwtHeaderName = jwtHeaderName
	verifyResult := go_jwt.Jwt.VerifyJwt(newParam.PubKey, ctx.GetHeader(jwtHeaderName))
	if !verifyResult {
		go_error.Throw(`jwt verify error`, this.errorCode)
	}
	out.JwtPayload = go_jwt.Jwt.DecodePayloadOfJwtBody(ctx.GetHeader(jwtHeaderName))
	if out.JwtPayload[`user_id`] == nil {
		go_error.Throw(`jwt verify error`, this.errorCode)
	}

	userId := go_reflect.Reflect.ToUint64(out.JwtPayload[`user_id`])
	out.UserId = &userId

	util.UpdateCtxValuesErrorMsg(ctx, `jwtAuth`, userId)
}
