package _type

import (
	"github.com/pefish/go-core/api-strategy/type"
)

type IGlobalStrategy interface {
	Init(param interface{})  // 同步的初始化函数
	_type.IStrategy
}
