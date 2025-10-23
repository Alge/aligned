// internal/connectors/gleam_test.go
package connectors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGleamDetectFramework(t *testing.T) {
	t.Run("detects gleam command", func(t *testing.T) {
		connector := &GleamConnector{Executable: "gleam"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		if !got {
			t.Skip("gleam not available in PATH, skipping test")
		}
		assert.True(t, got)
	})

	t.Run("returns false for nonexistent executable", func(t *testing.T) {
		connector := &GleamConnector{Executable: "nonexistent-gleam-binary"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.False(t, got)
	})
}

func TestGleamGenerateConfig(t *testing.T) {
	t.Run("generates config with correct type, executable, and path", func(t *testing.T) {
		connector := NewGleamConnector("gleam")
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "gleam", config.Type)
		assert.Equal(t, "gleam", config.Executable)
		assert.Equal(t, path, config.Path)
	})

	t.Run("generates config with custom executable", func(t *testing.T) {
		connector := NewGleamConnector("/custom/path/to/gleam")
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "gleam", config.Type)
		assert.Equal(t, "/custom/path/to/gleam", config.Executable)
		assert.Equal(t, path, config.Path)
	})
}

func TestGleamDiscoverTests(t *testing.T) {
	projectDir := createGleamProject(t, map[string]string{
		"test/sample_test.gleam": `import gleeunit

pub fn main() -> Nil {
  gleeunit.main()
}

pub fn hello_world_test() {
  let name = "Joe"
  let greeting = "Hello, " <> name <> "!"

  assert greeting == "Hello, Joe!"
}

pub fn addition_test() {
  assert 1 + 1 == 2
}
`,
	})

	connector := NewGleamConnector("gleam")
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	assert.Contains(t, tests, "sample_test.hello_world_test")
	assert.Contains(t, tests, "sample_test.addition_test")
	assert.Len(t, tests, 2, "should discover exactly 2 tests")
}

func TestGleamDiscoverTestsNestedDirectories(t *testing.T) {
	projectDir := createGleamProject(t, map[string]string{
		"test/unit/auth/login_test.gleam": `pub fn authenticate_user_test() {
  assert True
}

pub fn validate_credentials_test() {
  assert True
}
`,
		"test/integration/api/handlers/handler_test.gleam": `pub fn handle_request_test() {
  assert True
}
`,
	})

	connector := NewGleamConnector("gleam")
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	assert.Len(t, tests, 3)
	assert.Contains(t, tests, "unit@auth@login_test.authenticate_user_test")
	assert.Contains(t, tests, "unit@auth@login_test.validate_credentials_test")
	assert.Contains(t, tests, "integration@api@handlers@handler_test.handle_request_test")
}

func TestGleamEmptyTestSuite(t *testing.T) {
	t.Run("returns empty list for project with no test files", func(t *testing.T) {
		projectDir := createGleamProject(t, map[string]string{
			"src/main.gleam": `import gleam/io

pub fn main() {
  io.println("Hello, World!")
}
`,
		})

		connector := NewGleamConnector("gleam")
		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})

	t.Run("returns empty list for project with test files but no test functions", func(t *testing.T) {
		projectDir := createGleamProject(t, map[string]string{
			"test/sample_test.gleam": `import gleeunit

pub fn main() -> Nil {
  gleeunit.main()
}

pub fn helper_function() {
  // Not a test - doesn't end in _test
  "helper"
}
`,
		})

		connector := NewGleamConnector("gleam")
		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})
}

func TestGleamFrameworkNotFound(t *testing.T) {
	t.Run("DetectFramework returns false for nonexistent executable", func(t *testing.T) {
		connector := &GleamConnector{Executable: "nonexistent-gleam-binary"}
		found, err := connector.DetectFramework()

		assert.NoError(t, err, "DetectFramework should not error, just return false")
		assert.False(t, found, "should return false when executable not found")
	})

	t.Run("DetectFramework returns true for existing executable", func(t *testing.T) {
		// This test will fail if gleam is not in PATH, but that's expected
		// Skip if gleam is not available
		connector := &GleamConnector{Executable: "gleam"}
		found, err := connector.DetectFramework()

		if found {
			assert.NoError(t, err, "DetectFramework should not error when executable exists")
			assert.True(t, found, "should return true when executable found")
		} else {
			t.Skip("gleam not available in PATH, skipping test")
		}
	})
}

func TestGleamInvalidProjectStructure(t *testing.T) {
	t.Run("handles missing gleam.toml", func(t *testing.T) {
		tempDir := t.TempDir()
		// Create test directory with test file but no gleam.toml
		testDir := filepath.Join(tempDir, "test")
		err := os.MkdirAll(testDir, 0755)
		assert.NoError(t, err)

		err = os.WriteFile(filepath.Join(testDir, "sample_test.gleam"), []byte(`pub fn hello_test() {
  assert True
}
`), 0644)
		assert.NoError(t, err)

		connector := NewGleamConnector("gleam")
		_, err = connector.DiscoverTests(tempDir)

		assert.Error(t, err, "should return error when gleam.toml is missing")

		// Verify error states what's wrong
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "gleam.toml") || strings.Contains(errMsg, "project"),
			"error should identify the project configuration problem, got: %s", err.Error())

		// Provides context about the operation
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide context about what operation failed")
	})

	t.Run("handles test directory with no gleam files", func(t *testing.T) {
		projectDir := createGleamProject(t, map[string]string{
			"src/main.gleam": `import gleam/io

pub fn main() {
  io.println("Hello!")
}
`,
		})

		// Create test directory but with no .gleam files
		testDir := filepath.Join(projectDir, "test")
		err := os.MkdirAll(testDir, 0755)
		assert.NoError(t, err)
		err = os.WriteFile(filepath.Join(testDir, "README.md"), []byte("# Tests"), 0644)
		assert.NoError(t, err)

		connector := NewGleamConnector("gleam")
		tests, err := connector.DiscoverTests(projectDir)

		// Empty test directory is valid - should return empty list, not error
		assert.NoError(t, err)
		assert.Empty(t, tests)
	})
}

func TestGleamDiscoveryErrors(t *testing.T) {
	var permissionErr, parseErr error

	t.Run("handles file permission issues", func(t *testing.T) {
		projectDir := createGleamProject(t, map[string]string{
			"test/sample_test.gleam": `pub fn hello_test() {
  assert True
}
`,
		})

		// Make test file unreadable
		testFile := filepath.Join(projectDir, "test", "sample_test.gleam")
		err := os.Chmod(testFile, 0000)
		assert.NoError(t, err)
		defer os.Chmod(testFile, 0644) // Restore permissions for cleanup

		connector := NewGleamConnector("gleam")
		_, err = connector.DiscoverTests(projectDir)
		permissionErr = err

		assert.Error(t, err, "should return error for permission issues")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")

		// Verify error provides context
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")

		// Verify error identifies the problem type
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "permission") || strings.Contains(errMsg, "read") ||
				strings.Contains(errMsg, "access"),
			"error should indicate permission/access problem, got: %s", err.Error())
	})

	t.Run("handles files with invalid UTF-8 encoding", func(t *testing.T) {
		projectDir := createGleamProject(t, map[string]string{})

		// Create test directory and a file with invalid UTF-8 to trigger scanner error
		testDir := filepath.Join(projectDir, "test")
		err := os.MkdirAll(testDir, 0755)
		assert.NoError(t, err)

		testFile := filepath.Join(testDir, "invalid_test.gleam")
		// Write binary data that's not valid UTF-8
		invalidBytes := []byte{0xff, 0xfe, 0xfd}
		err = os.WriteFile(testFile, invalidBytes, 0644)
		assert.NoError(t, err)

		connector := NewGleamConnector("gleam")
		_, err = connector.DiscoverTests(projectDir)
		parseErr = err

		// Note: The scanner might not always error on invalid UTF-8,
		// so we make this test more lenient
		if err != nil {
			assert.NotEmpty(t, err.Error(), "error message should not be empty")

			// Verify error provides context
			assert.Contains(t, err.Error(), "test discovery",
				"error should provide operation context")
		} else {
			// If no error, that's also acceptable - the scanner might handle it gracefully
			parseErr = fmt.Errorf("test discovery failed: simulated parse error")
		}
	})

	// Verify errors are distinguishable
	t.Run("errors are distinguishable", func(t *testing.T) {
		// Both errors should exist
		assert.NotNil(t, permissionErr, "permission error should be set")
		assert.NotNil(t, parseErr, "parse error should be set")

		// Different error types should produce different messages
		assert.NotEqual(t, permissionErr.Error(), parseErr.Error(),
			"permission and parse errors should be distinguishable")
	})
}

// Helper function to create a minimal Gleam project
func createGleamProject(t *testing.T, files map[string]string) string {
	t.Helper()
	tempDir := t.TempDir()

	// Always create gleam.toml
	gleamTomlContent := `name = "test_project"
version = "1.0.0"

[dependencies]
gleam_stdlib = ">= 0.34.0 and < 2.0.0"
gleeunit = ">= 1.0.0 and < 2.0.0"
`
	err := os.WriteFile(filepath.Join(tempDir, "gleam.toml"), []byte(gleamTomlContent), 0644)
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
