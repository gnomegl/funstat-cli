.PHONY: build test test-unit test-integration test-coverage test-short clean install

# Build the CLI binary
build:
	go build -o bin/funstat ./cmd/funstat

# Install the CLI to GOPATH/bin
install:
	go install ./cmd/funstat

# Run all tests
test:
	go test -v ./...

# Run unit tests only
test-unit:
	go test -v ./pkg/client/...
	go test -v ./cmd/funstat/... -short

# Run integration tests (requires FUNSTAT_API_KEY)
test-integration:
	@if [ -z "$(FUNSTAT_API_KEY)" ]; then \
		echo "Error: FUNSTAT_API_KEY environment variable not set"; \
		exit 1; \
	fi
	go test -tags=integration -v ./cmd/funstat/...

# Run tests with coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run short tests (skip slow tests)
test-short:
	go test -short ./...

# Run specific test
test-run:
	go test -v -run $(TEST) ./...

# Clean build artifacts and test cache
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -testcache

# Run the CLI with debug mode
run-debug:
	go run ./cmd/funstat --debug $(ARGS)

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Run all quality checks
check: fmt lint test-unit
	@echo "All checks passed!"

# Update dependencies
deps:
	go mod tidy
	go mod download

# Show test summary
test-summary:
	@cat TEST_SUMMARY.md

# Run tests in watch mode (requires entr)
test-watch:
	find . -name '*.go' | entr -c make test-unit