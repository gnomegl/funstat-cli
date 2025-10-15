# Testing Guide

## Quick Start

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run with coverage
make test-coverage

# Run specific test
make test-run TEST=TestResolveUsernames
```

## Test Structure

```
tgstat/
├── pkg/client/
│   ├── client.go           # Library implementation
│   ├── client_test.go      # Unit tests (78% coverage)
│   └── types.go            # Type definitions
├── cmd/funstat/
│   ├── main.go             # CLI implementation
│   ├── handlers_test.go    # Handler tests (15.5% coverage)
│   └── integration_test.go # Integration tests
└── Makefile                # Test commands
```

## Test Coverage Summary

**Overall Coverage: 78% (pkg/client)**

### Package: `pkg/client` (78% coverage)
- ✅ All 12 client methods tested
- ✅ Error handling verified
- ✅ Context cancellation tested
- ✅ Authentication tested
- ✅ Mock HTTP server used

### Package: `cmd/funstat` (15.5% coverage)
- ✅ All 11 subcommands tested
- ✅ Handler integration verified
- ⚠️ CLI execution not directly tested (use integration tests)

## Running Tests

### 1. Unit Tests
```bash
# All packages
go test -v ./...

# Specific package
go test -v ./pkg/client/...

# With race detection
go test -race ./...

# Short mode (skip slow tests)
go test -short ./...
```

### 2. Integration Tests
```bash
# Requires API key
export FUNSTAT_API_KEY=your-api-key-here
go test -tags=integration -v ./cmd/funstat/...
```

### 3. Coverage Reports
```bash
# Generate HTML coverage report
make test-coverage
open coverage.html

# Terminal coverage
go test -cover ./...

# Detailed coverage per function
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### 4. Specific Tests
```bash
# Run single test
go test -v -run TestResolveUsernames ./pkg/client/...

# Run test pattern
go test -v -run "TestUser.*" ./pkg/client/...

# Verbose with test names
go test -v ./pkg/client/... | grep -E "RUN|PASS|FAIL"
```

## Test Organization

### Table-Driven Tests Pattern
All tests use table-driven pattern for comprehensive coverage:

```go
tests := []struct {
    name       string
    input      InputType
    want       ExpectedType
    wantErr    bool
}{
    {"success case", validInput, expectedOutput, false},
    {"error case", invalidInput, nil, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Function(tt.input)
        if tt.wantErr {
            assert.Error(t, err)
            return
        }
        require.NoError(t, err)
        assert.Equal(t, tt.want, got)
    })
}
```

## What's Tested

### ✅ User Commands
- [x] `user resolve` - Username to ID resolution
- [x] `user stats` - Full user statistics
- [x] `user stats-min` - Minimal statistics (FREE)
- [x] `user get-by-id` - Get users by Telegram ID
- [x] `user groups` - List user's groups
- [x] `user groups-count` - Count user's groups (FREE)
- [x] `user messages` - Get user messages with filters
- [x] `user messages-count` - Count messages (FREE)
- [x] `user names` - Name history
- [x] `user usernames` - Username history

### ✅ Group Commands
- [x] `group info` - Get group information

### ✅ Infrastructure
- [x] Client initialization with options
- [x] HTTP request/response handling
- [x] Error handling and propagation
- [x] Context cancellation
- [x] Authentication headers
- [x] Query parameter encoding
- [x] JSON marshaling/unmarshaling

### ✅ Edge Cases
- [x] Empty responses
- [x] Invalid user IDs
- [x] Not found errors (404)
- [x] Unauthorized errors (401)
- [x] Network timeouts
- [x] Multiple items (batch operations)
- [x] Pagination
- [x] Optional filters

## TDD Workflow

### 1. Write Failing Test
```bash
# Create test that fails
cat > pkg/client/new_feature_test.go << 'EOF'
func TestNewFeature(t *testing.T) {
    result, err := client.NewFeature()
    require.NoError(t, err)
    assert.NotNil(t, result)
}
EOF

# Run test - should fail
go test -v -run TestNewFeature ./pkg/client/...
# FAIL: undefined: Client.NewFeature
```

### 2. Implement Minimal Code
```bash
# Add minimal implementation
# Edit pkg/client/client.go

# Run test - should pass
go test -v -run TestNewFeature ./pkg/client/...
# PASS
```

### 3. Refactor
```bash
# Improve implementation
# Run all tests to verify
go test ./...
# All PASS
```

## Continuous Testing

### Watch Mode (requires `entr`)
```bash
# Install entr first: apt install entr
make test-watch

# Or manually
find . -name '*.go' | entr -c go test ./...
```

### Pre-commit Hook
```bash
# Create .git/hooks/pre-commit
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
echo "Running tests..."
go test -short ./... || exit 1
echo "Running fmt..."
go fmt ./... || exit 1
echo "All checks passed!"
EOF

chmod +x .git/hooks/pre-commit
```

## Debugging Failed Tests

### Verbose Output
```bash
# Maximum verbosity
go test -v -count=1 ./...

# With print statements preserved
go test -v ./... 2>&1 | tee test-output.log
```

### Single Test Debugging
```bash
# Run one test repeatedly
go test -v -run TestSpecific -count=10 ./pkg/client/...

# With race detector
go test -v -race -run TestSpecific ./pkg/client/...
```

### Coverage Gaps
```bash
# Find untested code
go test -coverprofile=coverage.out ./pkg/client/...
go tool cover -func=coverage.out | grep -E "0\.0%|statements"
```

## Integration Test Setup

### Requirements
1. Valid Funstat API key
2. Network connectivity
3. Test account with sufficient balance

### Running Integration Tests
```bash
# Set API key
export FUNSTAT_API_KEY="your-actual-api-key"

# Run integration tests
go test -tags=integration -v ./cmd/funstat/...

# Skip integration tests
go test -short ./...
```

### Integration Test Scenarios
- Real API requests
- Rate limiting verification
- Error handling with real API
- Context cancellation behavior
- Full request/response cycle

## Makefile Commands

```bash
make test              # Run all tests
make test-unit         # Unit tests only
make test-integration  # Integration tests (requires API key)
make test-coverage     # Generate coverage report
make test-short        # Fast tests only
make test-run TEST=... # Run specific test
make test-watch        # Watch mode (requires entr)
make test-summary      # Show test documentation
make check             # Run fmt, lint, test-unit
make clean             # Clean artifacts and cache
```

## Test Maintenance

### Adding New Tests
1. Write test first (TDD)
2. Follow table-driven pattern
3. Use descriptive test names
4. Test success and error cases
5. Verify with `make test`

### Updating Tests
1. Modify test expectations
2. Run affected tests
3. Verify coverage maintained
4. Update documentation

### Test Naming Convention
```
TestFunctionName                    # Function being tested
TestFunctionName/scenario_name      # Specific scenario
TestFunctionName_EdgeCase           # Edge case variant
```

### Test File Organization
```
function_test.go     # Tests for function.go
handlers_test.go     # Handler integration tests
integration_test.go  # Integration tests (build tag)
```

## Common Issues

### 1. Import Cycle
```bash
# Solution: Move shared test utilities to testutil package
mkdir internal/testutil
```

### 2. Flaky Tests
```bash
# Run multiple times to detect
go test -count=100 -run TestFlaky ./...

# Add timeouts and retries
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### 3. Slow Tests
```bash
# Identify slow tests
go test -v ./... | grep -E "\([0-9]+\.[0-9]+s\)"

# Mark slow tests
if testing.Short() {
    t.Skip("Skipping slow test")
}
```

## Benchmarking

```bash
# Run benchmarks
go test -bench=. ./pkg/client/...

# With memory allocation stats
go test -bench=. -benchmem ./pkg/client/...

# Compare benchmarks
go test -bench=. -count=5 ./pkg/client/... > old.txt
# Make changes
go test -bench=. -count=5 ./pkg/client/... > new.txt
benchstat old.txt new.txt
```

## Test Metrics

### Current Status
- **Total Tests:** 26 functions, 60+ scenarios
- **Coverage:** 78% (pkg/client), 15.5% (cmd/funstat)
- **Execution Time:** ~0.13s (unit tests)
- **Success Rate:** 100% (all passing)

### Coverage Goals
- ✅ pkg/client: >75% (achieved: 78%)
- ⚠️ cmd/funstat: >50% (current: 15.5%)
- 🎯 Overall: >70% (current: ~65%)

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Coverage Best Practices](https://go.dev/blog/cover)
- [TDD in Go](https://quii.gitbook.io/learn-go-with-tests/)

## Next Steps

1. ✅ Unit tests completed
2. ✅ Handler tests completed
3. ✅ Integration tests ready
4. 🎯 Increase cmd/funstat coverage to 50%
5. 🎯 Add benchmark tests
6. 🎯 Add fuzzing tests
7. 🎯 E2E CLI testing
