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
	go_logger.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateCtxValuesErrorMsg(ctx, `apiMsg`, apiMsg)
	go_logger.Logger.Debug(ctx.Request().Header)
	ctx.Next()
}

func CatchError(ctx iris.Context) {
	if err := recover(); err != nil {
		lang := ctx.GetHeader(`lang`)
		if lang == `` {
			lang = `zh`
		}
		var apiResult api_channel_builder.ApiResult
		if _, ok := err.(go_error.ErrorInfo); !ok {
			errorMessage := ``
			if _, ok := err.(error); !ok {
				errorMessage = err.(string)
			} else {
				errorMessage = err.(error).Error()
			}
			go_logger.Logger.Error(`system_error: ` + errorMessage+"\n"+ctx.Values().Get(`error_msg`).(string)+"\n"+go_stack.Stack.GetStack(go_stack.Option{Skip: 0, Count: 7}))
			ctx.StatusCode(iris.StatusOK)
			if go_application.Application.Debug {
				apiResult = api_channel_builder.ApiResult{
					Msg:  errorMessage,
					Code: 1,
					Data: nil,
				}
			} else {
				apiResult = api_channel_builder.ApiResult{
					Msg:  ``,
					Code: 1,
					Data: nil,
				}
			}
			ctx.JSON(apiResult)
		} else {
			ctx.StatusCode(iris.StatusOK)
			errorInfoStruct := err.(go_error.ErrorInfo)
			errMsg := `error: ` + errorInfoStruct.ErrorMessage
			if errorInfoStruct.Err != nil {
				errMsg += "\nsystem_error: " + errorInfoStruct.Err.Error()
			}
			go_logger.Logger.Error(errMsg +"\n"+ctx.Values().GetString(`error_msg`)+"\n"+go_stack.Stack.GetStack(go_stack.Option{Skip: 0, Count: 7}))
			if go_application.Application.Debug {
				apiResult = api_channel_builder.ApiResult{
					Msg:  errorInfoStruct.ErrorMessage,
					Code: errorInfoStruct.ErrorCode,
					Data: errorInfoStruct.Data,
				}
			} else {
				apiResult = api_channel_builder.ApiResult{
					Msg:  ``,
					Code: errorInfoStruct.ErrorCode,
					Data: errorInfoStruct.Data,
				}
			}
			ctx.JSON(apiResult)
		}
	}
}
