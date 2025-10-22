.PHONY: fmt test build clean coverage install lint align help

fmt:
	go fmt ./...

test:
	go test ./...

build:
	go build -o bin/align ./cmd/align

clean:
	rm -rf bin/ coverage.out

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

install:
	go install ./cmd/align

lint:
	golangci-lint run

align:
	./bin/align check spec/

help:
	@echo "Available targets:"
	@echo "  fmt       - Format Go code"
	@echo "  test      - Run tests"
	@echo "  build     - Build align binary"
	@echo "  clean     - Remove build artifacts"
	@echo "  coverage  - Run tests with coverage report"
	@echo "  install   - Install align to GOPATH/bin"
	@echo "  lint      - Run linter"
	@echo "  align     - Validate specs using align"
	@echo "  help      - Show this help message"