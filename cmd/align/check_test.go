package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupTestProject creates a temp directory with config, go.mod, and optional test/spec files
func setupTestProject(t *testing.T, testFiles map[string]string, specContent string) (tempDir, specPath string) {
	t.Helper()
	
	tempDir = t.TempDir()
	
	// Create config
	configContent := `connectors:
  - type: go
    executable: go
    path: .
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)
	
	// Create go.mod
	goModContent := `module testproject
go 1.23
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	assert.NoError(t, err)
	
	// Create test files
	for filename, content := range testFiles {
		err = os.WriteFile(filepath.Join(tempDir, filename), []byte(content), 0644)
		assert.NoError(t, err)
	}
	
	// Create spec file if provided
	if specContent != "" {
		specPath = filepath.Join(tempDir, "spec.md")
		err = os.WriteFile(specPath, []byte(specContent), 0644)
		assert.NoError(t, err)
	}
	
	return tempDir, specPath
}

func TestCheckSuccessCase(t *testing.T) {
	testFiles := map[string]string{
		"example_test.go": `package example
import "testing"
func TestFeature(t *testing.T) {}
`,
	}

	specContent := `# Test Spec

## Feature
Description of the feature.

**Test:** ` + "`testproject.TestFeature`" + `
`

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specPath}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 when all specs covered")
	assert.Contains(t, strings.ToLower(stdout.String()), "covered", "should display success message")
}

func TestCheckMissingTestReferences(t *testing.T) {
	testFiles := map[string]string{
		"example_test.go": `package example
import "testing"
func TestSomething(t *testing.T) {}
`,
	}

	specContent := `# Test Spec

## Feature One
Parent section.

### Subfeature
This is a leaf without a test reference.
`

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specPath}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1")
	assert.Contains(t, stdout.String(), "Subfeature", "should report which specs are missing test references")
	assert.Contains(t, strings.ToLower(stdout.String()), "missing test", "should report which specs are missing test references")
}

func TestCheckTestsNotFound(t *testing.T) {
	testFiles := map[string]string{
		"example_test.go": `package example
import "testing"
func TestSomethingElse(t *testing.T) {}
`,
	}

	specContent := `# Test Spec

## Feature
Description.

**Test:** ` + "`testproject.TestNonexistent`" + `
`

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specPath}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1")
	assert.Contains(t, stdout.String(), "TestNonexistent", "should report which test references cannot be found")
	assert.Contains(t, strings.ToLower(stdout.String()), "not found", "should report which test references cannot be found")
}

func TestCheckLoadsSingleFile(t *testing.T) {
	testFiles := map[string]string{
		"example_test.go": `package example
import "testing"
func TestFeature(t *testing.T) {}
`,
	}

	specContent := `# Test Spec

## Feature
Description.

**Test:** ` + "`testproject.TestFeature`" + `
`

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specPath}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should load and validate single specification file")
}

func TestCheckLoadsDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create config
	configContent := `connectors:
  - type: go
    executable: go
    path: .
`
	err := os.WriteFile(filepath.Join(tempDir, ".align.yml"), []byte(configContent), 0644)
	assert.NoError(t, err)
	
	// Create go.mod
	goModContent := `module testproject
go 1.23
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	assert.NoError(t, err)
	
	// Create test file
	testContent := `package example
import "testing"
func TestFeatureOne(t *testing.T) {}
func TestFeatureTwo(t *testing.T) {}
`
	err = os.WriteFile(filepath.Join(tempDir, "example_test.go"), []byte(testContent), 0644)
	assert.NoError(t, err)
	
	// Create spec directory with multiple files
	specDir := filepath.Join(tempDir, "specs")
	err = os.MkdirAll(specDir, 0755)
	assert.NoError(t, err)
	
	// First spec file
	spec1Content := `# Spec One

## Feature One
Description.

**Test:** ` + "`testproject.TestFeatureOne`" + `
`
	err = os.WriteFile(filepath.Join(specDir, "spec1.md"), []byte(spec1Content), 0644)
	assert.NoError(t, err)
	
	// Second spec file
	spec2Content := `# Spec Two

## Feature Two
Description.

**Test:** ` + "`testproject.TestFeatureTwo`" + `
`
	err = os.WriteFile(filepath.Join(specDir, "spec2.md"), []byte(spec2Content), 0644)
	assert.NoError(t, err)
	
	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specDir}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should recursively load all .md files in directory")
	assert.Contains(t, strings.ToLower(stdout.String()), "covered")
	// Should contain both spec titles
	assert.Contains(t, stdout.String(), "Spec One")
	assert.Contains(t, stdout.String(), "Spec Two")
}

func TestCheckCollapsedOutput(t *testing.T) {
	// Create a project with all passing tests
	testFiles := map[string]string{
		"example_test.go": `package example
import "testing"
func TestFeatureOne(t *testing.T) {}
func TestFeatureTwo(t *testing.T) {}
func TestSubFeatureOne(t *testing.T) {}
func TestSubFeatureTwo(t *testing.T) {}
`,
	}

	specContent := `# Main Spec

## Feature Group
Parent section with multiple features.

### Feature One
First feature.
**Test:** ` + "`testproject.TestFeatureOne`" + `

### Feature Two
Second feature.
**Test:** ` + "`testproject.TestFeatureTwo`" + `

## Sub Features

### Sub Feature One
**Test:** ` + "`testproject.TestSubFeatureOne`" + `

### Sub Feature Two
**Test:** ` + "`testproject.TestSubFeatureTwo`" + `
`

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specPath}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode)
	output := stdout.String()

	// Should show collapsed format with single top-level count
	assert.Contains(t, output, "Main Spec")
	assert.Contains(t, output, "(4/4 passed)") // Total count for Main Spec

	// Should NOT show individual feature details when collapsed
	assert.NotContains(t, output, "First feature")
	assert.NotContains(t, output, "Second feature")
	assert.NotContains(t, output, "Feature Group") // Collapsed under Main Spec
	assert.NotContains(t, output, "Sub Features")  // Collapsed under Main Spec
}

func TestCheckExpandsFailures(t *testing.T) {
	// Create project with one passing and one failing test
	testFiles := map[string]string{
		"example_test.go": `package example
import "testing"
func TestPassingFeature(t *testing.T) {}
// TestFailingFeature is intentionally missing
`,
	}

	specContent := `# Main Spec

## Passing Group

### Passing Feature
This passes.
**Test:** ` + "`testproject.TestPassingFeature`" + `

## Failing Group

### Missing Test Feature
This will fail.
**Test:** ` + "`testproject.TestFailingFeature`" + `

### Missing Reference Feature
This has no test reference.
`

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specPath}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode) // Should fail due to missing tests
	output := stdout.String()

	// Main Spec should be expanded because it contains failures
	assert.Contains(t, output, "✗")
	assert.Contains(t, output, "Main Spec")

	// Passing group should be collapsed
	assert.Contains(t, output, "✓")
	assert.Contains(t, output, "Passing Group")
	assert.Contains(t, output, "(1/1 passed)")
	assert.NotContains(t, output, "Passing Feature") // Should be collapsed

	// Failing group should be expanded
	assert.Contains(t, output, "✗")
	assert.Contains(t, output, "Failing Group")
	assert.Contains(t, output, "Missing Test Feature")
	assert.Contains(t, output, "Test not found")
	assert.Contains(t, output, "Missing Reference Feature")
	assert.Contains(t, output, "Missing test reference")
}

func TestCheckVerboseOutput(t *testing.T) {
	// Same setup as collapsed test
	testFiles := map[string]string{
		"example_test.go": `package example
import "testing"
func TestFeatureOne(t *testing.T) {}
func TestFeatureTwo(t *testing.T) {}
`,
	}

	specContent := `# Main Spec

## Feature Group

### Feature One
First feature.
**Test:** ` + "`testproject.TestFeatureOne`" + `

### Feature Two
Second feature.
**Test:** ` + "`testproject.TestFeatureTwo`" + `
`

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", "-v", specPath}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode)
	output := stdout.String()

	// Verbose mode should show all details
	assert.Contains(t, output, "Feature One")
	assert.Contains(t, output, "Feature Two")
	assert.Contains(t, output, "TestFeatureOne")
	assert.Contains(t, output, "TestFeatureTwo")

	// Should NOT show collapsed counts in verbose mode
	assert.NotContains(t, output, "(2/2 passed)")
}

func TestCheckInterfaceValidation(t *testing.T) {
	// Create specs with both valid and invalid implementations
	specContent := `# Framework

## API Interface [INTERFACE]

### Get Data
Get data method.

### Process Data
Process data method.

## Valid Implementation [IMPLEMENTS: API Interface]

### Get Data
Valid impl has this.
**Test:** ` + "`testproject.TestGetData`" + `

### Process Data
Valid impl has this too.
**Test:** ` + "`testproject.TestProcessData`" + `

### Extra Method
Implementations can have extra methods.
**Test:** ` + "`testproject.TestExtra`" + `

## Invalid Implementation [IMPLEMENTS: API Interface]

### Get Data
Has only one of the required methods.
**Test:** ` + "`testproject.TestGetDataInvalid`" + `
`

	testFiles := map[string]string{
		"impl_test.go": `package impl
import "testing"
func TestGetData(t *testing.T) {}
func TestProcessData(t *testing.T) {}
func TestExtra(t *testing.T) {}
func TestGetDataInvalid(t *testing.T) {}
`,
	}

	tempDir, specPath := setupTestProject(t, testFiles, specContent)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"check", specPath}, &stdout, &stderr)

	// Should fail due to invalid implementation
	assert.Equal(t, 1, exitCode, "should exit with code 1 when implementations are missing required interface sections")

	output := stdout.String()

	// Should report what's missing
	assert.Contains(t, output, "process data", "should report which sections are missing")
	assert.Contains(t, output, "Invalid Implementation", "should report which sections are missing")
}