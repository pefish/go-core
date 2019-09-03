package util

import (
	"github.com/pefish/go-config"
	"github.com/pefish/go-core/service-driver"
)

type DepositAddressServiceClass struct {
	baseUrl string
	apiConfig map[string]interface{}
	driver *service_driver.ServiceDriverClass
}

var DepositAddressService = DepositAddressServiceClass{}

var _ = service_driver.ServiceDriver.Register(&DepositAddressService)

func (this *DepositAddressServiceClass) Init(driver *service_driver.ServiceDriverClass) {
	this.driver = driver
	this.apiConfig = go_config.Config.GetMap(`depositAddressApi`)
	this.baseUrl = this.apiConfig[`baseUrl`].(string)
}

func (this *DepositAddressServiceClass) GetBaseUrl() string {
	return this.baseUrl
}

func (this *DepositAddressServiceClass) ValidateAddress(series string, address string) {
	path := this.apiConfig[`validateAddressPath`].(string)
	this.driver.PostJson(this.GetBaseUrl() + path, map[string]interface{}{
		`series`:  series,
		`address`: address,
	})
}

func (this *DepositAddressServiceClass) GetAddress(series string, address string) {
	path := this.apiConfig[`getAddressPath`].(string)
	this.driver.PostJson(this.GetBaseUrl() + path, map[string]interface{}{
		`series`:  series,
		`address`: address,
	})
}
