// internal/connectors/vitest_test.go
package connectors

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVitestDetectFramework(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	t.Run("detects vitest command", func(t *testing.T) {
		connector := &VitestConnector{Executable: "vitest"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("returns false for nonexistent executable", func(t *testing.T) {
		connector := &VitestConnector{Executable: "nonexistent-vitest-binary"}
		got, err := connector.DetectFramework()

		assert.NoError(t, err)
		assert.False(t, got)
	})
}

func TestVitestGenerateConfig(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	t.Run("generates config with correct type, executable, and path", func(t *testing.T) {
		connector := &VitestConnector{Executable: "vitest"}
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "vitest", config.Type)
		assert.Equal(t, "vitest", config.Executable)
		assert.Equal(t, path, config.Path)
	})

	t.Run("generates config with custom executable", func(t *testing.T) {
		connector := &VitestConnector{Executable: "/custom/path/to/vitest"}
		path := "/path/to/project"

		config := connector.GenerateConfig(path)

		assert.Equal(t, "vitest", config.Type)
		assert.Equal(t, "/custom/path/to/vitest", config.Executable)
		assert.Equal(t, path, config.Path)
	})
}

func TestVitestDiscoverTests(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	// This test verifies the connector uses `vitest list --json` by testing
	// the behavior that would only work with those specific flags:
	// - Tests are collected, not run (requires list command)
	// - Output is JSON format (requires --json flag)
	// - All tests are discovered with full describe hierarchy

	projectDir := createVitestProject(t, map[string]string{
		"package.json": `{
  "name": "test-project",
  "type": "module",
  "devDependencies": {
    "vitest": "*"
  }
}`,
		"src/example.test.js": `import { describe, it, expect } from 'vitest'

describe('Math operations', () => {
  it('adds 1 + 2 to equal 3', () => {
    expect(1 + 2).toBe(3)
  })

  it('subtracts 5 - 3 to equal 2', () => {
    expect(5 - 3).toBe(2)
  })
})`,
	})

	connector = &VitestConnector{Executable: "vitest"}
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	// Verify both tests are discovered with full path and describe hierarchy
	assert.Contains(t, tests, "src/example.test.js > Math operations > adds 1 + 2 to equal 3")
	assert.Contains(t, tests, "src/example.test.js > Math operations > subtracts 5 - 3 to equal 2")
	// Verify tests were collected, not executed
	assert.Len(t, tests, 2, "should discover exactly 2 tests")
}

func TestVitestDiscoverTestsNestedDirectories(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	// This test verifies vitest discovers tests in deeply nested directory structures
	// and preserves the full path in test identifiers

	projectDir := createVitestProject(t, map[string]string{
		"package.json": `{
  "name": "test-project",
  "type": "module",
  "devDependencies": {
    "vitest": "*"
  }
}`,
		"tests/unit/auth/auth.test.js": `import { describe, it, expect } from 'vitest'

describe('Authentication', () => {
  it('validates login', () => {
    expect(true).toBe(true)
  })
})`,
		"tests/integration/api/handlers/api.test.js": `import { describe, it, expect } from 'vitest'

describe('API Handler', () => {
  it('handles request', () => {
    expect(true).toBe(true)
  })
})`,
	})

	connector = &VitestConnector{Executable: "vitest"}
	tests, err := connector.DiscoverTests(projectDir)

	assert.NoError(t, err)
	assert.Len(t, tests, 2)
	// Verify nested directories are discovered with full paths preserved
	assert.Contains(t, tests, "tests/unit/auth/auth.test.js > Authentication > validates login")
	assert.Contains(t, tests, "tests/integration/api/handlers/api.test.js > API Handler > handles request")
}

func TestVitestEmptyTestSuite(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	t.Run("returns empty list for project with no tests", func(t *testing.T) {
		projectDir := createVitestProject(t, map[string]string{
			"package.json": `{
  "name": "test-project",
  "type": "module",
  "devDependencies": {
    "vitest": "*"
  }
}`,
			"src/main.js": `export function main() {
  console.log("Hello, World!")
}`,
		})

		connector := &VitestConnector{Executable: "vitest"}

		// Skip test if vitest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("vitest not installed, skipping test")
		}

		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})

	t.Run("returns empty list for project with only non-test files", func(t *testing.T) {
		projectDir := createVitestProject(t, map[string]string{
			"package.json": `{
  "name": "test-project",
  "type": "module",
  "devDependencies": {
    "vitest": "*"
  }
}`,
			"src/utils.js": `export function add(a, b) {
  return a + b
}`,
		})

		connector := &VitestConnector{Executable: "vitest"}

		// Skip test if vitest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("vitest not installed, skipping test")
		}

		tests, err := connector.DiscoverTests(projectDir)

		assert.NoError(t, err)
		assert.Empty(t, tests)
	})
}

func TestVitestFrameworkNotFound(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	projectDir := createVitestProject(t, map[string]string{
		"package.json": `{
  "name": "test-project",
  "type": "module",
  "devDependencies": {
    "vitest": "*"
  }
}`,
		"src/example.test.js": `import { describe, it, expect } from 'vitest'

describe('Example', () => {
  it('works', () => {
    expect(true).toBe(true)
  })
})`,
	})

	connector = &VitestConnector{Executable: "nonexistent-vitest-binary"}
	_, err = connector.DiscoverTests(projectDir)

	assert.Error(t, err, "should return error when vitest command not found")

	// Verify error is "clear" per specification:
	// 1. Identifies the specific problematic value
	assert.Contains(t, err.Error(), "nonexistent-vitest-binary",
		"error should identify the executable that was not found")

	// 2. States what's wrong in plain language (not just "exit status 1")
	errMsg := strings.ToLower(err.Error())
	assert.True(t,
		strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "no such"),
		"error should clearly state the problem, got: %s", err.Error())

	// 3. Provides context about the operation
	assert.Contains(t, err.Error(), "test discovery",
		"error should provide context about what operation failed")
}

func TestVitestInvalidProjectStructure(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	t.Run("handles import errors from missing dependencies", func(t *testing.T) {
		projectDir := createVitestProject(t, map[string]string{
			"package.json": `{
  "name": "test-project",
  "type": "module",
  "devDependencies": {
    "vitest": "*"
  }
}`,
			"src/example.test.js": `import { nonexistentFunction } from './nonexistent-module'
import { describe, it, expect } from 'vitest'

describe('Example', () => {
  it('fails to load', () => {
    expect(true).toBe(true)
  })
})`,
		})

		connector := &VitestConnector{Executable: "vitest"}

		// Skip test if vitest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("vitest not installed, skipping test")
		}

		_, err = connector.DiscoverTests(projectDir)

		assert.Error(t, err, "should return error when project has import errors")
		assert.Contains(t, err.Error(), "test discovery",
			"error should provide context about the operation that failed")
	})
}

func TestVitestDiscoveryErrors(t *testing.T) {
	connector := DefaultVitestConnector()
	found, err := connector.DetectFramework()
	if !found || err != nil {
		t.Skip("vitest not installed, skipping vitest connector tests")
	}

	t.Run("handles syntax errors in test files", func(t *testing.T) {
		projectDir := createVitestProject(t, map[string]string{
			"package.json": `{
  "name": "test-project",
  "type": "module",
  "devDependencies": {
    "vitest": "*"
  }
}`,
			"src/broken.test.js": `import { describe, it, expect } from 'vitest'

describe('Broken test', () => {
  it('has syntax error', () => {
    expect(true).toBe(true
  // Missing closing parenthesis
})`,
		})

		connector := &VitestConnector{Executable: "vitest"}

		// Skip test if vitest is not installed
		detected, err := connector.DetectFramework()
		if err != nil || !detected {
			t.Skip("vitest not installed, skipping test")
		}

		_, err = connector.DiscoverTests(projectDir)

		assert.Error(t, err, "should return error when test file has syntax errors")
		// Error message should include context from vitest's output
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "test discovery") || strings.Contains(errMsg, "failed"),
			"error should indicate test discovery failed, got: %s", err.Error())
	})
}

// Helper function to create a vitest test project
func createVitestProject(t *testing.T, files map[string]string) string {
	t.Helper()
	tempDir := t.TempDir()

	// Create test files
	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if dir != tempDir {
			err := os.MkdirAll(dir, 0755)
			assert.NoError(t, err)
		}
		err := os.WriteFile(fullPath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	return tempDir
}
