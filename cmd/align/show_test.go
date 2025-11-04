package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowParses(t *testing.T) {
	// Create temp directory with spec file
	tempDir := t.TempDir()
	specContent := `# Test Specification

## Feature One
Description of feature one.

**Test:** ` + "`TestFeatureOne`" + `
`
	specPath := filepath.Join(tempDir, "test.md")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	assert.NoError(t, err)

	var stdout, stderr bytes.Buffer
	run([]string{"show", specPath}, &stdout, &stderr)

	// Verify it parsed successfully (no errors, output produced)
	assert.Empty(t, stderr.String(), "should parse without errors")
	assert.NotEmpty(t, stdout.String(), "should produce output after parsing")
}

func TestShowDisplaysTitles(t *testing.T) {
	// Create temp directory with spec file
	tempDir := t.TempDir()
	specContent := `# Test Specification

## Feature One
Description of feature one.

### Subfeature A
Details about subfeature A.

**Test:** ` + "`TestSubfeatureA`" + `

## Feature Two
Description of feature two.

**Test:** ` + "`TestFeatureTwo`" + `
`
	specPath := filepath.Join(tempDir, "test.md")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	assert.NoError(t, err)

	var stdout, stderr bytes.Buffer
	run([]string{"show", specPath}, &stdout, &stderr)

	output := stdout.String()

	// Verify section titles are displayed
	assert.Contains(t, output, "Test Specification", "should display section titles")
	assert.Contains(t, output, "Feature One", "should display section titles")
	assert.Contains(t, output, "Feature Two", "should display section titles")
	assert.Contains(t, output, "Subfeature A", "should display section titles")

	// Verify test references are displayed (per spec/output_format.md § 3.2)
	// Spec promises: "displays test references in the format 'Test: test_name'"
	// Note: Output contains ANSI color codes, so we check for "Test:" and test names separately
	assert.Contains(t, output, "Test:", "should display 'Test:' label")
	assert.Contains(t, output, "TestSubfeatureA",
		"should display test reference name for Subfeature A")
	assert.Contains(t, output, "TestFeatureTwo",
		"should display test reference name for Feature Two")
}

func TestShowDisplaysHierarchy(t *testing.T) {
	// Create temp directory with spec file
	tempDir := t.TempDir()
	specContent := `# Test Specification

## Feature One
Description of feature one.

### Subfeature A
Details about subfeature A.

#### Deep Item
Even deeper.

**Test:** ` + "`TestDeepItem`" + `
`
	specPath := filepath.Join(tempDir, "test.md")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	assert.NoError(t, err)

	var stdout, stderr bytes.Buffer
	run([]string{"show", specPath}, &stdout, &stderr)

	output := stdout.String()

	// Verify hierarchical structure is displayed
	assert.Contains(t, output, "Test Specification", "should show root section")
	assert.Contains(t, output, "Feature One", "should show parent section")
	assert.Contains(t, output, "Subfeature A", "should show child section")
	assert.Contains(t, output, "Deep Item", "should show deeply nested section")

	// Check that sections appear in hierarchical order
	rootIdx := strings.Index(output, "Test Specification")
	featureIdx := strings.Index(output, "Feature One")
	subfeatureIdx := strings.Index(output, "Subfeature A")
	deepIdx := strings.Index(output, "Deep Item")
	assert.Greater(t, featureIdx, rootIdx, "child should appear after parent")
	assert.Greater(t, subfeatureIdx, featureIdx, "nested child should appear after parent")
	assert.Greater(t, deepIdx, subfeatureIdx, "deeply nested child should appear after parent")

	// Verify middle dot (·) indentation is used per spec
	assert.Contains(t, output, "·", "should use middle dot for indentation")

	// Verify indentation depth increases with nesting
	// Level 2 (Feature One) should have one middle dot: "· "
	// Level 3 (Subfeature A) should have two middle dots: "· · "
	// Level 4 (Deep Item) should have three middle dots: "· · · "
	lines := strings.Split(output, "\n")
	foundLevel2 := false
	foundLevel3 := false
	foundLevel4 := false
	for _, line := range lines {
		if strings.Contains(line, "Feature One") {
			assert.True(t, strings.Contains(line, "· ") && !strings.Contains(line, "· · "),
				"level 2 should have one middle dot, got: %s", line)
			foundLevel2 = true
		}
		if strings.Contains(line, "Subfeature A") {
			assert.True(t, strings.Contains(line, "· · "),
				"level 3 should have two middle dots, got: %s", line)
			foundLevel3 = true
		}
		if strings.Contains(line, "Deep Item") {
			assert.True(t, strings.Contains(line, "· · · "),
				"level 4 should have three middle dots, got: %s", line)
			foundLevel4 = true
		}
	}
	assert.True(t, foundLevel2, "should find level 2 section")
	assert.True(t, foundLevel3, "should find level 3 section")
	assert.True(t, foundLevel4, "should find level 4 section")
}

func TestShowExitCode(t *testing.T) {
	// Create temp directory with spec file
	tempDir := t.TempDir()
	specContent := `# Test Specification

## Feature One
**Test:** ` + "`TestFeatureOne`" + `
`
	specPath := filepath.Join(tempDir, "test.md")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	assert.NoError(t, err)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"show", specPath}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should exit with code 0 on success")
}

func TestShowFileNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"show", "nonexistent.md"}, &stdout, &stderr)

	assert.Equal(t, 1, exitCode, "should exit with code 1 when file not found")
	assert.Contains(t, strings.ToLower(stderr.String()), "not found", "should report file not found")
}

func TestShowEmptyFile(t *testing.T) {
	// Create temp directory with empty spec file
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "empty.md")
	err := os.WriteFile(specPath, []byte(""), 0644)
	assert.NoError(t, err)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"show", specPath}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should successfully process empty files with exit code 0")
}

func TestShowDirectory(t *testing.T) {
	// Create temp directory with spec files
	tempDir := t.TempDir()

	// Create a spec file in the directory
	specContent := `# Test Specification

## Feature One
**Test:** ` + "`TestFeatureOne`" + `
`
	specPath := filepath.Join(tempDir, "test.md")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	assert.NoError(t, err)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"show", tempDir}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should successfully process directory with exit code 0")
	assert.Empty(t, stderr.String(), "should parse directory without errors")
	assert.NotEmpty(t, stdout.String(), "should produce output when showing directory")
}

func TestShowDirectoryDisplaysAll(t *testing.T) {
	// Create temp directory with multiple spec files
	tempDir := t.TempDir()

	// Create first spec file
	spec1Content := `# Specification One

## Feature A
**Test:** ` + "`TestFeatureA`" + `
`
	spec1Path := filepath.Join(tempDir, "spec1.md")
	err := os.WriteFile(spec1Path, []byte(spec1Content), 0644)
	assert.NoError(t, err)

	// Create second spec file
	spec2Content := `# Specification Two

## Feature B
**Test:** ` + "`TestFeatureB`" + `
`
	spec2Path := filepath.Join(tempDir, "spec2.md")
	err = os.WriteFile(spec2Path, []byte(spec2Content), 0644)
	assert.NoError(t, err)

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"show", tempDir}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "should successfully process directory")
	output := stdout.String()

	// Verify both specifications are displayed
	assert.Contains(t, output, "Specification One", "should display first spec")
	assert.Contains(t, output, "Feature A", "should display section from first spec")
	assert.Contains(t, output, "TestFeatureA", "should display test from first spec")

	assert.Contains(t, output, "Specification Two", "should display second spec")
	assert.Contains(t, output, "Feature B", "should display section from second spec")
	assert.Contains(t, output, "TestFeatureB", "should display test from second spec")
}
