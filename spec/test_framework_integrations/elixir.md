# Elixir Connector [IMPLEMENTS: Test Framework Connector Interface]

The Elixir connector integrates Aligned with Elixir's ExUnit testing framework. It uses Mix's `test --trace` command to discover tests, ensuring accurate test identification while respecting Mix project configuration.

## Framework Detection

### Detect framework presence

Check if the `mix` command is available in PATH using `exec.LookPath()`. Return true if found, false if not found. Return error only for unexpected failures during detection (not for absence of mix).

**Test:** `Alge/aligned/internal/connectors.TestElixirDetectFramework`

## Configuration Initialization

### Generate default configuration

Create a ConnectorConfig with type "elixir", executable "mix", and the provided path. Can be initialized via `align init elixir-exunit [path]`.

**Test:** `Alge/aligned/internal/connectors.TestElixirGenerateConfig`

### List in init help

The elixir-exunit connector appears in `align init help` output with its name and description.

**Test:** `Alge/aligned/cmd/align.TestInitListsElixirConnector`

## Test Discovery

### Discover tests in project

Execute `mix test --trace` in the specified path to discover all tests. Parse the trace output to extract test identifiers in the format `file:Module:test name` (e.g., `test/sample_test.exs:SampleTest:test greets the world`). Return the list of test names with file path and module context.

**Test:** `Alge/aligned/internal/connectors.TestElixirDiscoverTests`

### Handle nested directories

Correctly discover tests in nested directory structures such as `test/unit/auth/` and `test/integration/api/handlers/`. The test identifiers preserve the full path relative to the project root, including all directory levels and module hierarchies.

**Test:** `Alge/aligned/internal/connectors.TestElixirDiscoverTestsNestedDirectories`

### Handle empty test suite gracefully

When an Elixir project contains no test files or no test cases, return an empty list without error. Mix outputs "There are no tests to run" which is a valid state, not a failure condition.

**Test:** `Alge/aligned/internal/connectors.TestElixirEmptyTestSuite`

### Report framework not found

When the mix executable is not found in PATH, return a clear error message indicating which executable was not found and suggesting installation steps (e.g., install Elixir which includes mix).

**Test:** `Alge/aligned/internal/connectors.TestElixirFrameworkNotFound`

### Report invalid project structure

Return a clear error when mix cannot run tests due to structural issues such as:
- Missing `mix.exs` file (not a Mix project)
- Invalid Mix project configuration
- Missing test helper files
- Compilation errors in project files

The error message distinguishes between project structure errors and compilation errors.

**Test:** `Alge/aligned/internal/connectors.TestElixirInvalidProjectStructure`

### Handle discovery errors

Return meaningful errors when test discovery fails due to:
- Compilation errors in test files
- Syntax errors in Elixir test code
- Missing dependencies in mix.exs
- Module definition errors
- Invalid test macros or setup blocks

Error messages include relevant context from Mix's output to aid debugging. Different error types are distinguishable from the error message content.

**Test:** `Alge/aligned/internal/connectors.TestElixirDiscoveryErrors`
