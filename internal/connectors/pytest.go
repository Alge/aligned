package connectors

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Alge/aligned/internal/config"
)

type PytestConnector struct {
	Executable string
}

// NewPytestConnector creates a new PytestConnector with the specified executable
func NewPytestConnector(executable string) *PytestConnector {
	if executable == "" {
		executable = "pytest"
	}
	return &PytestConnector{
		Executable: executable,
	}
}

// DefaultPytestConnector returns a PytestConnector with default configuration
func DefaultPytestConnector() *PytestConnector {
	return &PytestConnector{
		Executable: "pytest",
	}
}

// DetectFramework checks if the pytest executable is available
func (p *PytestConnector) DetectFramework() (bool, error) {
	_, err := exec.LookPath(p.Executable)
	return err == nil, nil
}

// GenerateConfig creates a default connector configuration for Pytest
func (p *PytestConnector) GenerateConfig(path string) config.ConnectorConfig {
	return config.ConnectorConfig{
		Type:       "pytest",
		Executable: p.Executable,
		Path:       path,
	}
}

// DiscoverTests discovers pytest tests in the given path with a default timeout
func (p *PytestConnector) DiscoverTests(path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return p.DiscoverTestsWithContext(ctx, path)
}

// DiscoverTestsWithContext discovers pytest tests in the given path with a context
func (p *PytestConnector) DiscoverTestsWithContext(ctx context.Context, path string) ([]string, error) {
	cmd := exec.CommandContext(ctx, p.Executable, "--collect-only", "-q")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("test discovery timed out: %w", ctx.Err())
		}

		// Exit status 5 with "no tests collected" and NO errors is not an error - it means empty test suite
		// Check for error indicators like "ERROR" or "error during collection"
		lowerOutput := strings.ToLower(outputStr)
		if strings.Contains(lowerOutput, "no tests collected") &&
		   !strings.Contains(lowerOutput, "error") {
			return []string{}, nil
		}

		// Include output to help user understand the problem
		return nil, fmt.Errorf("%s test discovery failed: %w\nOutput: %s", p.Executable, err, outputStr)
	}

	return parsePytestOutput(outputStr), nil
}

// parsePytestOutput extracts test node IDs from pytest --collect-only -q output
func parsePytestOutput(output string) []string {
	var tests []string
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and summary lines
		if line == "" || strings.Contains(line, " test") || strings.HasPrefix(line, "===") {
			continue
		}

		// Test node IDs contain "::" separator
		if strings.Contains(line, "::") {
			tests = append(tests, line)
		}
	}

	return tests
}
