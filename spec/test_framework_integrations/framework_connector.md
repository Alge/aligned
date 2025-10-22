# Test Framework Connector Interface [INTERFACE]

The connector interface defines how Aligned integrates with different test frameworks. Connectors must implement two methods: framework detection and test discovery.

## Framework Detection

### Detect framework presence

Check if the test framework is available and the project uses this framework. Return true if the framework is detected, false if not, and an error only for unexpected failures during detection.

## Configuration Initialization

### Generate default configuration

Create a default connector configuration for the framework. Takes a path to the project root and returns a ConnectorConfig with appropriate defaults for the framework (type, executable, and path). The connector can be initialized via `align init [language]-[framework] [path]` where the language-framework combination uniquely identifies the connector. All available connectors are listed when running `align init help` or `align init` without parameters.

## Test Discovery

### Discover tests in project

Execute framework-specific discovery command and return a list of test names in the framework's native format. Takes a path to the project root and returns all discoverable tests.

### Handle nested directories

Correctly discover tests in nested directory structures, returning tests from all subdirectories within the project.

### Handle empty test suite gracefully

When a project has no tests, return an empty list rather than an error. This is a valid state, not a failure.

### Report framework not found

Return a clear error when the framework executable is not found in PATH or the expected location. The error message must indicate which framework was not found.

### Report invalid project structure

Return a clear error when the project structure doesn't match framework expectations (e.g., missing configuration files, invalid module structure).

### Handle discovery errors

Return meaningful error messages when test discovery fails due to compilation errors, permission issues, or other problems. Distinguish between different error types in the error message.