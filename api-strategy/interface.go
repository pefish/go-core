package api_strategy

import api_session "github.com/pefish/go-core/api-session"

type InterfaceStrategy interface {
	Execute(out *api_session.ApiSessionClass, param interface{})
	GetName() string
	GetDescription() string
	GetErrorCode() uint64
}

