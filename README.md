# Go-Core Web Framework

<a>![](https://img.shields.io/badge/Go%20Coverage-72%25-brightgreen.svg?longCache=true&style=flat)</a>

Go-Core is a fast, simple and young web framework for Go.

It provides some features. like
1. Api Strategy
2. Swagger Generator
3. Return Hook
4. Cors
5. Api Error Catch
6. ...

## Quick start

```go
package main

import (
	"github.com/pefish/go-core/api"
	_type2 "github.com/pefish/go-core/api-session/type"
	"github.com/pefish/go-core/api-strategy/type"
	global_api_strategy2 "github.com/pefish/go-core/driver/global-api-strategy"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	api_strategy "github.com/pefish/go-core/pkg/api-strategy"
	global_api_strategy3 "github.com/pefish/go-core/pkg/global-api-strategy"
	"github.com/pefish/go-core/service"
	"github.com/pefish/go-error"
	"log"
	"time"
)

func main() {
	service.Service.SetName(`test service`) // set service name
	service.Service.SetPath(`/api/test`)
	global_api_strategy.ParamValidateStrategyInstance.SetErrorCode(2005)
	service.Service.SetRoutes([]*api.Api{
		{
			Description: "this is a test api",
			Path:        "/v1/test_api",
			Method:      `POST`,
			Strategies: []_type.StrategyData{
				{
					Strategy: &api_strategy.IpFilterStrategy,
					Param: api_strategy.IpFilterParam{
						GetValidIp: func(apiSession _type2.IApiSession) []string {
							return []string{`127.0.0.1`}
						},
					},
					Disable: true,
				},
			},
			ParamType:  global_api_strategy.ALL_TYPE,
			Controller: func(apiSession _type2.IApiSession) (i interface{}, info *go_error.ErrorInfo) {
				var params struct{
					Test string `json:"test" validate:"is-mobile"`
				}
				apiSession.ScanParams(&params)
				//return nil, go_error.WrapWithAll(errors.New("haha"), 2000, map[string]interface{}{
				//	"haha": "u7ytu7",
				//})
				return params.Test, nil
			},
			Params: struct {
				Test string `json:"test" validate:"required,is-mobile"`
			}{},
		},
	})
	global_api_strategy3.GlobalRateLimitStrategy.SetErrorCode(10000)
	global_api_strategy2.GlobalApiStrategyDriverInstance.Register(global_api_strategy2.GlobalStrategyData{
		Strategy: &global_api_strategy3.GlobalRateLimitStrategy,
		Param:    global_api_strategy3.GlobalRateLimitStrategyParam{
			FillInterval: 1000 * time.Millisecond,
		},
		Disable:  false,
	})
	service.Service.SetPort(8080)

	err := service.Service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

```sh
$ go run .
```

```shell script
curl -H "Content-Type: application/json" http://0.0.0.0:8080/api/test/v1/test_api -d "{}"

{"msg":"Key: 'Test'; Error:Field validation for 'test' failed on the 'required' tag; sql-inject-check,required,is-mobile","code":2005,"data":{"field":"test"}}

curl -H "Content-Type: application/json" http://0.0.0.0:8080/api/test/v1/test_api -d "{\"test\": \"yrte\"}"

{"msg":"Key: 'Test'; Error:Field validation for 'test' failed on the 'is-mobile' tag; sql-inject-check,required,is-mobile","code":2005,"data":{"field":"test"}}

curl -H "Content-Type: application/json" http://0.0.0.0:8080/api/test/v1/test_api -d "{\"test\": \"18317034426\"}"

{"msg":"","code":0,"data":"18317034426"}
```

## Document

[doc](https://godoc.org/github.com/pefish/go-core)

## Security Vulnerabilities

If you discover a security vulnerability within Go-Core, please send an e-mail to [pefish@qq.com](mailto:pefish@qq.com). All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the [BSD 3-clause license](LICENSE), just like the Go project itself.
