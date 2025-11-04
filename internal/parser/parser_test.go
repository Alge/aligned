// internal/parser/parser_test.go
package parser

import (
	"testing"
)

func TestParseMarkdownHeadings(t *testing.T) {
	// This test verifies full markdown parsing workflow including:
	// 1. Heading extraction (title and level)
	// 2. Tree structure building (parent-child relationships)
	// 3. Test reference extraction during parsing (integration of ExtractTestReference)
	//
	// By checking that child1.TestName and child2.TestName are correctly populated,
	// this test proves that ExtractTestReference is called and integrated into the
	// parsing workflow, not just tested in isolation.

	input := `# Specification Parsing

## Parse Markdown headings
Read a Markdown file and extract heading text and level (# = 1, ## = 2, ### = 3).

**Test:** ` + "`TestParseMarkdownHeadings`" + `

## Extract test reference from specification
Find lines matching "**Test:** ` + "`test_name`" + `" and extract the test name.

**Test:** ` + "`TestExtractTestReference`" + `
`

	result, err := ParseMarkdown(input)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	// Should have 1 root section ("Specification Parsing")
	if len(result.Sections) != 1 {
		t.Fatalf("got %d root sections, want 1", len(result.Sections))
	}

	root := result.Sections[0]

	// Check root section
	if root.Level != 1 {
		t.Errorf("root.Level = %d, want 1", root.Level)
	}
	if root.Title != "Specification Parsing" {
		t.Errorf("root.Title = %q, want %q", root.Title, "Specification Parsing")
	}

	// Check children (should have 2 subsections)
	if len(root.Children) != 2 {
		t.Fatalf("root has %d children, want 2", len(root.Children))
	}

	// First child: "Parse Markdown headings"
	child1 := root.Children[0]
	if child1.Level != 2 {
		t.Errorf("child1.Level = %d, want 2", child1.Level)
	}
	if child1.Title != "Parse Markdown headings" {
		t.Errorf("child1.Title = %q, want %q", child1.Title, "Parse Markdown headings")
	}
	// Verify test reference extraction is integrated into parsing
	if child1.TestName != "TestParseMarkdownHeadings" {
		t.Errorf("child1.TestName = %q, want %q", child1.TestName, "TestParseMarkdownHeadings")
	}

	// Second child: "Extract test reference from specification"
	child2 := root.Children[1]
	if child2.Level != 2 {
		t.Errorf("child2.Level = %d, want 2", child2.Level)
	}
	if child2.Title != "Extract test reference from specification" {
		t.Errorf("child2.Title = %q, want %q", child2.Title, "Extract test reference from specification")
	}
	// Verify test reference extraction is integrated into parsing
	if child2.TestName != "TestExtractTestReference" {
		t.Errorf("child2.TestName = %q, want %q", child2.TestName, "TestExtractTestReference")
	}
}

func TestExtractTestReference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard format with backticks",
			input:    "**Test:** `TestSomething`",
			expected: "TestSomething",
		},
		{
			name:     "without backticks",
			input:    "**Test:** TestSomething",
			expected: "", // Backticks are now required
		},
		{
			name: "with surrounding text",
			input: `Some description here.

**Test:** ` + "`TestName`" + `

More text after.`,
			expected: "TestName",
		},
		{
			name:     "no test reference",
			input:    "Just regular content without a test.",
			expected: "",
		},
		{
			name:     "qualified test name with package",
			input:    "**Test:** `package.TestSomething`",
			expected: "package.TestSomething",
		},
		{
			name:     "qualified test name with module path",
			input:    "**Test:** `module/package.TestSomething`",
			expected: "module/package.TestSomething",
		},
		{
			name:     "qualified test name without backticks",
			input:    "**Test:** testproject.TestFeature",
			expected: "", // Backticks are now required
		},
		{
			name:     "test name with spaces (ExUnit)",
			input:    "**Test:** `test/file.exs:Module:test description with spaces`",
			expected: "test/file.exs:Module:test description with spaces",
		},
		{
			name:     "test name with special characters (Vitest)",
			input:    "**Test:** `src/file.test.js > describe > test name`",
			expected: "src/file.test.js > describe > test name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTestReference(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractTestReference() = %q, want %q", result, tt.expected)
			}
		})
	}
}