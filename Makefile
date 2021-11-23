run:
	go run ./cmd/serve/main.go
lint:
	golangci-lint run ./...