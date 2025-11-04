// internal/parser/parser.go
package parser

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/Alge/aligned/internal/spec"
)

// ParseMarkdown parses markdown content into a Specification
func ParseMarkdown(content string) (*spec.Specification, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	var sections []*spec.Section
	var currentContent strings.Builder
	var lastSection *spec.Section

	headingPattern := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line is a heading
		if matches := headingPattern.FindStringSubmatch(line); matches != nil {
			// Save content to previous section if exists
			if lastSection != nil {
				lastSection.Content = strings.TrimSpace(currentContent.String())
				lastSection.TestName = ExtractTestReference(lastSection.Content)
			}

			// Create new section
			level := len(matches[1]) // Count # symbols
			title := strings.TrimSpace(matches[2])

			section := &spec.Section{
				Level:    level,
				Title:    title,
				Children: []*spec.Section{},
			}

			sections = append(sections, section)
			lastSection = section
			currentContent.Reset()
		} else {
			// Accumulate content for current section
			currentContent.WriteString(line)
			currentContent.WriteString("\n")
		}
	}

	// Don't forget the last section
	if lastSection != nil {
		lastSection.Content = strings.TrimSpace(currentContent.String())
		lastSection.TestName = ExtractTestReference(lastSection.Content)
	}

	// Build tree structure from flat list
	tree := buildTree(sections)

	return &spec.Specification{
		Sections: tree,
	}, nil
}

// buildTree converts a flat list of sections into a hierarchical tree
func buildTree(sections []*spec.Section) []*spec.Section {
	if len(sections) == 0 {
		return []*spec.Section{}
	}

	var roots []*spec.Section
	var stack []*spec.Section // Stack of potential parents

	for _, section := range sections {
		// Pop stack until we find appropriate parent
		for len(stack) > 0 && stack[len(stack)-1].Level >= section.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// This is a root section
			roots = append(roots, section)
		} else {
			// This is a child of the section on top of stack
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, section)
			section.Parent = parent
		}

		// Push current section onto stack
		stack = append(stack, section)
	}

	return roots
}

// ExtractTestReference finds and extracts the test name from content
// Looks for pattern: **Test:** `TestName`
// Returns empty string if no test reference found
func ExtractTestReference(content string) string {
	// Pattern matches: **Test:** followed by required backticks containing test reference
	// The test reference must be in backticks and can contain any characters except backticks
	// This supports various test name formats including those with spaces:
	// - Go: package.TestName
	// - Pytest: tests/test_file.py::TestClass::test_method
	// - ExUnit: test/file.exs:Module:test description with spaces
	// - Gleam: module@submodule.function_name_test
	// - Vitest: src/file.test.js > describe > test name
	pattern := regexp.MustCompile("(?m)^\\*\\*[Tt]est:\\*\\*\\s*`([^`]+)`")

	matches := pattern.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}