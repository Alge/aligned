package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alge/aligned/internal/config"
	"github.com/Alge/aligned/internal/connectors"
	"github.com/Alge/aligned/internal/logger"
	"github.com/Alge/aligned/internal/parser"
	"github.com/Alge/aligned/internal/spec"
)

func check(args []string, stdout, stderr io.Writer) int {
	// Check for verbose flag
	verbose := false
	specPath := ""
	
	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else if !strings.HasPrefix(arg, "-") {
			specPath = arg
		}
	}
	
	if specPath == "" {
		fmt.Fprintln(stderr, "Usage: align check [-v] <spec-file-or-directory>")
		return 1
	}
	
	// Load configuration
	configPath := filepath.Join(".", ".align.yml")
	cfg, err := config.LoadConfiguration(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(stderr, "Error: .align.yml not found")
			return 1
		}
		fmt.Fprintf(stderr, "Error: Invalid configuration: %v\n", err)
		return 1
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(stderr, "Error: Invalid configuration: %v\n", err)
		return 1
	}
	
	// Load specification (file or directory)
	specification, err := loadSpecification(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(stderr, "Error: Spec path not found: %s\n", specPath)
			return 1
		}
		fmt.Fprintf(stderr, "Error loading spec: %v\n", err)
		return 1
	}
	
	// Discover all tests
	var allTests []string
	for _, connectorCfg := range cfg.Connectors {
		var connector connectors.Connector
		
		switch connectorCfg.Type {
		case "go":
			executable := connectorCfg.Executable
			if executable == "" {
				executable = "go"
			}
			connector = connectors.NewGoConnector(executable)
		default:
			fmt.Fprintf(stderr, "Error: Unsupported connector type: %s\n", connectorCfg.Type)
			return 1
		}
		
		tests, err := connector.DiscoverTests(connectorCfg.Path)
		if err != nil {
			fmt.Fprintf(stderr, "Error discovering tests: %v\n", err)
			return 1
		}
		
		allTests = append(allTests, tests...)
	}
	
	// Create a set of discovered tests for quick lookup
	testSet := make(map[string]bool)
	for _, test := range allTests {
		testSet[test] = true
	}
	
	// Check coverage
	hasErrors := false
	missingReferences := []string{}
	testsNotFound := []string{}
	
	log := logger.Debug()
	
	for _, leaf := range specification.AllLeaves() {
		log.Debug("checking leaf",
			"title", leaf.Title,
			"isLeaf", leaf.IsLeaf(),
			"requiresTest", leaf.RequiresTest(),
			"hasParent", leaf.Parent != nil,
		)
		
		if leaf.Parent != nil {
			log.Debug("parent info",
				"parentTitle", leaf.Parent.Title,
				"parentIsInterface", leaf.Parent.IsInterface(),
			)
		}
		
		// Use RequiresTest() to handle interfaces properly
		if leaf.RequiresTest() {
			if !leaf.HasTest() {
				missingReferences = append(missingReferences, leaf.Title)
				hasErrors = true
				log.Debug("missing test reference", "title", leaf.Title)
			} else if !testSet[leaf.TestName] {
				testsNotFound = append(testsNotFound, leaf.TestName)
				hasErrors = true
				log.Debug("test not found", "title", leaf.Title, "testName", leaf.TestName)
			}
		} else {
			log.Debug("test not required", "title", leaf.Title)
		}
	}
	
	// Validate interface implementations
	interfaceErrors := specification.ValidateInterfaces()
	if len(interfaceErrors) > 0 {
		hasErrors = true
		log.Debug("interface validation errors", "count", len(interfaceErrors))
	}
	
	// Report results
	fmt.Fprintln(stdout, "Specification coverage report:")
	fmt.Fprintln(stdout, "")
	
	if verbose {
		// Show full tree in verbose mode
		printSpecificationWithStatusAndErrors(specification, testSet, interfaceErrors, stdout)
	} else {
		// Show collapsed view by default
		printSpecificationCollapsedWithErrors(specification, testSet, interfaceErrors, stdout)
	}
	
	fmt.Fprintln(stdout, "")
	
	if hasErrors {
		if len(missingReferences) > 0 {
			fmt.Fprintf(stdout, "%s%d specifications missing test references%s\n", colorRed, len(missingReferences), colorReset)
		}
		
		if len(testsNotFound) > 0 {
			fmt.Fprintf(stdout, "%s%d test references not found%s\n", colorRed, len(testsNotFound), colorReset)
		}
		
		if len(interfaceErrors) > 0 {
			fmt.Fprintf(stdout, "%s%d interface implementation errors:%s\n", colorRed, len(interfaceErrors), colorReset)
			for impl, missing := range interfaceErrors {
				fmt.Fprintf(stdout, "  %s is missing: %s\n", impl, strings.Join(missing, ", "))
			}
		}
		
		return 1
	}
	
	fmt.Fprintf(stdout, "%sAll specifications covered ✓%s\n", colorGreen, colorReset)
	return 0
}

func loadSpecification(path string) (*spec.Specification, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	
	if !info.IsDir() {
		// Single file
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		
		specification, err := parser.ParseMarkdown(string(content))
		if err != nil {
			return nil, err
		}
		specification.FilePath = path
		
		return specification, nil
	}
	
	// Directory - use the new ParseDirectory function
	return parser.ParseDirectory(path)
}

// printSpecificationCollapsedWithErrors is a wrapper that passes interface errors
func printSpecificationCollapsedWithErrors(specification *spec.Specification, testSet map[string]bool, interfaceErrors map[string][]string, stdout io.Writer) {
	for _, section := range specification.Sections {
		printSectionCollapsedWithErrors(section, 0, testSet, interfaceErrors, stdout)
	}
}

// printSectionCollapsedWithErrors handles display with interface error checking
func printSectionCollapsedWithErrors(section *spec.Section, indent int, testSet map[string]bool, interfaceErrors map[string][]string, stdout io.Writer) {
	prefix := colorGray + strings.Repeat("· ", indent) + colorReset
	
	// Check if this section or any descendant has errors (including interface errors)
	sectionHasError := checkSectionHasErrorWithInterface(section, testSet, interfaceErrors)
	
	// Count total and passing specs in this section
	total, passing := countSectionCoverage(section, testSet)
	
	var statusIcon string
	if sectionHasError {
		statusIcon = fmt.Sprintf("%s✗%s", colorRed, colorReset)
	} else {
		statusIcon = fmt.Sprintf("%s✓%s", colorGreen, colorReset)
	}
	
	// Print section title
	fmt.Fprintf(stdout, "%s%s %s%s%s", prefix, statusIcon, colorBlue, section.Title, colorReset)
	
	if section.IsLeaf() {
		// Leaf sections show their test info only if they require tests
		if section.RequiresTest() {
			if !section.HasTest() {
				fmt.Fprintf(stdout, " %s(Missing test reference)%s\n", colorRed, colorReset)
			} else if !testSet[section.TestName] {
				fmt.Fprintf(stdout, " %s(Test not found: %s)%s\n", colorRed, section.TestName, colorReset)
			} else {
				fmt.Fprintf(stdout, " %s(%s)%s\n", colorGray, section.TestName, colorReset)
			}
		} else {
			// Interface leaf sections don't need tests
			fmt.Fprintln(stdout, "")
		}
		return
	}
	
	// For non-leaf sections:
	// - If it has errors, expand to show the problems
	// - If it has no errors, collapse and show count
	if sectionHasError {
		// Has errors - expand to show them
		fmt.Fprintln(stdout, "")
		for _, child := range section.Children {
			printSectionCollapsedWithErrors(child, indent+1, testSet, interfaceErrors, stdout)
		}
	} else {
		// No errors - collapse and show count
		fmt.Fprintf(stdout, " %s(%d/%d passed)%s\n", colorGray, passing, total, colorReset)
		// Children are not printed - they're collapsed
	}
}

// printSpecificationWithStatusAndErrors is a wrapper for verbose mode
func printSpecificationWithStatusAndErrors(specification *spec.Specification, testSet map[string]bool, interfaceErrors map[string][]string, stdout io.Writer) {
	for _, section := range specification.Sections {
		printSectionWithStatusAndErrors(section, 0, testSet, interfaceErrors, stdout)
	}
}

// printSectionWithStatusAndErrors handles verbose display with interface error checking
func printSectionWithStatusAndErrors(section *spec.Section, indent int, testSet map[string]bool, interfaceErrors map[string][]string, stdout io.Writer) {
	// Print indentation with middle dots (gray)
	prefix := colorGray + strings.Repeat("· ", indent) + colorReset
	
	// Determine section status (including interface errors)
	sectionHasError := checkSectionHasErrorWithInterface(section, testSet, interfaceErrors)
	var statusIcon string
	if sectionHasError {
		statusIcon = fmt.Sprintf("%s✗%s", colorRed, colorReset)
	} else {
		statusIcon = fmt.Sprintf("%s✓%s", colorGreen, colorReset)
	}
	
	// Print section title with status icon
	fmt.Fprintf(stdout, "%s%s %s%s%s", prefix, statusIcon, colorBlue, section.Title, colorReset)
	
	// Additional info for leaf sections
	if section.IsLeaf() {
		if section.RequiresTest() {
			if !section.HasTest() {
				// Missing test reference
				fmt.Fprintf(stdout, " %s(Missing test reference)%s\n", colorRed, colorReset)
			} else if !testSet[section.TestName] {
				// Test not found
				fmt.Fprintf(stdout, " %s(Test not found: %s)%s\n", colorRed, section.TestName, colorReset)
			} else {
				// Test found
				fmt.Fprintf(stdout, " %s(%s)%s\n", colorGray, section.TestName, colorReset)
			}
		} else {
			// Interface leaf sections don't need tests
			fmt.Fprintln(stdout, "")
		}
	} else {
		// Non-leaf section
		fmt.Fprintln(stdout, "")
	}
	
	// Print children recursively
	for _, child := range section.Children {
		printSectionWithStatusAndErrors(child, indent+1, testSet, interfaceErrors, stdout)
	}
}

// checkSectionHasErrorWithInterface checks for both test errors and interface errors
func checkSectionHasErrorWithInterface(section *spec.Section, testSet map[string]bool, interfaceErrors map[string][]string) bool {
	// Check if this is an implementation with validation errors
	if _, hasError := interfaceErrors[section.Title]; hasError {
		return true
	}
	
	// Check if this section itself has a test error
	if section.IsLeaf() && section.RequiresTest() {
		if !section.HasTest() || !testSet[section.TestName] {
			return true
		}
	}
	
	// Check all children recursively
	for _, child := range section.Children {
		if checkSectionHasErrorWithInterface(child, testSet, interfaceErrors) {
			return true
		}
	}
	
	return false
}

// Legacy functions kept for compatibility but now just call the new versions
func printSpecificationCollapsed(specification *spec.Specification, testSet map[string]bool, stdout io.Writer) {
	printSpecificationCollapsedWithErrors(specification, testSet, make(map[string][]string), stdout)
}

func printSectionCollapsed(section *spec.Section, indent int, testSet map[string]bool, stdout io.Writer) {
	printSectionCollapsedWithErrors(section, indent, testSet, make(map[string][]string), stdout)
}

func printSpecificationWithStatus(specification *spec.Specification, testSet map[string]bool, stdout io.Writer) {
	printSpecificationWithStatusAndErrors(specification, testSet, make(map[string][]string), stdout)
}

func printSectionWithStatus(section *spec.Section, indent int, testSet map[string]bool, stdout io.Writer) {
	printSectionWithStatusAndErrors(section, indent, testSet, make(map[string][]string), stdout)
}

// countSectionCoverage counts total and passing specs in a section
func countSectionCoverage(section *spec.Section, testSet map[string]bool) (total int, passing int) {
	if section.IsLeaf() && section.RequiresTest() {
		total = 1
		if section.HasTest() && testSet[section.TestName] {
			passing = 1
		}
		return
	}
	
	for _, child := range section.Children {
		childTotal, childPassing := countSectionCoverage(child, testSet)
		total += childTotal
		passing += childPassing
	}
	return
}

// checkSectionHasError returns true if this section or any of its descendants has an error
func checkSectionHasError(section *spec.Section, testSet map[string]bool) bool {
	// Check if this section itself has an error
	if section.IsLeaf() && section.RequiresTest() {
		if !section.HasTest() || !testSet[section.TestName] {
			return true
		}
	}
	
	// Check all children recursively
	for _, child := range section.Children {
		if checkSectionHasError(child, testSet) {
			return true
		}
	}
	
	return false
}