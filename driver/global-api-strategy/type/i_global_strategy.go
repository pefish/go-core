package _type

import _type2 "github.com/pefish/go-core/api-strategy/type"

type IGlobalStrategy interface {
	Init(param interface{})  // 同步的初始化函数
	_type2.IStrategy
}
