package util

import (
	"fmt"
	api_session "github.com/pefish/go-core/api-session"
)

func UpdateSessionErrorMsg(apiSession *api_session.ApiSessionClass, key string, data interface{}) {
	errorMsg := apiSession.Datas[`error_msg`]
	if errorMsg == nil {
		apiSession.Datas[`error_msg`] = fmt.Sprintf("%s: %v\n", key, data)
	} else {
		apiSession.Datas[`error_msg`] = fmt.Sprintf("%s%s: %v\n", errorMsg.(string), key, data)
	}
}
