package middlewares

import (
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/util"
	"github.com/pefish/go-logger"
	"github.com/kataras/iris"
)

func ErrorHandle(ctx iris.Context) {
	defer api_channel_builder.ApiChannelBuilder.CatchError(ctx)
	apiMsg := fmt.Sprintf(`%s %s %s`, ctx.RemoteAddr(), ctx.Path(), ctx.Method())
	p_logger.Logger.Info(fmt.Sprintf(`---------------- %s ----------------`, apiMsg))
	util.UpdateCtxValuesErrorMsg(ctx, `apiMsg`, apiMsg)
	if p_application.Application.Debug {
		p_logger.Logger.Info(ctx.Request().Header)
	}
	ctx.Next()
}
