package external_service

import (
	"encoding/json"
	"github.com/pefish/go-core/api"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
)


// 让外部服务可以通过这个基类调用内部功能
type BaseExternalServiceClass struct {

}


func (this *BaseExternalServiceClass) Init(driver *ExternalServiceDriverClass) {

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
	result := api.ApiResult{}
	_, err := go_http.NewHttpRequester().PostForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if err != nil {
		panic(err)
	}
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
	result := api.ApiResult{}
	_, err := go_http.NewHttpRequester().GetForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if err != nil {
		panic(err)
	}
	if result.Code != 0 {
		go_error.Throw(result.Msg, result.Code)
	}
	return result.Data
}
