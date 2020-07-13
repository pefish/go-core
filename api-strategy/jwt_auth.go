package api_strategy

import (
	jwt2 "github.com/dgrijalva/jwt-go"
	_type "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-core/driver/logger"
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
	disableUserId       bool
	errorMsg            string
}

var JwtAuthApiStrategy = JwtAuthStrategyClass{
	errorCode:           go_error.INTERNAL_ERROR_CODE,
	errorMsg:            `Unauthorized`,
}

type JwtAuthParam struct {
}

func (jwtAuth *JwtAuthStrategyClass) GetName() string {
	return `jwtAuth`
}

func (jwtAuth *JwtAuthStrategyClass) GetDescription() string {
	return `jwt auth`
}

func (jwtAuth *JwtAuthStrategyClass) SetErrorCode(code uint64) {
	jwtAuth.errorCode = code
}

func (jwtAuth *JwtAuthStrategyClass) SetErrorMessage(msg string) {
	jwtAuth.errorMsg = msg
}

func (jwtAuth *JwtAuthStrategyClass) GetErrorCode() uint64 {
	return jwtAuth.errorCode
}

func (jwtAuth *JwtAuthStrategyClass) SetNoCheckExpire() {
	jwtAuth.noCheckExpire = true
}

func (jwtAuth *JwtAuthStrategyClass) DisableUserId() {
	jwtAuth.disableUserId = true
}

func (jwtAuth *JwtAuthStrategyClass) SetPubKey(pubKey string) {
	jwtAuth.pubKey = pubKey
}

func (jwtAuth *JwtAuthStrategyClass) SetHeaderName(headerName string) {
	jwtAuth.headerName = headerName
}

func (jwtAuth *JwtAuthStrategyClass) Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, jwtAuth.GetName())

	out.SetJwtHeaderName(jwtAuth.headerName)
	jwt := out.Header(jwtAuth.headerName)

	verifyResult, token, err := go_jwt.Jwt.VerifyJwt(jwtAuth.pubKey, jwt, jwtAuth.noCheckExpire)
	if err != nil {
		return &go_error.ErrorInfo{
			InternalErrorMessage: jwtAuth.errorMsg,
			ErrorMessage: jwtAuth.errorMsg,
			ErrorCode: jwtAuth.errorCode,
			Err: err,
		}
	}
	if !verifyResult {
		return &go_error.ErrorInfo{
			InternalErrorMessage: `jwt verify error or jwt expired`,
			ErrorMessage: jwtAuth.errorMsg,
			ErrorCode: jwtAuth.errorCode,
		}
	}
	jwtBody := token.Claims.(jwt2.MapClaims)
	out.SetJwtBody(jwtBody)
	if !jwtAuth.disableUserId {
		jwtPayload := jwtBody[`payload`].(map[string]interface{})
		if jwtPayload[`user_id`] == nil {
			return &go_error.ErrorInfo{
				InternalErrorMessage: `jwt verify error, user_id not exist`,
				ErrorMessage: jwtAuth.errorMsg,
				ErrorCode: jwtAuth.errorCode,
			}
		}

		userId := go_reflect.Reflect.MustToUint64(jwtPayload[`user_id`])
		out.SetUserId(userId)

		util.UpdateSessionErrorMsg(out, `jwtAuth`, userId)
	}
	return nil
}
