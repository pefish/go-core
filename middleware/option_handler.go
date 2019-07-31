package middleware

import (
	"github.com/kataras/iris"
)

func OptionHandle (ctx iris.Context) {
	if ctx.Method() == `OPTIONS` {
		ctx.StatusCode(200)
		return
	}
	ctx.Next()
}
