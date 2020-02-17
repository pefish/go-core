package _interface

// 实现访问Api的目的。通过接口访问，达到避免循环引用的目的
type InterfaceApi interface {
	GetDescription() string
	GetParamType() string
	GetParams() interface{}
}
