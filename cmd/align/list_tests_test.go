package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTestsDiscovery(t *testing.T) {
	// This test verifies that connectors from config are actually used for discovery
	// by ensuring that tests are only found when the correct connector type is configured

	t.Run("uses configured go connector to discover tests", func(t *testing.T) {
		tempDir := t.TempDir()
		configContent := `connectors:
  - type: go
    executable: go
    path: .
`
		err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
		assert.NoError(t, err)

		// Create a Go test file
		testContent := `package example
import "testing"
func TestExample(t *testing.T) {}
func TestAnother(t *testing.T) {}
`
		err = os.WriteFile(filepath.Join(tempDir, "example_test.go"), []byte(testContent), 0644)
		assert.NoError(t, err)

		goModContent := `module testproject
go 1.23
`
		err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
		assert.NoError(t, err)

		// Change to temp directory
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		var stdout, stderr bytes.Buffer
		run([]string{"list-tests"}, &stdout, &stderr)

		output := stdout.String()
		// Verify tests were discovered using the Go connector
		assert.NotEmpty(t, output, "should discover tests using configured connector")
		assert.Contains(t, output, "TestExample", "should find Go tests")
		assert.Contains(t, output, "TestAnother", "should find Go tests")
	})

	t.Run("respects connector path configuration", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create subdirectory with tests
		subDir := filepath.Join(tempDir, "subproject")
		err := os.Mkdir(subDir, 0755)
		assert.NoError(t, err)

		// Config points to subdirectory
		configContent := `connectors:
  - type: go
    executable: go
    path: ./subproject
`
		err = os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
		assert.NoError(t, err)

		// Create test in subdirectory
		testContent := `package subexample
import "testing"
func TestSubproject(t *testing.T) {}
`
		err = os.WriteFile(filepath.Join(subDir, "sub_test.go"), []byte(testContent), 0644)
		assert.NoError(t, err)

		goModContent := `module subproject
go 1.23
`
		err = os.WriteFile(filepath.Join(subDir, "go.mod"), []byte(goModContent), 0644)
		assert.NoError(t, err)

		// Create test in root directory (should NOT be discovered)
		rootTestContent := `package root
import "testing"
func TestRoot(t *testing.T) {}
`
		err = os.WriteFile(filepath.Join(tempDir, "root_test.go"), []byte(rootTestContent), 0644)
		assert.NoError(t, err)

		rootGoModContent := `module root
go 1.23
`
		err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(rootGoModContent), 0644)
		assert.NoError(t, err)

		// Change to temp directory
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		var stdout, stderr bytes.Buffer
		run([]string{"list-tests"}, &stdout, &stderr)

		output := stdout.String()
		// Verify only subdirectory tests are found (path config is respected)
		assert.Contains(t, output, "TestSubproject",
			"should find tests in configured path")
		assert.NotContains(t, output, "TestRoot",
			"should not find tests outside configured path - proves connector path config is used")
	})
}

func TestListTestsPrintsNames(t *testing.T) {
	// Create temp directory with valid config
	tempDir := t.TempDir()
	configContent := `connectors:
  - type: go
    executable: go
    path: .
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Create a simple test file
	testContent := `package example
import "testing"
func TestExample(t *testing.T) {}
func TestAnother(t *testing.T) {}
`
	err = os.WriteFile(filepath.Join(tempDir, "example_test.go"), []byte(testContent), 0644)
	assert.NoError(t, err)

	// Create go.mod
	goModContent := `module testproject
go 1.23
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"list-tests"}, &stdout, &stderr)

	output := stdout.String()
	assert.Contains(t, output, "TestExample", "output should include discovered test names")
	assert.Contains(t, output, "TestAnother", "output should include discovered test names")
}

func TestListTestsExitCode(t *testing.T) {
	// Create temp directory with valid config
	tempDir := t.TempDir()
	configContent := `connectors:
  - type: go
    executable: go
    path: .
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Create a simple test file
	testContent := `package example
import "testing"
func TestExample(t *testing.T) {}
`
	err = os.WriteFile(filepath.Join(tempDir, "example_test.go"), []byte(testContent), 0644)
	assert.NoError(t, err)

	// Create go.mod
	goModContent := `module testproject
go 1.23
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"list-tests"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 on success")
}

func TestListTestsConfigMissing(t *testing.T) {
	// Create temp directory without config
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"list-tests"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 when config missing")
	assert.Contains(t, strings.ToLower(stderr.String()), "not found")
}

func TestListTestsConfigInvalid(t *testing.T) {
	// Create temp directory with invalid config
	tempDir := t.TempDir()

	// Create malformed YAML config
	invalidConfigContent := `connectors:
  - type: go
    executable: go
    path .   # Missing colon - invalid YAML
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(invalidConfigContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"list-tests"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 when config invalid")
}
func TestGoConnectorRegisteredInListTests(t *testing.T) {
	tempDir := t.TempDir()
	configContent := "connectors:\n  - type: go\n    path: .\n"
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"list-tests"}, &stdout, &stderr)

	stderrStr := stderr.String()
	assert.NotContains(t, strings.ToLower(stderrStr), "unsupported connector type",
		"go connector should be registered in list-tests command")
}

func TestPytestConnectorRegisteredInListTests(t *testing.T) {
	tempDir := t.TempDir()
	configContent := "connectors:\n  - type: pytest\n    path: .\n"
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"list-tests"}, &stdout, &stderr)

	stderrStr := stderr.String()
	assert.NotContains(t, strings.ToLower(stderrStr), "unsupported connector type",
		"pytest connector should be registered in list-tests command")
}

func TestElixirConnectorRegisteredInListTests(t *testing.T) {
	tempDir := t.TempDir()
	configContent := "connectors:\n  - type: elixir\n    path: .\n"
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"list-tests"}, &stdout, &stderr)

	stderrStr := stderr.String()
	assert.NotContains(t, strings.ToLower(stderrStr), "unsupported connector type",
		"elixir connector should be registered in list-tests command")
}

func TestGleamConnectorRegisteredInListTests(t *testing.T) {
	tempDir := t.TempDir()
	configContent := "connectors:\n  - type: gleam\n    path: .\n"
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"list-tests"}, &stdout, &stderr)

	stderrStr := stderr.String()
	assert.NotContains(t, strings.ToLower(stderrStr), "unsupported connector type",
		"gleam connector should be registered in list-tests command")
}

func TestVitestConnectorRegisteredInListTests(t *testing.T) {
	tempDir := t.TempDir()
	configContent := "connectors:\n  - type: vitest\n    path: .\n"
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"list-tests"}, &stdout, &stderr)

	stderrStr := stderr.String()
	assert.NotContains(t, strings.ToLower(stderrStr), "unsupported connector type",
		"vitest connector should be registered in list-tests command")
}
