# Gleam Connector [IMPLEMENTS: Test Framework Connector Interface]

The Gleam connector integrates Aligned with Gleam's gleeunit testing framework. It discovers tests by parsing Gleam source files in the `test/` directory to find public functions ending in `_test`.

## Framework Detection

### Detect framework presence

Check if the `gleam` command is available in PATH using `exec.LookPath()`. Return true if found, false if not found. Return error only for unexpected failures during detection (not for absence of gleam).

**Test:** `Alge/aligned/internal/connectors.TestGleamDetectFramework`

## Configuration Initialization

### Generate default configuration

Create a ConnectorConfig with type "gleam", executable "gleam", and the provided path. Can be initialized via `align init gleam-gleeunit [path]`.

**Test:** `Alge/aligned/internal/connectors.TestGleamGenerateConfig`

### List in init help

The gleam-gleeunit connector appears in `align init help` output with its name and description.

**Test:** `Alge/aligned/cmd/align.TestInitListsGleamConnector`

## Command Integration

### Register in check command

The gleam connector is registered in the check command, allowing configurations with type "gleam" to successfully discover tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestGleamConnectorRegisteredInCheck`

### Register in list-tests command

The connector is registered in the list-tests command, allowing configurations with this connector type to successfully list tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestGleamConnectorRegisteredInListTests`

## Test Discovery

### Discover tests in project

Find all `.gleam` files in the `test/` directory and parse them to identify public functions ending in `_test`. Return test names in the format `module_name.function_name` (e.g., `test_discovery_sample_test.hello_world_test`, `math_test.multiply_test`). Module names are derived from file paths by removing the `.gleam` extension and replacing directory separators with `@` for nested paths.

**Test:** `Alge/aligned/internal/connectors.TestGleamDiscoverTests`

### Handle nested directories

Correctly discover tests in nested directory structures such as `test/unit/auth/` and `test/integration/api/handlers/`. Test module names preserve the full path structure (e.g., `unit@auth@login_test.authenticate_user_test`).

**Test:** `Alge/aligned/internal/connectors.TestGleamDiscoverTestsNestedDirectories`

### Handle empty test suite gracefully

When a Gleam project contains no test files or no test functions ending in `_test`, return an empty list without error. This is a valid project state, not a failure condition.

**Test:** `Alge/aligned/internal/connectors.TestGleamEmptyTestSuite`

### Report framework not found

When the gleam executable is not found in PATH, return a clear error message indicating which executable was not found and suggesting installation steps (e.g., install Gleam via asdf, homebrew, or the official installer).

**Test:** `Alge/aligned/internal/connectors.TestGleamFrameworkNotFound`

### Report invalid project structure

Return a clear error when the project structure doesn't match Gleam expectations:
- Missing `gleam.toml` file (not a Gleam project)
- Invalid Gleam project configuration
- Test directory exists but contains no `.gleam` files

The error message distinguishes between missing project configuration and empty test directories.

**Test:** `Alge/aligned/internal/connectors.TestGleamInvalidProjectStructure`

### Handle discovery errors

Return meaningful errors when test discovery fails due to:
- File permission issues preventing reading of test files
- Invalid Gleam syntax in test files (parse errors)
- Test files that cannot be read or accessed

Error messages include relevant context to aid debugging. Different error types are distinguishable from the error message content.

**Test:** `Alge/aligned/internal/connectors.TestGleamDiscoveryErrors`
