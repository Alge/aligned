package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCreatesFile(t *testing.T) {
	// Create temp directory without config
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"init", "go-test", "./"}, &stdout, &stderr)

	// Verify config file was created
	configPath := filepath.Join(tempDir, ".align.yml")
	assert.FileExists(t, configPath, "should create .align.yml file")
}

func TestInitWritesConnectorType(t *testing.T) {
	// Create temp directory without config
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"init", "go-test", "./"}, &stdout, &stderr)

	// Verify config contains connector type
	configPath := filepath.Join(tempDir, ".align.yml")
	content, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "type: go", "should write connector type to config")
}

func TestInitWritesPath(t *testing.T) {
	// Create temp directory without config
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"init", "go-test", "./"}, &stdout, &stderr)

	// Verify config contains path
	configPath := filepath.Join(tempDir, ".align.yml")
	content, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "path: ./", "should write path to config")
}

func TestInitSuccessMessage(t *testing.T) {
	// Create temp directory without config
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	run([]string{"init", "go-test", "./"}, &stdout, &stderr)

	assert.Contains(t, strings.ToLower(stdout.String()), "created", "should display success message")
}

func TestInitExitCode(t *testing.T) {
	// Create temp directory without config
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init", "go-test", "./"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 on success")
}

func TestInitConfigExists(t *testing.T) {
	// Create temp directory with existing config
	tempDir := t.TempDir()
	configContent := `connectors:
  - type: go
    path: ./
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init", "go-test", "./"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 when config already exists")
	assert.Contains(t, strings.ToLower(stderr.String()), "already exists", "should report file already exists")
}

func TestInitMissingArguments(t *testing.T) {
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	t.Run("only connector type", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		exitCode := run([]string{"init", "go-test"}, &stdout, &stderr)

		assert.Equal(t, 1, exitCode, "should exit with code 1 when path missing")
		assert.Contains(t, strings.ToLower(stderr.String()), "usage", "should show usage message")
	})
}

func TestInitUnsupportedConnector(t *testing.T) {
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init", "unsupported", "./"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 for unsupported connector")
	assert.Contains(t, strings.ToLower(stderr.String()), "unsupported", "should report unsupported connector")
}

func TestInitNoArgsShowsHelp(t *testing.T) {
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when showing help")
	output := stdout.String()
	assert.Contains(t, strings.ToLower(output), "usage", "should display usage information")
	assert.Contains(t, strings.ToLower(output), "supported connectors", "should list supported connectors")
}

func TestInitHelpShowsHelp(t *testing.T) {
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init", "help"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when showing help")
	output := stdout.String()
	assert.Contains(t, strings.ToLower(output), "usage", "should display usage information")
	assert.Contains(t, strings.ToLower(output), "supported connectors", "should list supported connectors")
}

func TestInitListsGoConnector(t *testing.T) {
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when showing help")
	output := strings.ToLower(stdout.String())

	// Verify go-test connector is listed
	assert.Contains(t, output, "go-test", "should list go-test connector")
	assert.Contains(t, output, "go with built-in testing", "should describe go-test connector")
}

func TestInitListsPytestConnector(t *testing.T) {
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when showing help")
	output := strings.ToLower(stdout.String())

	// Verify python-pytest connector is listed
	assert.Contains(t, output, "python-pytest", "should list python-pytest connector")
	assert.Contains(t, output, "python with pytest", "should describe python-pytest connector")
}

func TestInitListsElixirConnector(t *testing.T) {
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"init"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when showing help")
	output := strings.ToLower(stdout.String())

	// Verify elixir-exunit connector is listed
	assert.Contains(t, output, "elixir-exunit", "should list elixir-exunit connector")
	assert.Contains(t, output, "elixir with exunit", "should describe elixir-exunit connector")
}
