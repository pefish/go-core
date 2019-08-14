package api_strategy

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
	"github.com/pefish/go-jwt"
	"github.com/pefish/go-reflect"
)

type JwtAuthStrategyClass struct {
	errorCode           uint64
	pubKey              string
	headerName          string
	noCheckExpire       bool
	jwtErrorErrorCode   uint64
	jwtExpiredErrorCode uint64
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

func (this *JwtAuthStrategyClass) SetJwtErrorErrorCode(code uint64) {
	this.jwtErrorErrorCode = code
}

func (this *JwtAuthStrategyClass) SetJwtExpiredErrorCode(code uint64) {
	this.jwtExpiredErrorCode = code
}

func (this *JwtAuthStrategyClass) SetNoCheckExpire() {
	this.noCheckExpire = true
}

func (this *JwtAuthStrategyClass) SetPubKey(pubKey string) {
	this.pubKey = pubKey
}

func (this *JwtAuthStrategyClass) SetHeaderName(headerName string) {
	this.headerName = headerName
}

func (this *JwtAuthStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	defer func() {
		if err := recover(); err != nil {
			go_error.Throw(`jwt verify error`, this.errorCode)
		}
	}()
	out.JwtHeaderName = this.headerName
	jwt := out.Ctx.GetHeader(this.headerName)

	verifyResult := false
	if this.noCheckExpire == true {
		verifyResult = go_jwt.Jwt.VerifyJwtSkipClaimsValidation(this.pubKey, jwt)
		if !verifyResult {
			go_error.Throw(`jwt verify error`, this.jwtErrorErrorCode)
		}
	} else {
		verifyResult = go_jwt.Jwt.VerifyJwt(this.pubKey, jwt)
		if !verifyResult {
			go_error.Throw(`jwt verify error or jwt expired`, this.errorCode)
		}
	}
	out.JwtPayload = go_jwt.Jwt.DecodePayloadOfJwtBody(jwt)
	if out.JwtPayload[`user_id`] == nil {
		go_error.Throw(`jwt verify error, user_id not exist`, this.errorCode)
	}

	userId := go_reflect.Reflect.ToUint64(out.JwtPayload[`user_id`])
	out.UserId = userId

	util.UpdateCtxValuesErrorMsg(out.Ctx, `jwtAuth`, userId)
}
