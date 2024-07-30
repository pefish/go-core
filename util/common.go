package util

import (
	"fmt"

	i_core "github.com/pefish/go-interface/i-core"
)

func UpdateSessionErrorMsg(apiSession i_core.IApiSession, key string, data interface{}) {
	errorMsg := apiSession.Data(`error_msg`)
	if errorMsg == nil {
		apiSession.SetData(`error_msg`, fmt.Sprintf("%s: %v\n", key, data))
	} else {
		apiSession.SetData(`error_msg`, fmt.Sprintf("%s%s: %v\n", errorMsg.(string), key, data))
	}
}
