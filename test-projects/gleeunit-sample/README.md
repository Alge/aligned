# Gleeunit Sample Project

This is a minimal Gleam project using gleeunit, demonstrating Aligned's Gleam connector.

## Setup

```bash
# Download dependencies
gleam deps download

# Build the project
gleam build
```

## Running Tests

```bash
# Run all tests
gleam test

# Build and run tests separately
gleam build
gleam test
```

## Project Structure

The project demonstrates:
- **Public test functions** - Functions ending in `_test` are automatically discovered
- **Pattern matching** - Using Gleam's pattern matching and Result types
- **Pipe operator** - Chaining assertions with `|>`

## Test Discovery

Aligned discovers tests by parsing `.gleam` files in the `test/` directory to find public functions ending in `_test`.

Test names are in the format `module_name.function_name`:
```
calculator_test.add_positive_numbers_test
calculator_test.add_negative_numbers_test
calculator_test.subtract_numbers_test
calculator_test.multiply_numbers_test
calculator_test.divide_numbers_test
calculator_test.divide_by_zero_test
```

## Using with Aligned

From the repository root:

```bash
# Build Aligned
make build

# Initialize Aligned config for this project
cd test-projects/gleeunit-sample
../../bin/align init gleam-gleeunit .

# Check specification coverage
../../bin/align check spec.md
```
