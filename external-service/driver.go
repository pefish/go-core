package external_service

type ServiceDriverClass struct {
	externalServices map[string]ExternalServiceInterface
}

var ServiceDriver = ServiceDriverClass{
	externalServices: map[string]ExternalServiceInterface{},
}

func (this *ServiceDriverClass) Startup() {
	for _, v := range this.externalServices {
		v.Init(this)
	}
}

func (this *ServiceDriverClass) Register(name string, svc ExternalServiceInterface) bool {
	this.externalServices[name] = svc
	return true
}

func (this *ServiceDriverClass) Call(name string, method string) interface{} {
	// TODO 调用name外部服务
	return nil
}