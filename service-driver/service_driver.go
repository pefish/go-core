package service_driver

import (
	"encoding/json"
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
)

type ServiceDriverClass struct {
	ExternalServices []ExternalServiceInterface
}

var ServiceDriver = ServiceDriverClass{}

func (this *ServiceDriverClass) Init() {
	for _, v := range this.ExternalServices {
		v.Init(this)
	}
}

func (this *ServiceDriverClass) Register(svc ExternalServiceInterface) bool {
	this.ExternalServices = append(this.ExternalServices, svc)
	return true
}

func (this *ServiceDriverClass) PostJsonForStruct(url string, params map[string]interface{}, struct_ interface{}) {
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

func (this *ServiceDriverClass) PostJson(url string, params map[string]interface{}) interface{} {
	result := api_channel_builder.ApiResult{}
	go_http.Http.PostJsonForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if result.Code != 0 {
		go_error.Throw(result.Msg, result.Code)
	}
	return result.Data
}

func (this *ServiceDriverClass) GetJsonForStruct(url string, params map[string]interface{}, struct_ interface{}) {
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

func (this *ServiceDriverClass) GetJson(url string, params map[string]interface{}) interface{} {
	result := api_channel_builder.ApiResult{}
	go_http.Http.GetForStruct(go_http.RequestParam{
		Url:    url,
		Params: params,
	}, &result)
	if result.Code != 0 {
		go_error.Throw(result.Msg, result.Code)
	}
	return result.Data
}
