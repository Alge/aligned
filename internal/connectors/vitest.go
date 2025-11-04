package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Alge/aligned/internal/config"
)

type VitestConnector struct {
	Executable string
}

// VitestTestItem represents a single test from vitest list --json output
type VitestTestItem struct {
	Name string `json:"name"`
	File string `json:"file"`
}

// NewVitestConnector creates a new VitestConnector with the specified executable
func NewVitestConnector(executable string) *VitestConnector {
	if executable == "" {
		executable = "vitest"
	}
	return &VitestConnector{
		Executable: executable,
	}
}

// DefaultVitestConnector returns a VitestConnector with default configuration
func DefaultVitestConnector() *VitestConnector {
	return &VitestConnector{
		Executable: "vitest",
	}
}

// DetectFramework checks if the vitest executable is available
func (v *VitestConnector) DetectFramework() (bool, error) {
	_, err := exec.LookPath(v.Executable)
	return err == nil, nil
}

// GenerateConfig creates a default connector configuration for Vitest
func (v *VitestConnector) GenerateConfig(path string) config.ConnectorConfig {
	return config.ConnectorConfig{
		Type:       "vitest",
		Executable: v.Executable,
		Path:       path,
	}
}

// DiscoverTests discovers vitest tests in the given path with a default timeout
func (v *VitestConnector) DiscoverTests(path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return v.DiscoverTestsWithContext(ctx, path)
}

// DiscoverTestsWithContext discovers vitest tests in the given path with a context
func (v *VitestConnector) DiscoverTestsWithContext(ctx context.Context, path string) ([]string, error) {
	cmd := exec.CommandContext(ctx, v.Executable, "list", "--json")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("test discovery timed out: %w", ctx.Err())
		}

		// Check if vitest is not found
		if strings.Contains(outputStr, "command not found") || strings.Contains(err.Error(), "executable file not found") {
			return nil, fmt.Errorf("%s not found in PATH. Install it with: npm install -D vitest", v.Executable)
		}

		// Include output to help user understand the problem
		return nil, fmt.Errorf("%s test discovery failed: %w\nOutput: %s", v.Executable, err, outputStr)
	}

	// Parse the JSON output
	tests, err := parseVitestJSON(outputStr, path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vitest output: %w\nOutput: %s", err, outputStr)
	}

	return tests, nil
}

// parseVitestJSON extracts test names from vitest list --json output
// Format: [{name: "...", file: "..."}]
// Returns: ["relative/path/file.test.js > test name"]
func parseVitestJSON(output string, basePath string) ([]string, error) {
	// Handle empty output
	trimmed := strings.TrimSpace(output)
	if trimmed == "" || trimmed == "[]" {
		return []string{}, nil
	}

	var items []VitestTestItem
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Handle empty test suite
	if len(items) == 0 {
		return []string{}, nil
	}

	var tests []string
	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", basePath, err)
	}

	for _, item := range items {
		// Convert absolute file path to relative path
		relPath, err := filepath.Rel(absBasePath, item.File)
		if err != nil {
			// If we can't get relative path, use the filename
			relPath = filepath.Base(item.File)
		}

		// Combine file path and test name: "path/to/file.test.js > test name"
		testIdentifier := relPath + " > " + item.Name
		tests = append(tests, testIdentifier)
	}

	return tests, nil
}
