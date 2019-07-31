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
}

var JwtAuthApiStrategy = JwtAuthStrategyClass{}

type JwtAuthParam struct {
	ErrorCode     uint64
	JwtHeaderName string
	PubKey        string
}

func (this *JwtAuthStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	newParam := param.(JwtAuthParam)
	defer func() {
		if err := recover(); err != nil {
			p_error.Throw(`jwt verify error`, newParam.ErrorCode)
		}
	}()
	jwtHeaderName := newParam.JwtHeaderName
	out.JwtHeaderName = jwtHeaderName
	verifyResult := p_jwt.Jwt.VerifyJwt(newParam.PubKey, ctx.GetHeader(jwtHeaderName))
	if !verifyResult {
		p_error.ThrowInternal(`jwt verify error`)
	}
	out.JwtPayload = p_jwt.Jwt.DecodePayloadOfJwtBody(ctx.GetHeader(jwtHeaderName))
	if out.JwtPayload[`user_id`] == nil {
		p_error.ThrowInternal(`jwt verify error`)
	}

	userId := p_reflect.Reflect.ToUint64(out.JwtPayload[`user_id`])
	out.UserId = &userId

	util.UpdateCtxValuesErrorMsg(ctx, `jwtAuth`, userId)
}
