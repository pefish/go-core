package external_service

import "github.com/pefish/go-core/external-service"

func init() {
	external_service.ServiceDriver.Register(`deposit_address`, &DepositAddressService)
}

type DepositAddressServiceClass struct {
	baseUrl string
	external_service.BaseExternalServiceClass
}

var DepositAddressService = DepositAddressServiceClass{}


func (this *DepositAddressServiceClass) Init(driver *external_service.ServiceDriverClass) {
	this.baseUrl = `http://baidu.com`
}

func (this *DepositAddressServiceClass) Test(series string, address string) interface{} {
	path := ``
	return this.PostJson(this.baseUrl + path, map[string]interface{}{
		`series`:  series,
		`address`: address,
	})
}
