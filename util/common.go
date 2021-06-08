package util

import (
	"fmt"
	_type "github.com/pefish/go-core-type/api-session"
)

func UpdateSessionErrorMsg(apiSession _type.IApiSession, key string, data interface{}) {
	errorMsg := apiSession.Data(`error_msg`)
	if errorMsg == nil {
		apiSession.SetData(`error_msg`, fmt.Sprintf("%s: %v\n", key, data))
	} else {
		apiSession.SetData(`error_msg`, fmt.Sprintf("%s%s: %v\n", errorMsg.(string), key, data))
	}
}
