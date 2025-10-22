package connectors

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Alge/aligned/internal/config"
)

type GoConnector struct {
	Executable string
}

func NewGoConnector(executable string) *GoConnector {
    if executable == "" {
        executable = "go"
    }
    return &GoConnector{
        Executable: executable,
    }
}

// DefaultGoConnector returns a GoConnector with default configuration
func DefaultGoConnector() *GoConnector {
    return &GoConnector{
        Executable: "go",
    }
}

// DetectFramework checks if the Go executable is available
func (g *GoConnector) DetectFramework() (bool, error) {
	_, err := exec.LookPath(g.Executable)
	return err == nil, nil
}

// GenerateConfig creates a default connector configuration for Go
func (g *GoConnector) GenerateConfig(path string) config.ConnectorConfig {
	return config.ConnectorConfig{
		Type:       "go",
		Executable: g.Executable,
		Path:       path,
	}
}

// ValidateCompatibility checks if the Go version is 1.13 or higher
func (g *GoConnector) ValidateCompatibility() (bool, error) {
	cmd := exec.Command(g.Executable, "version")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to get go version: %w", err)
	}

	version := string(output)
	// Example output: "go version go1.23.2 linux/amd64"
	re := regexp.MustCompile(`go version go(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 3 {
		return false, fmt.Errorf("failed to parse go version from: %s", version)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])

	// go test -list was added in Go 1.13
	if major > 1 || (major == 1 && minor >= 13) {
		return true, nil
	}

	return false, nil
}

// ValidateConfiguration validates the connector configuration
func (g *GoConnector) ValidateConfiguration() error {
	if g.Executable == "" {
		return fmt.Errorf("executable path cannot be empty")
	}
	return nil
}

// DiscoverTests discovers Go tests in the given path with a default timeout
func (g *GoConnector) DiscoverTests(path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return g.DiscoverTestsWithContext(ctx, path)
}

// DiscoverTestsWithContext discovers Go tests in the given path with a context
func (g *GoConnector) DiscoverTestsWithContext(ctx context.Context, path string) ([]string, error) {
	cmd := exec.CommandContext(ctx, g.Executable, "test", "-list=.", "./...")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("test discovery timed out: %w", ctx.Err())
		}
		// Include output to help user understand the problem
		return nil, fmt.Errorf("%s test discovery failed: %w\nOutput: %s", g.Executable, err, string(output))
	}

	return parseGoTestOutput(string(output)), nil
}

// parseGoTestOutput extracts package-qualified test names from go test -list output
func parseGoTestOutput(output string) []string {
	var tests []string
	var currentTests []string
	var moduleName string
	
	scanner := bufio.NewScanner(strings.NewReader(output))
	testPattern := regexp.MustCompile(`^(Test[A-Za-z0-9_]+)$`)
	okPattern := regexp.MustCompile(`^ok\s+(\S+)`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check if this is a test name
		if matches := testPattern.FindStringSubmatch(line); matches != nil {
			currentTests = append(currentTests, matches[1])
			continue
		}
		
		// Check if this is an "ok" line with package
		if matches := okPattern.FindStringSubmatch(line); matches != nil {
			fullPath := matches[1]
			
			// First package path tells us the module name
			if moduleName == "" {
				parts := strings.Split(fullPath, "/")
				moduleName = parts[0]
			}
			
			// Strip module prefix to get relative path
			relativePath := strings.TrimPrefix(fullPath, moduleName+"/")
			if relativePath == moduleName {
				// Root package - use module name
				relativePath = moduleName
			}
			
			// Add all accumulated tests with package prefix
			for _, test := range currentTests {
				tests = append(tests, relativePath+"."+test)
			}
			currentTests = nil
		}
	}
	
	return tests
}