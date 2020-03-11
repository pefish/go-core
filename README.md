# Go-Core Web Framework

[![view examples](https://img.shields.io/badge/learn%20by-examples-0C8EC5.svg?style=for-the-badge&logo=go)](https://github.com/pefish/go-core/tree/master/_examples)

Go-Core is a fast, simple and young web framework for Go.

It provides some features. like
1. Api Strategy
2. Swagger Generator
3. Return Hook
4. Cors
5. Api Error Catch
6. And so on

## 📖 Learning Iris

<details>
<summary>Quick start</summary>

```go
package main

import (
	go_core "github.com/pefish/go-core"
	"github.com/pefish/go-core/api"
	api_session "github.com/pefish/go-core/api-session"
	api_strategy "github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/driver/logger"
	global_api_strategy "github.com/pefish/go-core/global-api-strategy"
	go_logger "github.com/pefish/go-logger"
	"time"
)

func PostTest(apiSession *api_session.ApiSessionClass) interface{} {
	return `haha, this is return value`
}

func main() {
	go_core.Service.SetName(`test service`) // set service name
	logger.LoggerDriver.Register(go_logger.NewLogger(go_logger.WithIsDebug(true))) // register logger
	go_core.Service.SetPath(`/api/test`)
	global_api_strategy.ParamValidateStrategy.SetErrorCode(2005)
	go_core.Service.SetRoutes([]*api.Api{
		{
			Description: "this is a test api",
			Path:        "/v1/test_api",
			Method:      `POST`,
			Strategies: []api_strategy.StrategyData{
				{
					Strategy: &api_strategy.RateLimitApiStrategy,
					Param: api_strategy.RateLimitParam{
						Limit: 1 * time.Second,
					},
				},
			},
			ParamType:  global_api_strategy.ALL_TYPE,
			Controller: PostTest,
		},
	})
	go_core.Service.SetPort(8080)

	go_core.Service.Run()
}
```

```sh
$ go run main.go
```

</details>

## 🛡 Security Vulnerabilities

If you discover a security vulnerability within Go-Core, please send an e-mail to [pefish@qq.com](mailto:pefish@qq.com). All security vulnerabilities will be promptly addressed.

## 📝 License

This project is licensed under the [BSD 3-clause license](LICENSE), just like the Go project itself.
