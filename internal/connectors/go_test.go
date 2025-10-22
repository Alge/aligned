// internal/connectors/go_test.go
package connectors

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestGoDetectFramework(t *testing.T) {
	t.Run("detects go command", func(t *testing.T) {
		connector := &GoConnector{Executable: "go"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("returns false for nonexistent executable", func(t *testing.T) {
		connector := &GoConnector{Executable: "nonexistent-go-binary"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.False(t, got)
	})
}

func TestGoValidateCompatibility(t *testing.T) {
	t.Run("validates go 1.13 or higher", func(t *testing.T) {
		connector := &GoConnector{Executable: "go"}
		compatible, err := connector.ValidateCompatibility()

		assert.NoError(t, err)
		assert.True(t, compatible, "current go version should be compatible (1.13+)")
	})
}

func TestGoValidateConfiguration(t *testing.T) {
	t.Run("validates valid configuration", func(t *testing.T) {
		connector := &GoConnector{Executable: "go"}
		err := connector.ValidateConfiguration()

		assert.NoError(t, err)
	})

	t.Run("rejects empty executable", func(t *testing.T) {
		connector := &GoConnector{Executable: ""}
		err := connector.ValidateConfiguration()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "executable")
	})
}

func TestGoDefaultConfiguration(t *testing.T) {
	t.Run("returns default configuration", func(t *testing.T) {
		connector := DefaultGoConnector()

		assert.Equal(t, "go", connector.Executable)
	})
}

func TestGoGenerateConfig(t *testing.T) {
	t.Run("generates config with correct type, executable, and path", func(t *testing.T) {
		connector := NewGoConnector("go")
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "go", config.Type)
		assert.Equal(t, "go", config.Executable)
		assert.Equal(t, path, config.Path)
	})

	t.Run("generates config with custom executable", func(t *testing.T) {
		connector := NewGoConnector("/custom/path/to/go")
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "go", config.Type)
		assert.Equal(t, "/custom/path/to/go", config.Executable)
		assert.Equal(t, path, config.Path)
	})
}

func TestGoDiscoveryTimeout(t *testing.T) {
	t.Run("respects context timeout", func(t *testing.T) {
		projectDir := createGoProject(t, map[string]string{
			"example_test.go": `package example
import "testing"
func TestFoo(t *testing.T) {}
`,
		})

		connector := NewGoConnector("go")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tests, err := connector.DiscoverTestsWithContext(ctx, projectDir)

		assert.NoError(t, err)
		assert.Contains(t, tests, "testproject.TestFoo")
	})

	t.Run("returns error when context expires", func(t *testing.T) {
		projectDir := createGoProject(t, map[string]string{
			"example_test.go": `package example
import "testing"
func TestFoo(t *testing.T) {}
`,
		})

		connector := NewGoConnector("go")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Millisecond) // Ensure context expires

		_, err := connector.DiscoverTestsWithContext(ctx, projectDir)

		assert.Error(t, err)
	})
}

func TestGoDiscoveryErrors(t *testing.T) {
	// This test verifies errors are meaningful and distinguishable
	var directoryErr, moduleErr, compilationErr error

	t.Run("handles nonexistent directory", func(t *testing.T) {
		connector := NewGoConnector("go")
		_, err := connector.DiscoverTests("/nonexistent/path")
		directoryErr = err

		assert.Error(t, err, "should return error for nonexistent directory")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")

		// Verify error is "meaningful" - provides context
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")

		// Verify error identifies the problem type
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "directory") || strings.Contains(errMsg, "path") || strings.Contains(errMsg, "no such"),
			"error should indicate directory/path problem, got: %s", err.Error())
	})

	t.Run("handles invalid go.mod", func(t *testing.T) {
		tempDir := t.TempDir()
		invalidGoMod := `this is not valid go.mod syntax{{{`
		err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(invalidGoMod), 0644)
		assert.NoError(t, err)

		connector := NewGoConnector("go")
		_, err = connector.DiscoverTests(tempDir)
		moduleErr = err

		assert.Error(t, err, "should return error for invalid go.mod")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")

		// Verify error is "meaningful" - provides context
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")

		// Verify error identifies the problem type
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "module") || strings.Contains(errMsg, "go.mod") || strings.Contains(errMsg, "parse"),
			"error should indicate module/parsing problem, got: %s", err.Error())
	})

	t.Run("handles compilation errors", func(t *testing.T) {
		projectDir := createGoProject(t, map[string]string{
			"example_test.go": `package example
import "testing"
func TestFoo(t *testing.T) {
	this is invalid go code
}
`,
		})

		connector := NewGoConnector("go")
		_, err := connector.DiscoverTests(projectDir)
		compilationErr = err

		assert.Error(t, err, "should return error for compilation failure")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")

		// Verify error is "meaningful" - provides context and details
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")
		assert.Contains(t, err.Error(), "Output:",
			"error should include compiler output for debugging")

		// Verify error identifies the problem type
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "syntax") || strings.Contains(errMsg, "compile") ||
			strings.Contains(errMsg, "build") || strings.Contains(errMsg, "undefined") ||
			strings.Contains(errMsg, "expected") || strings.Contains(errMsg, "setup failed"),
			"error should indicate compilation problem, got: %s", err.Error())
	})

	// Verify errors are "distinguishable" - different error types produce different messages
	t.Run("errors are distinguishable", func(t *testing.T) {
		// All errors should exist
		assert.NotNil(t, directoryErr, "directory error should be set")
		assert.NotNil(t, moduleErr, "module error should be set")
		assert.NotNil(t, compilationErr, "compilation error should be set")

		// Different error types should produce different messages
		assert.NotEqual(t, directoryErr.Error(), moduleErr.Error(),
			"directory and module errors should be distinguishable")
		assert.NotEqual(t, directoryErr.Error(), compilationErr.Error(),
			"directory and compilation errors should be distinguishable")
		assert.NotEqual(t, moduleErr.Error(), compilationErr.Error(),
			"module and compilation errors should be distinguishable")
	})
}

func TestGoFrameworkNotFound(t *testing.T) {
	projectDir := createGoProject(t, map[string]string{
		"example_test.go": `package example
import "testing"
func TestFoo(t *testing.T) {}
`,
	})

	connector := &GoConnector{Executable: "nonexistent-go-binary"}
	_, err := connector.DiscoverTests(projectDir)

	assert.Error(t, err, "should return error when go command not found")

	// Verify error identifies the specific problematic value
	assert.Contains(t, err.Error(), "nonexistent-go-binary",
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

func TestGoInvalidProjectStructure(t *testing.T) {
	t.Run("handles missing go.mod", func(t *testing.T) {
		tempDir := t.TempDir()
		// Create a test file without go.mod
		err := os.WriteFile(filepath.Join(tempDir, "example_test.go"), []byte(`package example
import "testing"
func TestFoo(t *testing.T) {}
`), 0644)
		assert.NoError(t, err)

		connector := NewGoConnector("go")
		_, err = connector.DiscoverTests(tempDir)

		assert.Error(t, err, "should return error when go.mod is missing")

		// Verify error states what's wrong - should mention module problem
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "go.mod") || strings.Contains(errMsg, "module"),
			"error should identify the module/go.mod problem, got: %s", err.Error())

		// 2. Provides context about the operation
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide context about what operation failed")

		// 3. Includes helpful details from go test output
		assert.Contains(t, err.Error(), "Output:",
			"error should include go test output for debugging")
	})

	t.Run("handles directory without go files", func(t *testing.T) {
		tempDir := t.TempDir()

		connector := NewGoConnector("go")
		_, err := connector.DiscoverTests(tempDir)

		assert.Error(t, err, "should return error for empty directory")

		// Error should not be just "exit status 1"
		assert.NotContains(t, err.Error(), "exit status 1\n\n",
			"error should include helpful context, not just exit status")

		// Should provide context
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide context about what operation failed")
	})
}

func TestGoEmptyTestSuite(t *testing.T) {
	t.Run("returns empty list for project with no tests", func(t *testing.T) {
		projectDir := createGoProject(t, map[string]string{
			"main.go": `package main

func main() {
	println("Hello, World!")
}
`,
		})

		connector := NewGoConnector("go")
		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})

	t.Run("returns empty list for project with only non-test go files", func(t *testing.T) {
		projectDir := createGoProject(t, map[string]string{
			"utils.go": `package utils

func Add(a, b int) int {
	return a + b
}
`,
		})

		connector := NewGoConnector("go")
		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})
}

func TestGoPartialFailures(t *testing.T) {
	t.Run("handles mixed valid and invalid packages", func(t *testing.T) {
		projectDir := createGoProject(t, map[string]string{
			"good/good_test.go": `package good
import "testing"
func TestGood(t *testing.T) {}
`,
			"bad/bad_test.go": `package bad
import "testing"
func TestBad(t *testing.T) {
	this is invalid go code
}
`,
		})

		connector := NewGoConnector("go")
		_, err := connector.DiscoverTests(projectDir)

		// Current behavior: fails when any package fails
		// This test documents the current behavior
		assert.Error(t, err)
	})

	t.Run("handles package with valid and invalid test files", func(t *testing.T) {
		projectDir := createGoProject(t, map[string]string{
			"example_test.go": `package example
import "testing"
func TestValid(t *testing.T) {}
`,
			"bad_test.go": `package example
import "testing"
func TestInvalid(t *testing.T) {
	invalid syntax here
}
`,
		})

		connector := NewGoConnector("go")
		_, err := connector.DiscoverTests(projectDir)

		// Current behavior: fails when compilation fails
		assert.Error(t, err)
	})
}

// Helper function to create a minimal Go project
func createGoProject(t *testing.T, files map[string]string) string {
	t.Helper()
	tempDir := t.TempDir()
	
	// Always create go.mod
	goModContent := `module testproject

go 1.23
`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	assert.NoError(t, err)
	
	// Create test files
	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if dir != tempDir {
			err = os.MkdirAll(dir, 0755)
			assert.NoError(t, err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		assert.NoError(t, err)
	}
	
	return tempDir
}

func TestGoDiscoverTests(t *testing.T) {
	// This test verifies the connector uses `go test -list=. ./...` by testing the behavior
	// that would only work with those specific flags:
	// - Tests are listed, not run (requires -list flag)
	// - All tests matching pattern are found (requires . pattern)
	// - Tests in all packages are discovered (requires ./... recursive flag - see TestGoDiscoverTestsNestedDirectories)

	projectDir := createGoProject(t, map[string]string{
		"example_test.go": `package example
import "testing"
func TestFoo(t *testing.T) {}
func TestBar(t *testing.T) {}
`,
	})

	connector := NewGoConnector("go")
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	// Verify both tests are discovered (validates -list=. pattern matches all tests)
	assert.Contains(t, tests, "testproject.TestFoo")
	assert.Contains(t, tests, "testproject.TestBar")
	// Verify tests were listed, not executed (would timeout or fail if -list not used)
	assert.Len(t, tests, 2, "should discover exactly 2 tests")
}

func TestGoDiscoverTestsNestedDirectories(t *testing.T) {
	// This test verifies the ./... recursive flag is used by ensuring tests in deeply
	// nested packages are discovered. Without ./..., only tests in the current package
	// would be found.

	projectDir := createGoProject(t, map[string]string{
		"internal/auth/auth_test.go": `package auth
import "testing"
func TestLogin(t *testing.T) {}
`,
		"cmd/server/handlers/api/api_test.go": `package api
import "testing"
func TestAPIHandler(t *testing.T) {}
`,
	})

	connector := NewGoConnector("go")
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	assert.Len(t, tests, 2)
	// Verify nested packages are discovered (requires ./... recursive pattern)
	assert.Contains(t, tests, "internal/auth.TestLogin")
	assert.Contains(t, tests, "cmd/server/handlers/api.TestAPIHandler")
}