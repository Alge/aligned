# Go Test Sample Project

This is a minimal Go project used for testing and demonstrating Aligned's Go test connector.

## Setup

No additional setup needed if Go is installed:

```bash
# Verify Go installation
go version
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# List tests (what Aligned uses)
go test -list=. ./...
```

## Project Structure

The project demonstrates:
- **Package-based organization** - Tests in `calculator` package
- **Table-driven tests** - Multiple test functions
- **Standard Go testing** - Using `testing.T`

## Test Discovery Output

Running `go test -list=. ./...` produces:

```
TestAdd
TestAddNegative
TestSubtract
TestMultiply
ok  	example/calculator	0.001s
```

Aligned converts this to package-qualified test names:
```
calculator.TestAdd
calculator.TestAddNegative
calculator.TestSubtract
calculator.TestMultiply
```

## Using with Aligned

From the repository root:

```bash
# Build Aligned
make build

# Initialize Aligned config for this project
cd test-projects/go-test-sample
../../bin/align init go-test .

# Check specification coverage
../../bin/align check spec.md
```
