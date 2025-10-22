package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckconfSuccessMessage(t *testing.T) {
	// Create temp directory with valid config
	tempDir := t.TempDir()
	configContent := `connectors:
  - type: go
    executable: go
    path: ./
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"checkconf"}, &stdout, &stderr)

	assert.Contains(t, strings.ToLower(stdout.String()), "valid", "should display success message")
}

func TestCheckconfExitCode(t *testing.T) {
	// Create temp directory with valid config
	tempDir := t.TempDir()
	configContent := `connectors:
  - type: go
    executable: go
    path: ./
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"checkconf"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when config is valid")
}

func TestCheckconfVerbosePrintsDetails(t *testing.T) {
	// Create temp directory with valid config
	tempDir := t.TempDir()
	configContent := `connectors:
  - type: go
    executable: go
    path: ./
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"checkconf", "-v"}, &stdout, &stderr)

	output := stdout.String()
	assert.Contains(t, output, "type: go", "should print configuration details")
	assert.Contains(t, output, "executable: go", "should print configuration details")
	assert.Contains(t, output, "path: ./", "should print configuration details")
}

func TestCheckconfVerboseExitCode(t *testing.T) {
	// Create temp directory with valid config
	tempDir := t.TempDir()
	configContent := `connectors:
  - type: go
    executable: go
    path: ./
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"checkconf", "-v"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when using verbose flag")
}

func TestCheckconfMissingFile(t *testing.T) {
	// Create temp directory WITHOUT config file
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"checkconf"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 when config file missing")
	assert.Contains(t, strings.ToLower(stderr.String()), "not found", "should report file not found")
}

func TestCheckconfInvalidYAML(t *testing.T) {
	// Create temp directory with invalid YAML
	tempDir := t.TempDir()
	invalidContent := `connectors:
  - type: go
    executable: go
    path: ./
  invalid yaml here: [[[
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(invalidContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"checkconf"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 when YAML is invalid")
	assert.Contains(t, strings.ToLower(stderr.String()), "invalid", "should report invalid YAML")
}

func TestCheckconfEmptyConfig(t *testing.T) {
	// Create temp directory with empty config
	tempDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(""), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"checkconf"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 when config has no connectors")
	assert.Contains(t, strings.ToLower(stderr.String()), "no connectors", "should report no connectors")
}
