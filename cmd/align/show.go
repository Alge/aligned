package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Alge/aligned/internal/parser"
	"github.com/Alge/aligned/internal/spec"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

func show(args []string, stdout, stderr io.Writer) int {
	// Check arguments
	if len(args) < 1 {
		fmt.Fprintln(stderr, "Usage: align show <spec-file>")
		return 1
	}
	
	specPath := args[0]
	
	// Check if file exists
	content, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(stderr, "Error: Spec file not found: %s\n", specPath)
			return 1
		}
		fmt.Fprintf(stderr, "Error reading spec file: %v\n", err)
		return 1
	}
	
	// Parse the spec
	specification, err := parser.ParseMarkdown(string(content))
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing spec: %v\n", err)
		return 1
	}
	
	// Display the spec structure
	printSpecification(specification, stdout)
	
	return 0
}

func printSpecification(specification *spec.Specification, stdout io.Writer) {
	for _, section := range specification.Sections {
		printSection(section, 0, stdout)
	}
}

func printSection(section *spec.Section, indent int, stdout io.Writer) {
	// Print indentation with middle dots (gray)
	prefix := colorGray + strings.Repeat("· ", indent) + colorReset
	
	// Print section title (blue for headings)
	fmt.Fprintf(stdout, "%s%s%s%s\n", prefix, colorBlue, section.Title, colorReset)
	
	// If section has a test, show it (green)
	if section.TestName != "" {
		fmt.Fprintf(stdout, "%s  %sTest: %s%s%s\n", 
			colorGray+strings.Repeat("· ", indent)+colorReset,
			colorGray,
			colorGreen,
			section.TestName,
			colorReset)
	} else if section.IsLeaf() {
		// Leaf section without test - show warning
		fmt.Fprintf(stdout, "%s  %s⚠ Missing test reference%s\n",
			colorGray+strings.Repeat("· ", indent)+colorReset,
			colorRed,
			colorReset)
	}
	
	// Print children recursively
	for _, child := range section.Children {
		printSection(child, indent+1, stdout)
	}
}