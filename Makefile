
test: api-session/type/type.go
	mockgen github.com/pefish/go-http IHttp > mock/mock-go-http/i_http.go

	go build ./...
	go test -cover ./...

api-session/type/type.go:
	mockgen github.com/pefish/go-core-type/api-session IApiSession > mock/mock-api-session/mock.go
