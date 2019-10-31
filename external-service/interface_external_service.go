package external_service

type ExternalServiceInterface interface {
	Init(driver *ServiceDriverClass)

	PostJsonForStruct(url string, params map[string]interface{}, struct_ interface{})
	PostJson(url string, params map[string]interface{}) interface{}
	GetJsonForStruct(url string, params map[string]interface{}, struct_ interface{})
	GetJson(url string, params map[string]interface{}) interface{}
}
