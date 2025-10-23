package connectors

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Alge/aligned/internal/config"
)

type GleamConnector struct {
	Executable string
}

func NewGleamConnector(executable string) *GleamConnector {
	if executable == "" {
		executable = "gleam"
	}
	return &GleamConnector{
		Executable: executable,
	}
}

// DefaultGleamConnector returns a GleamConnector with default configuration
func DefaultGleamConnector() *GleamConnector {
	return &GleamConnector{
		Executable: "gleam",
	}
}

// DetectFramework checks if the Gleam executable is available
func (g *GleamConnector) DetectFramework() (bool, error) {
	_, err := exec.LookPath(g.Executable)
	return err == nil, nil
}

// GenerateConfig creates a default connector configuration for Gleam
func (g *GleamConnector) GenerateConfig(path string) config.ConnectorConfig {
	return config.ConnectorConfig{
		Type:       "gleam",
		Executable: g.Executable,
		Path:       path,
	}
}

// DiscoverTests discovers Gleam tests by parsing source files in the test/ directory
func (g *GleamConnector) DiscoverTests(path string) ([]string, error) {
	// Check if gleam.toml exists (validates this is a Gleam project)
	gleamTomlPath := filepath.Join(path, "gleam.toml")
	if _, err := os.Stat(gleamTomlPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("gleam test discovery failed: gleam.toml not found in project root")
	}

	// Look for test directory
	testDir := filepath.Join(path, "test")
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		// No test directory is valid - return empty list
		return []string{}, nil
	}

	var tests []string

	// Walk through test directory to find all .gleam files
	err := filepath.Walk(testDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			// Check if it's a permission error
			if os.IsPermission(err) {
				return fmt.Errorf("test discovery failed: permission denied reading %s", filePath)
			}
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .gleam files
		if !strings.HasSuffix(filePath, ".gleam") {
			return nil
		}

		// Parse the file to find test functions
		fileTests, err := parseGleamTestFile(filePath)
		if err != nil {
			return fmt.Errorf("test discovery failed: %w", err)
		}

		// Convert file path to module name
		moduleName, err := filePathToModuleName(filePath, testDir)
		if err != nil {
			return err
		}

		// Add module prefix to each test
		for _, test := range fileTests {
			tests = append(tests, moduleName+"."+test)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tests, nil
}

// filePathToModuleName converts a file path to a Gleam module name
// Example: test/unit/auth/login_test.gleam -> unit@auth@login_test
func filePathToModuleName(filePath, testDir string) (string, error) {
	// Get relative path from test directory
	relPath, err := filepath.Rel(testDir, filePath)
	if err != nil {
		return "", err
	}

	// Remove .gleam extension
	modulePath := strings.TrimSuffix(relPath, ".gleam")

	// Replace directory separators with @ (Gleam's module separator for nested paths)
	moduleName := strings.ReplaceAll(modulePath, string(filepath.Separator), "@")

	return moduleName, nil
}

// parseGleamTestFile parses a Gleam file and returns test function names
func parseGleamTestFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied reading file: %s", filePath)
		}
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	defer file.Close()

	var tests []string
	scanner := bufio.NewScanner(file)

	// Pattern to match public functions ending in _test
	// Example: pub fn hello_world_test() {
	testFuncPattern := regexp.MustCompile(`^\s*pub\s+fn\s+([a-z_][a-z0-9_]*_test)\s*\(`)

	for scanner.Scan() {
		line := scanner.Text()

		if matches := testFuncPattern.FindStringSubmatch(line); matches != nil {
			testName := matches[1]
			tests = append(tests, testName)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error parsing Gleam file %s: invalid syntax or encoding", filePath)
	}

	return tests, nil
}
