// internal/connectors/elixir_test.go
package connectors

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElixirDetectFramework(t *testing.T) {
	t.Run("detects mix command", func(t *testing.T) {
		connector := &ElixirConnector{Executable: "mix"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("returns false for nonexistent executable", func(t *testing.T) {
		connector := &ElixirConnector{Executable: "nonexistent-mix-binary"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.False(t, got)
	})
}

func TestElixirGenerateConfig(t *testing.T) {
	t.Run("generates config with correct type, executable, and path", func(t *testing.T) {
		connector := &ElixirConnector{Executable: "mix"}
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "elixir", config.Type)
		assert.Equal(t, "mix", config.Executable)
		assert.Equal(t, path, config.Path)
	})

	t.Run("generates config with custom executable", func(t *testing.T) {
		connector := &ElixirConnector{Executable: "/custom/path/to/mix"}
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "elixir", config.Type)
		assert.Equal(t, "/custom/path/to/mix", config.Executable)
		assert.Equal(t, path, config.Path)
	})
}

func TestElixirDiscoverTests(t *testing.T) {
	projectDir := createElixirProject(t, map[string]string{
		"test/example_test.exs": `defmodule ExampleTest do
  use ExUnit.Case

  test "foo" do
    assert true
  end

  test "bar" do
    assert true
  end
end
`,
	})

	connector := &ElixirConnector{Executable: "mix"}

	// Skip test if mix is not installed
	detected, err := connector.DetectFramework()
	if err != nil || !detected {
		t.Skip("mix not installed, skipping test")
	}

	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	assert.Contains(t, tests, "test/example_test.exs:ExampleTest:test foo")
	assert.Contains(t, tests, "test/example_test.exs:ExampleTest:test bar")
	assert.Len(t, tests, 2, "should discover exactly 2 tests")
}

func TestElixirDiscoverTestsNestedDirectories(t *testing.T) {
	projectDir := createElixirProject(t, map[string]string{
		"test/unit/auth/login_test.exs": `defmodule Unit.Auth.LoginTest do
  use ExUnit.Case

  test "validates credentials" do
    assert true
  end
end
`,
		"test/integration/api/handler_test.exs": `defmodule Integration.API.HandlerTest do
  use ExUnit.Case

  test "handles request" do
    assert true
  end
end
`,
	})

	connector := &ElixirConnector{Executable: "mix"}

	// Skip test if mix is not installed
	detected, err := connector.DetectFramework()
	if err != nil || !detected {
		t.Skip("mix not installed, skipping test")
	}

	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	assert.Len(t, tests, 2)
	assert.Contains(t, tests, "test/unit/auth/login_test.exs:Unit.Auth.LoginTest:test validates credentials")
	assert.Contains(t, tests, "test/integration/api/handler_test.exs:Integration.API.HandlerTest:test handles request")
}

func TestElixirEmptyTestSuite(t *testing.T) {
	t.Run("returns empty list for project with no tests", func(t *testing.T) {
		projectDir := createElixirProject(t, map[string]string{})

		connector := &ElixirConnector{Executable: "mix"}

		// Skip test if mix is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("mix not installed, skipping test")
		}

		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})
}

func TestElixirFrameworkNotFound(t *testing.T) {
	projectDir := createElixirProject(t, map[string]string{
		"test/example_test.exs": `defmodule ExampleTest do
  use ExUnit.Case

  test "foo" do
    assert true
  end
end
`,
	})

	connector := &ElixirConnector{Executable: "nonexistent-mix-binary"}
	_, err := connector.DiscoverTests(projectDir)

	assert.Error(t, err, "should return error when mix command not found")

	// Verify error is "clear" per specification
	assert.Contains(t, err.Error(), "nonexistent-mix-binary",
		"error should identify the executable that was not found")

	errMsg := strings.ToLower(err.Error())
	assert.True(t,
		strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "no such"),
		"error should clearly state the problem, got: %s", err.Error())

	assert.Contains(t, err.Error(), "test discovery",
		"error should provide context about what operation failed")
}

func TestElixirInvalidProjectStructure(t *testing.T) {
	t.Run("handles missing mix.exs", func(t *testing.T) {
		tempDir := t.TempDir()
		// Create test file without mix.exs
		testDir := filepath.Join(tempDir, "test")
		err := os.MkdirAll(testDir, 0755)
		assert.NoError(t, err)
		err = os.WriteFile(filepath.Join(testDir, "example_test.exs"), []byte(`defmodule ExampleTest do
  use ExUnit.Case
  test "foo" do
    assert true
  end
end
`), 0644)
		assert.NoError(t, err)

		connector := &ElixirConnector{Executable: "mix"}

		// Skip test if mix is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("mix not installed, skipping test")
		}

		_, err = connector.DiscoverTests(tempDir)

		assert.Error(t, err, "should return error when mix.exs is missing")
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide context about what operation failed")
	})
}

func TestElixirDiscoveryErrors(t *testing.T) {
	var syntaxErr, nonexistentErr error

	t.Run("handles syntax errors", func(t *testing.T) {
		projectDir := createElixirProject(t, map[string]string{
			"test/broken_test.exs": `defmodule BrokenTest do
  use ExUnit.Case

  test "syntax error" do
    this is not valid elixir code
  end
end
`,
		})

		connector := &ElixirConnector{Executable: "mix"}

		// Skip test if mix is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("mix not installed, skipping test")
		}

		_, err = connector.DiscoverTests(projectDir)
		syntaxErr = err

		assert.Error(t, err, "should return error for syntax errors")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")
		assert.Contains(t, err.Error(), "Output:",
			"error should include mix output for debugging")
	})

	t.Run("handles nonexistent directory", func(t *testing.T) {
		connector := &ElixirConnector{Executable: "mix"}

		// Skip test if mix is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("mix not installed, skipping test")
		}

		_, err = connector.DiscoverTests("/nonexistent/path")
		nonexistentErr = err

		assert.Error(t, err, "should return error for nonexistent directory")
		assert.NotEmpty(t, err.Error(), "error message should not be empty")
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide operation context")
	})

	t.Run("errors are distinguishable", func(t *testing.T) {
		if syntaxErr == nil || nonexistentErr == nil {
			t.Skip("mix not installed or some tests were skipped")
		}

		assert.NotEqual(t, syntaxErr.Error(), nonexistentErr.Error(),
			"syntax and directory errors should be distinguishable")
	})
}

// Helper function to create a minimal Elixir project
func createElixirProject(t *testing.T, files map[string]string) string {
	t.Helper()
	tempDir := t.TempDir()

	// Always create mix.exs
	mixExsContent := `defmodule TestProject.MixProject do
  use Mix.Project

  def project do
    [
      app: :test_project,
      version: "0.1.0",
      elixir: "~> 1.14",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  def application do
    [
      extra_applications: [:logger]
    ]
  end

  defp deps do
    []
  end
end
`
	err := os.WriteFile(filepath.Join(tempDir, "mix.exs"), []byte(mixExsContent), 0644)
	assert.NoError(t, err)

	// Create test_helper.exs
	testHelperContent := `ExUnit.start()`
	testDir := filepath.Join(tempDir, "test")
	err = os.MkdirAll(testDir, 0755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(testDir, "test_helper.exs"), []byte(testHelperContent), 0644)
	assert.NoError(t, err)

	// Create test files
	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if dir != tempDir && dir != testDir {
			err = os.MkdirAll(dir, 0755)
			assert.NoError(t, err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	return tempDir
}
