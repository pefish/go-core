package external_service

import (
	"encoding/json"
	"errors"
	"github.com/pefish/go-core/api"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
)


// 让外部服务可以通过这个基类调用内部功能
type BaseExternalServiceClass struct {

}


func (bes *BaseExternalServiceClass) Init(driver *ExternalServiceDriverClass) {

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
	_, err := go_http.NewHttpRequester().PostForStruct(go_http.RequestParam{
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
	_, err := go_http.NewHttpRequester().GetForStruct(go_http.RequestParam{
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
