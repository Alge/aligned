# Pytest Sample Project

This is a minimal Python project using pytest, demonstrating Aligned's pytest connector.

## Setup

```bash
# Create virtual environment (optional but recommended)
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

## Running Tests

```bash
# Run all tests
pytest

# Run tests with verbose output
pytest -v

# Collect tests without running (what Aligned uses)
pytest --collect-only -q
```

## Project Structure

The project demonstrates:
- **Class-based test organization** - Test classes for grouping related tests
- **Pytest fixtures and assertions** - Using pytest's powerful features
- **Exception testing** - Testing error conditions with `pytest.raises()`

## Test Discovery Output

Running `pytest --collect-only -q` produces:

```
tests/test_calculator.py::TestAddition::test_add_positive_numbers
tests/test_calculator.py::TestAddition::test_add_negative_numbers
tests/test_calculator.py::TestSubtraction::test_subtract_numbers
tests/test_calculator.py::TestMultiplication::test_multiply_numbers
tests/test_calculator.py::TestDivision::test_divide_numbers
tests/test_calculator.py::TestDivision::test_divide_by_zero
```

Aligned uses these node IDs directly as test references.

## Using with Aligned

From the repository root:

```bash
# Build Aligned
make build

# Initialize Aligned config for this project
cd test-projects/pytest-sample
../../bin/align init python-pytest .

# Check specification coverage
../../bin/align check spec.md
```
