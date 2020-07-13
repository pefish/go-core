package _interface

type IApi interface {
	GetDescription() string
	GetParamType() string
	GetParams() interface{}
}
