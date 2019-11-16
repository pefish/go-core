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
	disableUserId       bool
	errorMsg string
}

var JwtAuthApiStrategy = JwtAuthStrategyClass{
	errorCode: go_error.INTERNAL_ERROR_CODE,
	errorMsg: `Unauthorized`,
}

type JwtAuthParam struct {
}

func (this *JwtAuthStrategyClass) GetName() string {
	return `jwtAuth`
}

func (this *JwtAuthStrategyClass) GetDescription() string {
	return `jwt auth`
}

func (this *JwtAuthStrategyClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *JwtAuthStrategyClass) SetErrorMessage(msg string) {
	this.errorMsg = msg
}

func (this *JwtAuthStrategyClass) GetErrorCode() uint64 {
	return this.errorCode
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

func (this *JwtAuthStrategyClass) DisableUserId() {
	this.disableUserId = true
}

func (this *JwtAuthStrategyClass) SetPubKey(pubKey string) {
	this.pubKey = pubKey
}

func (this *JwtAuthStrategyClass) SetHeaderName(headerName string) {
	this.headerName = headerName
}

func (this *JwtAuthStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	out.JwtHeaderName = this.headerName
	jwt := out.Ctx.GetHeader(this.headerName)

	if this.noCheckExpire == true {
		verifyResult, err := go_jwt.Jwt.VerifyJwtSkipClaimsValidation(this.pubKey, jwt)
		if err != nil {
			go_error.ThrowWithInternalMsg(this.errorMsg, err.Error(), this.jwtErrorErrorCode)
		}
		if !verifyResult {
			go_error.ThrowWithInternalMsg(this.errorMsg, `jwt verify error`, this.jwtErrorErrorCode)
		}
	} else {
		verifyResult, err := go_jwt.Jwt.VerifyJwt(this.pubKey, jwt)
		if err != nil {
			go_error.ThrowWithInternalMsg(this.errorMsg, err.Error(), this.jwtErrorErrorCode)
		}
		if !verifyResult {
			go_error.ThrowWithInternalMsg(this.errorMsg,`jwt verify error or jwt expired`, this.errorCode)
		}
	}
	jwtBody, err := go_jwt.Jwt.DecodeBodyOfJwt(jwt)
	if err != nil {
		go_error.ThrowWithInternalMsg(this.errorMsg, err.Error(), this.jwtErrorErrorCode)
	}
	out.JwtBody = jwtBody
	if !this.disableUserId {
		jwtPayload := out.JwtBody[`payload`].(map[string]interface{})
		if jwtPayload[`user_id`] == nil {
			go_error.ThrowWithInternalMsg(this.errorMsg,`jwt verify error, user_id not exist`, this.errorCode)
		}

		userId := go_reflect.Reflect.MustToUint64(jwtPayload[`user_id`])
		out.UserId = userId

		util.UpdateCtxValuesErrorMsg(out.Ctx, `jwtAuth`, userId)
	}
}
