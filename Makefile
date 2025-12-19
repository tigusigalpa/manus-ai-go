.PHONY: test test-coverage lint fmt vet build clean examples

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Build examples
build-examples:
	cd examples/basic && go build -o ../../bin/basic main.go
	cd examples/file-upload && go build -o ../../bin/file-upload main.go
	cd examples/webhook && go build -o ../../bin/webhook main.go
	@echo "Examples built in bin/"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run all checks
check: fmt vet test

# Install dependencies
deps:
	go mod download
	go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  build-examples  - Build example programs"
	@echo "  clean           - Clean build artifacts"
	@echo "  check           - Run fmt, vet, and test"
	@echo "  deps            - Install dependencies"
