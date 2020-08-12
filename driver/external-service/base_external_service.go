package external_service

import (
	"encoding/json"
	"errors"
	go_application "github.com/pefish/go-application"
	"github.com/pefish/go-core/api"
	"github.com/pefish/go-core/driver/logger"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
	"time"
)


// 让外部服务可以通过这个基类调用内部功能
type BaseExternalServiceClass struct {
	timeout time.Duration
}


func (bes *BaseExternalServiceClass) Init(driver *ExternalServiceDriverClass) {
	bes.timeout = 20 * time.Second
}

func (bes *BaseExternalServiceClass) SetTimeout(timeout time.Duration) {
	bes.timeout = timeout
}

func (bes *BaseExternalServiceClass) PostJsonForStruct(url string, params map[string]interface{}, struct_ interface{}) *go_error.ErrorInfo {
	data, errInfo := bes.PostJson(url, params)
	if errInfo != nil {
		return errInfo
	}
	inrec, err := json.Marshal(data)
	if err != nil {
		return go_error.Wrap(err)
	}
	err = json.Unmarshal(inrec, struct_)
	if err != nil {
		return go_error.Wrap(err)
	}
	return nil
}

func (bes *BaseExternalServiceClass) PostJson(url string, params map[string]interface{}) (interface{}, *go_error.ErrorInfo) {
	result := api.ApiResult{}
	_, err := go_http.NewHttpRequester(go_http.WithIsDebug(go_application.Application.Debug), go_http.WithLogger(logger.LoggerDriver.Logger), go_http.WithTimeout(bes.timeout)).PostForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if err != nil {
		return nil, go_error.Wrap(err)
	}
	if result.Code != 0 {
		return nil, go_error.WrapWithAll(errors.New(result.Msg), result.Code, result.Data)
	}
	return result.Data, nil
}

func (bes *BaseExternalServiceClass) GetJsonForStruct(url string, params map[string]interface{}, struct_ interface{}) *go_error.ErrorInfo {
	data, errInfo := bes.GetJson(url, params)
	if errInfo != nil {
		return errInfo
	}
	inrec, err := json.Marshal(data)
	if err != nil {
		return go_error.Wrap(err)
	}
	err = json.Unmarshal(inrec, struct_)
	if err != nil {
		return go_error.Wrap(err)
	}
	return nil
}

func (bes *BaseExternalServiceClass) GetJson(url string, params map[string]interface{}) (interface{}, *go_error.ErrorInfo) {
	result := api.ApiResult{}
	_, err := go_http.NewHttpRequester(go_http.WithIsDebug(go_application.Application.Debug), go_http.WithLogger(logger.LoggerDriver.Logger), go_http.WithTimeout(bes.timeout)).GetForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if err != nil {
		return nil, go_error.Wrap(err)
	}
	if result.Code != 0 {
		return nil, go_error.WrapWithAll(errors.New(result.Msg), result.Code, result.Data)
	}
	return result.Data, nil
}
