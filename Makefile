
test:
	go mod tidy
	go build ./...
	go test -cover ./...
