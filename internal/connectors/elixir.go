package connectors

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/Alge/aligned/internal/config"
)

type ElixirConnector struct {
	Executable string
}

// NewElixirConnector creates a new ElixirConnector with the specified executable
func NewElixirConnector(executable string) *ElixirConnector {
	if executable == "" {
		executable = "mix"
	}
	return &ElixirConnector{
		Executable: executable,
	}
}

// DefaultElixirConnector returns an ElixirConnector with default configuration
func DefaultElixirConnector() *ElixirConnector {
	return &ElixirConnector{
		Executable: "mix",
	}
}

// DetectFramework checks if the mix executable is available
func (e *ElixirConnector) DetectFramework() (bool, error) {
	_, err := exec.LookPath(e.Executable)
	return err == nil, nil
}

// GenerateConfig creates a default connector configuration for Elixir
func (e *ElixirConnector) GenerateConfig(path string) config.ConnectorConfig {
	return config.ConnectorConfig{
		Type:       "elixir",
		Executable: e.Executable,
		Path:       path,
	}
}

// DiscoverTests discovers Elixir tests in the given path with a default timeout
func (e *ElixirConnector) DiscoverTests(path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return e.DiscoverTestsWithContext(ctx, path)
}

// DiscoverTestsWithContext discovers Elixir tests in the given path with a context
func (e *ElixirConnector) DiscoverTestsWithContext(ctx context.Context, path string) ([]string, error) {
	cmd := exec.CommandContext(ctx, e.Executable, "test", "--trace")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("test discovery timed out: %w", ctx.Err())
		}

		// "There are no tests to run" is not an error - it means empty test suite
		if strings.Contains(outputStr, "There are no tests to run") {
			return []string{}, nil
		}

		// Include output to help user understand the problem
		return nil, fmt.Errorf("%s test discovery failed: %w\nOutput: %s", e.Executable, err, outputStr)
	}

	return parseElixirTestOutput(outputStr), nil
}

// parseElixirTestOutput extracts test identifiers from mix test --trace output
// Format: file:Module:test name
func parseElixirTestOutput(output string) []string {
	var tests []string
	var currentFile string
	var currentModule string

	scanner := bufio.NewScanner(strings.NewReader(output))

	// Pattern for module line: "ModuleName [path/to/file.exs]"
	modulePattern := regexp.MustCompile(`^([A-Z][A-Za-z0-9.]*)\s+\[([^\]]+)\]`)

	// Pattern for test line: "  * test description [L#number]" (before timing)
	testPattern := regexp.MustCompile(`^\s+\*\s+(test\s+[^\[]+)\s+\[L#\d+\]`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a module declaration line
		if matches := modulePattern.FindStringSubmatch(line); matches != nil {
			currentModule = matches[1]
			currentFile = matches[2]
			continue
		}

		// Check if this is a test line (before timing info)
		if matches := testPattern.FindStringSubmatch(line); matches != nil && currentModule != "" && currentFile != "" {
			testName := strings.TrimSpace(matches[1])
			testID := fmt.Sprintf("%s:%s:%s", currentFile, currentModule, testName)
			tests = append(tests, testID)
		}
	}

	return tests
}
