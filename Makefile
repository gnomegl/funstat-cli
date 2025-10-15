.PHONY: build test clean install

# Build the CLI binary
build:
	go build -o bin/funstat ./cmd/funstat

# Install the CLI to GOPATH/bin
install:
	go install ./cmd/funstat

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run the CLI with debug mode
run-debug:
	go run ./cmd/funstat --debug $(ARGS)

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Update dependencies
deps:
	go mod tidy
	go mod download