# ExUnit Sample Project

This is a minimal Elixir project using ExUnit, demonstrating Aligned's ExUnit connector.

## Setup

```bash
# Install dependencies
mix deps.get

# Compile the project
mix compile
```

## Running Tests

```bash
# Run all tests
mix test

# Run tests with verbose output
mix test --trace

# List tests (what Aligned uses)
mix test --list
```

## Project Structure

The project demonstrates:
- **Describe blocks** - Grouping related tests with `describe`
- **Pattern matching in tests** - Using Elixir's pattern matching for assertions
- **ExUnit.Case** - Standard ExUnit test structure

## Test Discovery Output

Running `mix test --list` produces:

```
test/calculator_test.exs:6 - CalculatorTest Addition adds positive numbers
test/calculator_test.exs:10 - CalculatorTest Addition adds negative numbers
test/calculator_test.exs:16 - CalculatorTest Subtraction subtracts numbers correctly
test/calculator_test.exs:22 - CalculatorTest Multiplication multiplies numbers correctly
test/calculator_test.exs:28 - CalculatorTest Division divides numbers correctly
test/calculator_test.exs:32 - CalculatorTest Division handles division by zero
```

Aligned converts this to:
```
test/calculator_test.exs:6
test/calculator_test.exs:10
test/calculator_test.exs:16
test/calculator_test.exs:22
test/calculator_test.exs:28
test/calculator_test.exs:32
```

## Using with Aligned

From the repository root:

```bash
# Build Aligned
make build

# Initialize Aligned config for this project
cd test-projects/exunit-sample
../../bin/align init elixir-exunit .

# Check specification coverage
../../bin/align check spec.md
```
