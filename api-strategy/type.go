package api_strategy

import (
	api_session "github.com/pefish/go-core/api-session"
	go_error "github.com/pefish/go-error"
)

type StrategyData struct {
	Strategy IStrategy
	Param    interface{}
	Disable  bool
}

type IStrategy interface {
	Execute(out api_session.IApiSession, param interface{}) *go_error.ErrorInfo
	GetName() string
	GetDescription() string
	GetErrorCode() uint64
}
