package main

import (
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-core/swagger"
	"test/route"
)

func main() {
	service.Service.SetName(`test`)
	service.Service.SetPath(`/api/test`)
	service.Service.SetRoutes(route.TestRoute)
	swagger.GetSwaggerInstance().GeneSwagger(`www.zexchange.xyz`, `swagger.json`, `json`)
}
