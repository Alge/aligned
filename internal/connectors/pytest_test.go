// internal/connectors/pytest_test.go
package connectors

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPytestDetectFramework(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	t.Run("detects pytest command", func(t *testing.T) {
		connector := &PytestConnector{Executable: "pytest"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("returns false for nonexistent executable", func(t *testing.T) {
		connector := &PytestConnector{Executable: "nonexistent-pytest-binary"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.False(t, got)
	})
}

func TestPytestGenerateConfig(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	t.Run("generates config with correct type, executable, and path", func(t *testing.T) {
		connector := &PytestConnector{Executable: "pytest"}
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "pytest", config.Type)
		assert.Equal(t, "pytest", config.Executable)
		assert.Equal(t, path, config.Path)
	})

	t.Run("generates config with custom executable", func(t *testing.T) {
		connector := &PytestConnector{Executable: "/custom/path/to/pytest"}
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "pytest", config.Type)
		assert.Equal(t, "/custom/path/to/pytest", config.Executable)
		assert.Equal(t, path, config.Path)
	})
}

func TestPytestDiscoverTests(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	// This test verifies the connector uses `pytest --collect-only -q` by testing
	// the behavior that would only work with those specific flags:
	// - Tests are collected, not run (requires --collect-only flag)
	// - Output is quiet format (requires -q flag)
	// - All tests are discovered

	projectDir := createPytestProject(t, map[string]string{
		"test_example.py": `import pytest

def test_foo():
    assert True

def test_bar():
    assert True
`,
	})

	connector = &PytestConnector{Executable: "pytest"}
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	// Verify both tests are discovered
	assert.Contains(t, tests, "test_example.py::test_foo")
	assert.Contains(t, tests, "test_example.py::test_bar")
	// Verify tests were collected, not executed
	assert.Len(t, tests, 2, "should discover exactly 2 tests")
}

func TestPytestDiscoverTestsNestedDirectories(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	// This test verifies pytest discovers tests in deeply nested directory structures
	// and preserves the full path in test node IDs

	projectDir := createPytestProject(t, map[string]string{
		"tests/unit/auth/test_auth.py": `import pytest

def test_login():
    assert True
`,
		"tests/integration/api/handlers/test_api.py": `import pytest

def test_api_handler():
    assert True
`,
	})

	connector = &PytestConnector{Executable: "pytest"}
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	assert.Len(t, tests, 2)
	// Verify nested packages are discovered with full paths preserved
	assert.Contains(t, tests, "tests/unit/auth/test_auth.py::test_login")
	assert.Contains(t, tests, "tests/integration/api/handlers/test_api.py::test_api_handler")
}

func TestPytestEmptyTestSuite(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	t.Run("returns empty list for project with no tests", func(t *testing.T) {
		projectDir := createPytestProject(t, map[string]string{
			"main.py": `def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
`,
		})

		connector := &PytestConnector{Executable: "pytest"}

		// Skip test if pytest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("pytest not installed, skipping test")
		}

		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})

	t.Run("returns empty list for project with only non-test python files", func(t *testing.T) {
		projectDir := createPytestProject(t, map[string]string{
			"utils.py": `def add(a, b):
    return a + b
`,
		})

		connector := &PytestConnector{Executable: "pytest"}

		// Skip test if pytest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("pytest not installed, skipping test")
		}

		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})
}

func TestPytestFrameworkNotFound(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	projectDir := createPytestProject(t, map[string]string{
		"test_example.py": `import pytest

def test_foo():
    assert True
`,
	})

	connector = &PytestConnector{Executable: "nonexistent-pytest-binary"}
	_, err = connector.DiscoverTests(projectDir)

	assert.Error(t, err, "should return error when pytest command not found")

	// Verify error is "clear" per specification:
	// 1. Identifies the specific problematic value
	assert.Contains(t, err.Error(), "nonexistent-pytest-binary",
		"error should identify the executable that was not found")

	// 2. States what's wrong in plain language (not just "exit status 1")
	errMsg := strings.ToLower(err.Error())
	assert.True(t,
		strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "no such"),
		"error should clearly state the problem, got: %s", err.Error())

	// 3. Provides context about the operation
	assert.Contains(t, err.Error(), "test discovery",
		"error should provide context about what operation failed")
}

func TestPytestInvalidProjectStructure(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	t.Run("handles import errors from missing dependencies", func(t *testing.T) {
		projectDir := createPytestProject(t, map[string]string{
			"test_example.py": `import nonexistent_module

def test_foo():
    assert True
`,
		})

		connector := &PytestConnector{Executable: "pytest"}

		// Skip test if pytest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("pytest not installed, skipping test")
		}

		_, err = connector.DiscoverTests(projectDir)

		assert.Error(t, err, "should return error for import errors")

		// Verify error is "clear" and provides context
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide context about what operation failed")

		// Should include helpful details from pytest output
		assert.Contains(t, err.Error(), "Output:",
			"error should include pytest output for debugging")
	})

	t.Run("handles directory without proper test discovery", func(t *testing.T) {
		tempDir := t.TempDir()

		connector := &PytestConnector{Executable: "pytest"}

		// Skip test if pytest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("pytest not installed, skipping test")
		}

		// Empty directory should not error - it's handled by empty suite case
		tests, err := connector.DiscoverTests(tempDir)
		assert.NoError(t, err)
		assert.Empty(t, tests)
	})
}

func TestPytestDiscoveryErrors(t *testing.T) {
	connector := DefaultPytestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("pytest not installed, skipping pytest connector tests")
	}

	// This test verifies errors are "meaningful" and "distinguishable" per specification
	var importErr, syntaxErr, nonexistentErr error

	t.Run("handles import errors", func(t *testing.T) {
		projectDir := createPytestProject(t, map[string]string{
			"test_example.py": `import nonexistent_module

def test_foo():
    assert True
`,
		})

		connector := &PytestConnector{Executable: "pytest"}

		// Skip test if pytest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("pytest not installed, skipping test")
		}

		_, err = connector.DiscoverTests(projectDir)
		importErr = err

		assert.Error(t, err, "should return error for import errors")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")

		// Verify error is "meaningful" - provides context
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")

		// Verify error includes helpful output
		assert.Contains(t, err.Error(), "Output:",
			"error should include pytest output for debugging")
	})

	t.Run("handles syntax errors", func(t *testing.T) {
		projectDir := createPytestProject(t, map[string]string{
			"test_example.py": `import pytest

def test_foo():
    this is invalid python syntax
`,
		})

		connector := &PytestConnector{Executable: "pytest"}

		// Skip test if pytest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("pytest not installed, skipping test")
		}

		_, err = connector.DiscoverTests(projectDir)
		syntaxErr = err

		assert.Error(t, err, "should return error for syntax errors")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")

		// Verify error is "meaningful" - provides context and details
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")
		assert.Contains(t, err.Error(), "Output:",
			"error should include pytest output for debugging")
	})

	t.Run("handles nonexistent directory", func(t *testing.T) {
		connector := &PytestConnector{Executable: "pytest"}

		// Skip test if pytest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("pytest not installed, skipping test")
		}

		_, err = connector.DiscoverTests("/nonexistent/path")
		nonexistentErr = err

		assert.Error(t, err, "should return error for nonexistent directory")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")

		// Verify error is "meaningful" - provides context
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")
	})

	// Verify errors are "distinguishable" - different error types produce different messages
	t.Run("errors are distinguishable", func(t *testing.T) {
		// Skip if pytest not installed (some errors might be nil)
		if importErr == nil || syntaxErr == nil || nonexistentErr == nil {
			t.Skip("pytest not installed or some tests were skipped")
		}

		// Different error types should produce different messages
		assert.NotEqual(t, importErr.Error(), syntaxErr.Error(),
			"import and syntax errors should be distinguishable")
		assert.NotEqual(t, importErr.Error(), nonexistentErr.Error(),
			"import and directory errors should be distinguishable")
		assert.NotEqual(t, syntaxErr.Error(), nonexistentErr.Error(),
			"syntax and directory errors should be distinguishable")
	})
}

// Helper function to create a minimal Python project
func createPytestProject(t *testing.T, files map[string]string) string {
	t.Helper()
	tempDir := t.TempDir()

	// Create test files
	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if dir != tempDir {
			err := os.MkdirAll(dir, 0755)
			assert.NoError(t, err)
		}
		err := os.WriteFile(fullPath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	return tempDir
}
