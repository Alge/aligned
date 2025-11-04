# Go Connector [IMPLEMENTS: Test Framework Connector Interface]

The Go connector integrates Aligned with Go's built-in testing framework using `go test -list` for test discovery.

## Framework Detection

### Detect framework presence

Check if the `go` command is available in PATH. Return true if found, false if not found, error only for unexpected detection failures.

**Test:** `Alge/aligned/internal/connectors.TestGoDetectFramework`

## Configuration Initialization

### Generate default configuration

Create a ConnectorConfig with type "go", executable "go", and the provided path. Can be initialized via `align init go-test [path]`.

**Test:** `Alge/aligned/internal/connectors.TestGoGenerateConfig`

### List in init help

The go-test connector appears in `align init help` output with its name and description.

**Test:** `Alge/aligned/cmd/align.TestInitListsGoConnector`

## Command Integration

### Register in check command

The go connector is registered in the check command, allowing configurations with type "go" to successfully discover tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestGoConnectorRegisteredInCheck`

### Register in list-tests command

The connector is registered in the list-tests command, allowing configurations with this connector type to successfully list tests without "unsupported connector type" errors.

**Test:** `Alge/aligned/cmd/align.TestGoConnectorRegisteredInListTests`

## Test Discovery

### Discover tests in project

Run `go test -list=. ./...` to list all tests in the project and its subdirectories, returning package-qualified test names.

**Test:** `Alge/aligned/internal/connectors.TestGoDiscoverTests`

### Handle nested directories

Discover tests in nested packages like `internal/auth` and `cmd/server/handlers/api`, preserving package paths in test names.

**Test:** `Alge/aligned/internal/connectors.TestGoDiscoverTestsNestedDirectories`

### Handle empty test suite gracefully

When a Go project has no tests, return an empty list without error.

**Test:** `Alge/aligned/internal/connectors.TestGoEmptyTestSuite`

### Report framework not found

Return a clear error when the `go` command is not found in PATH, indicating which executable was not found.

**Test:** `Alge/aligned/internal/connectors.TestGoFrameworkNotFound`

### Report invalid project structure

Return a clear error when go.mod is missing or the project structure is invalid, helping users understand what's wrong with their setup.

**Test:** `Alge/aligned/internal/connectors.TestGoInvalidProjectStructure`

### Handle discovery errors

Return meaningful errors when test discovery fails due to compilation errors, permission issues, or other problems. Error messages distinguish between different failure types.

**Test:** `Alge/aligned/internal/connectors.TestGoDiscoveryErrors`