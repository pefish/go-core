package util

import (
	"fmt"
	api_session "github.com/pefish/go-core/api-session"
)

func UpdateSessionErrorMsg(apiSession api_session.InterfaceApiSession, key string, data interface{}) {
	errorMsg := apiSession.Date(`error_msg`)
	if errorMsg == nil {
		apiSession.SetData(`error_msg`, fmt.Sprintf("%s: %v\n", key, data))
	} else {
		apiSession.SetData(`error_msg`, fmt.Sprintf("%s%s: %v\n", errorMsg.(string), key, data))
	}
}
