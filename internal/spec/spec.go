package spec

import "strings"

// Specification represents a parsed specification document
type Specification struct {
	FilePath string
	Sections []*Section
}

// Section represents a section in the specification
type Section struct {
	Level    int        // Heading level (1, 2, 3...)
	Number   string     // Auto-generated: "1.1.1"
	Title    string     // "Parse Markdown headings"
	Content  string     // Everything between this heading and next
	TestName string     // "TestParseMarkdownHeadings" (empty if not a leaf)
	Children []*Section // Nested sections
	Parent   *Section   // Parent section (nil for root)
}

// IsLeaf returns true if this section has no children
func (s *Section) IsLeaf() bool {
	return len(s.Children) == 0
}

// HasTest returns true if this section has a test reference
func (s *Section) HasTest() bool {
	return s.TestName != ""
}

// AllLeaves returns all leaf sections in the tree
func (s *Specification) AllLeaves() []*Section {
	var leaves []*Section
	var walk func(*Section)
	walk = func(section *Section) {
		if section.IsLeaf() {
			leaves = append(leaves, section)
		}
		for _, child := range section.Children {
			walk(child)
		}
	}
	for _, root := range s.Sections {
		walk(root)
	}
	return leaves
}

// RequiredTests returns all test names that should exist
func (s *Specification) RequiredTests() []string {
	var tests []string
	for _, leaf := range s.AllLeaves() {
		if leaf.HasTest() {
			tests = append(tests, leaf.TestName)
		}
	}
	return tests
}

func (s *Section) IsInterface() bool {
	return strings.Contains(s.Title, "[INTERFACE]")
}

// RequiresTest returns true if this section requires a test reference
// Interfaces and their children don't require tests
func (s *Section) RequiresTest() bool {
	// Only leaf sections normally require tests
	if !s.IsLeaf() {
		return false
	}
	
	// Check if this section or any ancestor is an interface
	current := s
	for current != nil {
		if current.IsInterface() {
			return false
		}
		current = current.Parent
	}
	
	return true
}

func (s *Section) GetImplementedInterface() string {
	// Look for [IMPLEMENTS: InterfaceName] pattern
	if idx := strings.Index(s.Title, "[IMPLEMENTS:"); idx != -1 {
		// Find the closing bracket
		endIdx := strings.Index(s.Title[idx:], "]")
		if endIdx != -1 {
			// Extract the interface name between : and ]
			startIdx := idx + len("[IMPLEMENTS:")
			interfaceName := s.Title[startIdx : idx+endIdx]
			// Trim any whitespace
			return strings.TrimSpace(interfaceName)
		}
	}
	return ""
}

// ValidateImplementation checks if an implementation has all required sections from the interface
// Returns a list of missing section titles (normalized to lowercase)
func ValidateImplementation(impl *Section, iface *Section) []string {
	var missing []string
	
	// Get all child titles from interface (normalized to lowercase)
	requiredSections := make(map[string]bool)
	for _, child := range iface.Children {
		normalized := strings.ToLower(strings.TrimSpace(child.Title))
		requiredSections[normalized] = true
	}
	
	// Get all child titles from implementation (normalized to lowercase)
	implementedSections := make(map[string]bool)
	for _, child := range impl.Children {
		normalized := strings.ToLower(strings.TrimSpace(child.Title))
		implementedSections[normalized] = true
	}
	
	// Find missing sections
	for required := range requiredSections {
		if !implementedSections[required] {
			missing = append(missing, required)
		}
	}
	
	return missing
}

// ValidateInterfaces checks all implementations against their interfaces
// Returns a map of implementation titles to their missing sections
func (s *Specification) ValidateInterfaces() map[string][]string {
	errors := make(map[string][]string)
	
	// Collect all interfaces and implementations
	interfaces := make(map[string]*Section)
	implementations := make(map[*Section]string)
	
	var walk func(*Section)
	walk = func(section *Section) {
		// Check if it's an interface
		if section.IsInterface() {
			// Extract clean name without [INTERFACE] marker
			name := strings.TrimSpace(section.Title)
			if idx := strings.Index(name, "[INTERFACE]"); idx != -1 {
				name = strings.TrimSpace(name[:idx])
			}
			interfaces[name] = section
		}
		
		// Check if it's an implementation
		if implInterface := section.GetImplementedInterface(); implInterface != "" {
			implementations[section] = implInterface
		}
		
		// Recurse through children
		for _, child := range section.Children {
			walk(child)
		}
	}
	
	// Walk all sections
	for _, root := range s.Sections {
		walk(root)
	}
	
	// Validate each implementation
	for impl, interfaceName := range implementations {
		iface, found := interfaces[interfaceName]
		if !found {
			errors[impl.Title] = []string{"interface '" + interfaceName + "' not found"}
			continue
		}
		
		// Check for missing sections
		missing := ValidateImplementation(impl, iface)
		if len(missing) > 0 {
			errors[impl.Title] = missing
		}
	}
	
	return errors
}