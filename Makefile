run:
	go run ./cmd/servelocal/serve.go
lint:
	golangci-lint run ./...