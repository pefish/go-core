package middleware

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-error"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-stack"
)

func ErrorHandle(ctx iris.Context) {
	defer CatchError(ctx)
	apiMsg := fmt.Sprintf(`%s %s %s`, ctx.RemoteAddr(), ctx.Path(), ctx.Method())
	p_logger.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateCtxValuesErrorMsg(ctx, `apiMsg`, apiMsg)
	if p_application.Application.Debug {
		p_logger.Logger.Info(ctx.Request().Header)
	}
	ctx.Next()
}

func CatchError(ctx iris.Context) {
	if err := recover(); err != nil {
		lang := ctx.GetHeader(`lang`)
		if lang == `` {
			lang = `zh`
		}
		var apiResult api_channel_builder.ApiResult
		if _, ok := err.(p_error.ErrorInfo); !ok {
			errorMessage := ``
			if _, ok := err.(error); !ok {
				errorMessage = err.(string)
			} else {
				errorMessage = err.(error).Error()
			}
			p_logger.Logger.Error(`system_error: ` + errorMessage+"\n"+ctx.Values().Get(`error_msg`).(string)+"\n"+go_stack.Stack.GetStack(go_stack.Option{Skip: 2, Count: 7}))
			ctx.StatusCode(iris.StatusOK)
			if p_application.Application.Debug {
				apiResult = api_channel_builder.ApiResult{
					ErrorMessage: &errorMessage,
					ErrorCode:    1,
					Data:         nil,
				}
			} else {
				apiResult = api_channel_builder.ApiResult{
					ErrorMessage: nil,
					ErrorCode:    1,
					Data:         nil,
				}
			}
			ctx.JSON(apiResult)
		} else {
			ctx.StatusCode(iris.StatusOK)
			errorInfoStruct := err.(p_error.ErrorInfo)
			errMsg := `error: ` + errorInfoStruct.ErrorMessage
			if errorInfoStruct.Err != nil {
				errMsg += "\nsystem_error: " + errorInfoStruct.Err.Error()
			}
			p_logger.Logger.Error(errMsg +"\n"+ctx.Values().GetString(`error_msg`)+"\n"+go_stack.Stack.GetStack(go_stack.Option{Skip: 2, Count: 7}))
			if p_application.Application.Debug {
				apiResult = api_channel_builder.ApiResult{
					ErrorMessage: &errorInfoStruct.ErrorMessage,
					ErrorCode:    errorInfoStruct.ErrorCode,
					Data:         errorInfoStruct.Data,
				}
			} else {
				apiResult = api_channel_builder.ApiResult{
					ErrorMessage: nil,
					ErrorCode:    errorInfoStruct.ErrorCode,
					Data:         errorInfoStruct.Data,
				}
			}
			ctx.JSON(apiResult)
		}
	}
}
