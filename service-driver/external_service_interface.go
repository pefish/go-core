package service_driver

type ExternalServiceInterface interface {
	GetBaseUrl() string
	Init(driver *ServiceDriverClass)
}
