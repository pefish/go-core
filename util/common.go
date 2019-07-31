package util

import (
	"fmt"
	"github.com/kataras/iris"
)

func UpdateCtxValuesErrorMsg(ctx iris.Context, key string, data interface{}) {
	errorMsg := ctx.Values().GetString(`error_msg`)
	if errorMsg == `` {
		ctx.Values().Set(`error_msg`, fmt.Sprintf("%s: %v\n", key, data))
	} else {
		ctx.Values().Set(`error_msg`, fmt.Sprintf("%s%s: %v\n", errorMsg, key, data))
	}
}
