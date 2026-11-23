.PHONY: test build build-cli clean help

help:
	@echo "Available targets:"
	@echo "  test       - Run all tests with verbose output and coverage"
	@echo "  build      - Build the project libraries"
	@echo "  build-cli  - Build the signer CLI binary"
	@echo "  clean      - Clean build artifacts"

test:
	@echo "Running tests..."
	go test ./... -v -cover

build:
	@echo "Building project libraries..."
	go build ./...

build-cli:
	@echo "Building signer CLI binary..."
	go build -o bin/signer main.go
	@echo "Built: bin/signer"

clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -rf bin/

.DEFAULT_GOAL := help
