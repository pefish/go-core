module test

go 1.13

replace github.com/pefish/go-core => ../

require (
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/kataras/iris v11.1.0+incompatible
	github.com/pefish/go-application v0.1.1
	github.com/pefish/go-config v0.1.4
	github.com/pefish/go-core v0.0.0-00010101000000-000000000000
	github.com/pefish/go-error v0.3.4
	github.com/pefish/go-logger v0.1.15
	go.opencensus.io v0.22.1
)
