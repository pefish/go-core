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
	errorCode       uint64
	pubKey          string
	headerName      string
	noExpireForever bool
}

var JwtAuthApiStrategy = JwtAuthStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
}

type JwtAuthParam struct {
}

func (this *JwtAuthStrategyClass) GetName() string {
	return `jwtAuth`
}

func (this *JwtAuthStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *JwtAuthStrategyClass) SetNoExpireForever() {
	this.noExpireForever = true
}

func (this *JwtAuthStrategyClass) SetPubKey(pubKey string) {
	this.pubKey = pubKey
}

func (this *JwtAuthStrategyClass) SetHeaderName(headerName string) {
	this.headerName = headerName
}

func (this *JwtAuthStrategyClass) Execute(ctx iris.Context, out *api_session.ApiSessionClass, param interface{}) {
	defer func() {
		if err := recover(); err != nil {
			go_error.Throw(`jwt verify error`, this.errorCode)
		}
	}()
	out.JwtHeaderName = this.headerName
	jwt := ctx.GetHeader(this.headerName)

	verifyResult := false
	if this.noExpireForever == true {
		verifyResult = go_jwt.Jwt.VerifyJwtSkipClaimsValidation(this.pubKey, jwt)
	} else {
		verifyResult = go_jwt.Jwt.VerifyJwt(this.pubKey, jwt)
	}
	if !verifyResult {
		go_error.Throw(`jwt verify error`, this.errorCode)
	}
	out.JwtPayload = go_jwt.Jwt.DecodePayloadOfJwtBody(jwt)
	if out.JwtPayload[`user_id`] == nil {
		go_error.Throw(`jwt verify error`, this.errorCode)
	}

	userId := go_reflect.Reflect.ToUint64(out.JwtPayload[`user_id`])
	out.UserId = userId

	util.UpdateCtxValuesErrorMsg(ctx, `jwtAuth`, userId)
}
