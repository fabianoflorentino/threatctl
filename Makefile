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

coverage:
	go test ./... -coverprofile=coverage.out
	@echo "Coverage summary:"
	go tool cover -func=coverage.out | sed -n '/total:/p'

check-coverage:
	@go test ./... -coverprofile=coverage.out >/dev/null || exit 1
	@cov=$(go tool cover -func=coverage.out | awk '/total:/ {print substr($3,1,length($3)-1)}'); \
	if [ -z "$${cov}" ]; then echo "cannot determine coverage"; exit 2; fi; \
	awk -v cov="$${cov}" 'BEGIN { if (cov+0 < 80) { print "coverage too low: " cov "%"; exit 1 } else { print "coverage OK: " cov "%" } }'
