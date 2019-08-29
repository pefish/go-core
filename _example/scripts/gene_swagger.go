package main

import (
	"github.com/pefish/go-core/swagger"
	"test/service"
)

func main() {
	service.TestService.Init()
	swagger.GetSwaggerInstance().SetService(&service.TestService).GeneSwagger(`www.zexchange.xyz`, `swagger.json`, `json`)
}
