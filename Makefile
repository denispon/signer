.PHONY: test build clean help

help:
	@echo "Available targets:"
	@echo "  test      - Run all tests with verbose output and coverage"
	@echo "  build     - Build the project"
	@echo "  clean     - Clean build artifacts"

test:
	@echo "Running tests..."
	go test ./... -v -cover

build:
	@echo "Building project..."
	go build ./...

clean:
	@echo "Cleaning build artifacts..."
	go clean

.DEFAULT_GOAL := help
