package external_service

// 接口驱动
type ExternalServiceDriverClass struct {
	externalServices map[string]IExternalService
}

var ExternalServiceDriver = ExternalServiceDriverClass{
	externalServices: map[string]IExternalService{},
}

func (this *ExternalServiceDriverClass) Startup() {
	for _, v := range this.externalServices {
		v.Init(this)
	}
}

func (this *ExternalServiceDriverClass) Register(name string, svc IExternalService) bool {
	this.externalServices[name] = svc
	return true
}

func (this *ExternalServiceDriverClass) Call(name string, method string) interface{} {
	// TODO 调用name外部服务
	return nil
}

type IExternalService interface {
	Init(driver *ExternalServiceDriverClass)

	PostJsonForStruct(url string, params map[string]interface{}, struct_ interface{})
	PostJson(url string, params map[string]interface{}) interface{}
	GetJsonForStruct(url string, params map[string]interface{}, struct_ interface{})
	GetJson(url string, params map[string]interface{}) interface{}
}
