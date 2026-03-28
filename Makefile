.PHONY: build test tidy

bin = bin

build: $(bin)
	go build -o $(bin)/threatctl ./cmd/threatctl
	go build -o $(bin)/genpcap ./cmd/genpcap

$(bin):
	mkdir -p $(bin)

tidy:
	go mod tidy

test:
	go test ./...
