.PHONY: build test tidy

build:
	go build ./cmd/threatctl

tidy:
	go mod tidy

test:
	go test ./...
