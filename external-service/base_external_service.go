package external_service

import (
	"encoding/json"
	_interface "github.com/pefish/go-core/interface"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
)

type BaseExternalServiceClass struct {

}


func (this *BaseExternalServiceClass) Init(driver *ServiceDriverClass) {

}

func (this *BaseExternalServiceClass) PostJsonForStruct(url string, params map[string]interface{}, struct_ interface{}) {
	data := this.PostJson(url, params)
	inrec, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(inrec, struct_)
	if err != nil {
		panic(err)
	}
}

func (this *BaseExternalServiceClass) PostJson(url string, params map[string]interface{}) interface{} {
	result := _interface.ApiResult{}
	go_http.Http.PostJsonForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if result.Code != 0 {
		go_error.Throw(result.Msg, result.Code)
	}
	return result.Data
}

func (this *BaseExternalServiceClass) GetJsonForStruct(url string, params map[string]interface{}, struct_ interface{}) {
	data := this.GetJson(url, params)
	inrec, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(inrec, struct_)
	if err != nil {
		panic(err)
	}
}

func (this *BaseExternalServiceClass) GetJson(url string, params map[string]interface{}) interface{} {
	result := _interface.ApiResult{}
	go_http.Http.GetForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if result.Code != 0 {
		go_error.Throw(result.Msg, result.Code)
	}
	return result.Data
}
