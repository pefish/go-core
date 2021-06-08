package external_service

import (
	_type "github.com/pefish/go-core/driver/external-service/type"
	"sync"
)

// 接口驱动
type ExternalServiceDriver struct {
	externalServices map[string]_type.IExternalService
	sync.Once
}

var ExternalServiceDriverInstance = ExternalServiceDriver{
	externalServices: map[string]_type.IExternalService{},
}

func (esd *ExternalServiceDriver) Startup() {
	esd.Do(func() {
		for _, v := range esd.externalServices {
			v.Init()
		}
	})
}

func (esd *ExternalServiceDriver) Register(name string, svc _type.IExternalService) bool {
	esd.externalServices[name] = svc
	return true
}

func (esd *ExternalServiceDriver) ExternalServices() map[string]_type.IExternalService {
	return esd.externalServices
}
