package main

import (
	"github.com/pefish/go-core/swagger"
	"test/service"
)

func main() {
	route.TestService.Init()
	swagger.GetSwaggerInstance().SetService(&route.TestService).GeneSwagger(`www.zexchange.xyz`, `swagger.json`, `json`)
}
