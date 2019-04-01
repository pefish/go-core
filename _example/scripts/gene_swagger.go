package main

import (
	"github.com/pefish/go-core/swagger"
	"test/src/service"
)

func main() {
	service.TestService.Init()
	swagger.GetSwaggerInstance().SetService(&service.TestService).GeneSwagger(`https://www.zexchange.xyz`, `swagger.json`, `json`)
}
