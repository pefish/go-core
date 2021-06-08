package _type

import (
	_type "github.com/pefish/go-core-type/api-session"
	go_error "github.com/pefish/go-error"
)

type StrategyData struct {
	Strategy IStrategy
	Param    interface{}
	Disable  bool
}

type IStrategy interface {
	Execute(out _type.IApiSession, param interface{}) *go_error.ErrorInfo
	GetName() string
	GetDescription() string
	GetErrorCode() uint64
}
