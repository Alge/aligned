# Pytest Connector [IMPLEMENTS: Test Framework Connector Interface]

The Pytest connector integrates Aligned with Python's pytest testing framework. It uses pytest's native `--collect-only` flag to discover tests, ensuring 100% accuracy and respecting pytest configuration files.

## Framework Detection

### Detect framework presence

Check if the `pytest` command is available in PATH using `exec.LookPath()`. Return true if found, false if not found. Return error only for unexpected failures during detection (not for absence of pytest).

**Test:** `Alge/aligned/internal/connectors.TestPytestDetectFramework`

## Configuration Initialization

### Generate default configuration

Create a ConnectorConfig with type "pytest", executable "pytest", and the provided path. Can be initialized via `align init python-pytest [path]`.

**Test:** `Alge/aligned/internal/connectors.TestPytestGenerateConfig`

### List in init help

The python-pytest connector appears in `align init help` output with its name and description.

**Test:** `Alge/aligned/cmd/align.TestInitListsPytestConnector`

## Command Integration

### Register in check command

The pytest connector is registered in the check command, allowing configurations with type "pytest" to successfully discover tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestPytestConnectorRegisteredInCheck`

### Register in list-tests command

The connector is registered in the list-tests command, allowing configurations with this connector type to successfully list tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestPytestConnectorRegisteredInListTests`

## Test Discovery

### Discover tests in project

Execute `pytest --collect-only -q` in the specified path to discover all tests. Parse the output to extract fully-qualified test node IDs in pytest format (e.g., `tests/test_auth.py::TestLogin::test_valid_credentials`). Return the list of test names without package prefixes or modification.

**Test:** `Alge/aligned/internal/connectors.TestPytestDiscoverTests`

### Handle nested directories

Correctly discover tests in nested directory structures such as `tests/unit/auth/` and `tests/integration/api/handlers/`. The test node IDs preserve the full path relative to the project root, including all directory levels and class hierarchies.

**Test:** `Alge/aligned/internal/connectors.TestPytestDiscoverTestsNestedDirectories`

### Handle empty test suite gracefully

When a Python project contains no test files or no test functions, return an empty list without error. This is a valid state, not a failure condition.

**Test:** `Alge/aligned/internal/connectors.TestPytestEmptyTestSuite`

### Report framework not found

When the pytest executable is not found in PATH, return a clear error message indicating which executable was not found and suggesting installation steps (e.g., `pip install pytest`).

**Test:** `Alge/aligned/internal/connectors.TestPytestFrameworkNotFound`

### Report invalid project structure

Return a clear error when pytest cannot collect tests due to structural issues such as:
- Import errors from missing `__init__.py` files in test packages
- Invalid PYTHONPATH preventing module imports
- Missing required project dependencies
- Syntax errors in conftest.py files

The error message distinguishes between collection errors and other failure types.

**Test:** `Alge/aligned/internal/connectors.TestPytestInvalidProjectStructure`

### Handle discovery errors

Return meaningful errors when test discovery fails due to:
- Import errors in test files or their dependencies
- Syntax errors in Python test code
- Missing test dependencies
- Fixture definition errors
- Duplicate test node IDs

Error messages include relevant context from pytest's output to aid debugging. Different error types are distinguishable from the error message content.

**Test:** `Alge/aligned/internal/connectors.TestPytestDiscoveryErrors`
