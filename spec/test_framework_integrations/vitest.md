# Vitest Connector [IMPLEMENTS: Test Framework Connector Interface]

The Vitest connector integrates Aligned with Vitest testing framework. It uses Vitest's native `list --json` command to discover tests, ensuring 100% accuracy and respecting Vitest configuration files.

## Framework Detection

### Detect framework presence

Check if the `vitest` command is available in PATH using `exec.LookPath()`. Return true if found, false if not found. Return error only for unexpected failures during detection (not for absence of vitest).

**Test:** `Alge/aligned/internal/connectors.TestVitestDetectFramework`

## Configuration Initialization

### Generate default configuration

Create a ConnectorConfig with type "vitest", executable "vitest", and the provided path. Can be initialized via `align init javascript-vitest [path]`.

**Test:** `Alge/aligned/internal/connectors.TestVitestGenerateConfig`

### List in init help

The javascript-vitest connector appears in `align init help` output with its name and description.

**Test:** `Alge/aligned/cmd/align.TestInitListsVitestConnector`

## Command Integration

### Register in check command

The vitest connector is registered in the check command, allowing configurations with type "vitest" to successfully discover tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestVitestConnectorRegisteredInCheck`

### Register in list-tests command

The connector is registered in the list-tests command, allowing configurations with this connector type to successfully list tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestVitestConnectorRegisteredInListTests`

## Test Discovery

### Discover tests in project

Execute `vitest list --json` in the specified path to discover all tests. Parse the JSON output to extract test names in the format `{relative_file_path} > {test_name}`. The test name includes the describe hierarchy (e.g., `src/example.test.js > Math operations > Addition > adds 1 + 2 to equal 3`). Return the list of fully-qualified test identifiers.

**Test:** `Alge/aligned/internal/connectors.TestVitestDiscoverTests`

### Handle nested directories

Correctly discover tests in nested directory structures such as `src/components/auth/` and `tests/integration/api/`. The test identifiers preserve the full relative path from the project root.

**Test:** `Alge/aligned/internal/connectors.TestVitestDiscoverTestsNestedDirectories`

### Handle empty test suite gracefully

When a JavaScript project contains no test files or no test functions, return an empty JSON array `[]` without error. This is a valid state, not a failure condition.

**Test:** `Alge/aligned/internal/connectors.TestVitestEmptyTestSuite`

### Report framework not found

When the vitest executable is not found in PATH, return a clear error message indicating which executable was not found and suggesting installation steps (e.g., `npm install -D vitest`).

**Test:** `Alge/aligned/internal/connectors.TestVitestFrameworkNotFound`

### Report invalid project structure

Return a clear error when vitest cannot collect tests due to structural issues such as:
- Import errors from missing dependencies
- Invalid module imports
- Missing required project dependencies
- Syntax errors in test files
- Invalid vitest configuration

The error message distinguishes between collection errors and other failure types.

**Test:** `Alge/aligned/internal/connectors.TestVitestInvalidProjectStructure`

### Handle discovery errors

Return meaningful errors when test discovery fails due to:
- Import errors in test files or their dependencies
- Syntax errors in JavaScript/TypeScript test code
- Missing test dependencies
- Invalid JSON output from vitest
- Malformed test structure

Error messages include relevant context from vitest's output to aid debugging. Different error types are distinguishable from the error message content.

**Test:** `Alge/aligned/internal/connectors.TestVitestDiscoveryErrors`
