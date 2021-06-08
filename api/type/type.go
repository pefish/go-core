package _type

type IApi interface {
	GetDescription() string
	GetParamType() string
	GetParams() interface{}
}
