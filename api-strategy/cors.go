package api_strategy

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"strings"
)

type CorsStrategyClass struct {
	AllowOrigins []string
}

var CorsApiStrategy = CorsStrategyClass{
	AllowOrigins: []string{`*`},
}


func (this *CorsStrategyClass) SetAllowedOrigins(origins []string) {
	this.AllowOrigins = origins
}

func (this *CorsStrategyClass) GetName() string {
	return `cors`
}

func (this *CorsStrategyClass) GetDescription() string {
	return `cors`
}

func (this *CorsStrategyClass) GetErrorCode() uint64 {
	return go_error.INTERNAL_ERROR_CODE
}

func (this *CorsStrategyClass) isOriginAllowed(origin string) bool {
	origin = strings.ToLower(origin)
	for _, o := range this.AllowOrigins {
		if o == origin || o == `*` {
			return true
		}

		if i := strings.IndexByte(o, '*'); i >= 0 {
			w := wildcard{o[0:i], o[i+1:]}
			if w.match(origin) {
				return true
			}
		}
	}
	return false
}

func (this *CorsStrategyClass) Execute(route *api_channel_builder.Route, out *api_session.ApiSessionClass, param interface{}) {
	origin := out.Ctx.GetHeader("Origin")

	out.Ctx.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

	if origin != "" && !this.isOriginAllowed(origin) {
		go_error.ThrowInternal(`origin is not be allowed`)
	}

	out.Ctx.Header("Access-Control-Allow-Origin", origin)
	out.Ctx.Header("Access-Control-Allow-Methods", out.Ctx.Method())
	out.Ctx.Header("Access-Control-Allow-Headers", "*")
	out.Ctx.Header("Access-Control-Allow-Credentials", "true")
}



type wildcard struct {
	prefix string
	suffix string
}

func (w wildcard) match(s string) bool {
	return len(s) >= len(w.prefix+w.suffix) && strings.HasPrefix(s, w.prefix) && strings.HasSuffix(s, w.suffix)
}
