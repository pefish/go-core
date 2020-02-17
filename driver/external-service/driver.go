package external_service

type ExternalServiceDriverClass struct {
	externalServices map[string]ExternalServiceInterface
}

var ExternalServiceDriver = ExternalServiceDriverClass{
	externalServices: map[string]ExternalServiceInterface{},
}

func (this *ExternalServiceDriverClass) Startup() {
	for _, v := range this.externalServices {
		v.Init(this)
	}
}

func (this *ExternalServiceDriverClass) Register(name string, svc ExternalServiceInterface) bool {
	this.externalServices[name] = svc
	return true
}

func (this *ExternalServiceDriverClass) Call(name string, method string) interface{} {
	// TODO 调用name外部服务
	return nil
}