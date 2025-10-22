package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateImplementationStructure(t *testing.T) {
	t.Run("valid implementation with all required sections", func(t *testing.T) {
		// Define interface
		iface := &Section{
			Title: "Connector Interface [INTERFACE]",
			Children: []*Section{
				{Title: "Detect framework presence"},
				{Title: "Discover tests"},
				{Title: "Parse test output"},
			},
		}
		
		// Valid implementation (has all required sections)
		impl := &Section{
			Title: "Go Connector [IMPLEMENTS: Connector Interface]",
			Children: []*Section{
				{Title: "Detect framework presence"},
				{Title: "Discover tests"},
				{Title: "Parse test output"},
			},
		}
		
		missing := ValidateImplementation(impl, iface)
		assert.Empty(t, missing, "Should have no missing sections")
	})
	
	t.Run("invalid implementation missing sections", func(t *testing.T) {
		// Define interface
		iface := &Section{
			Title: "Connector Interface [INTERFACE]",
			Children: []*Section{
				{Title: "Detect framework presence"},
				{Title: "Discover tests"},
				{Title: "Parse test output"},
			},
		}
		
		// Invalid implementation (missing some sections)
		impl := &Section{
			Title: "Python Connector [IMPLEMENTS: Connector Interface]",
			Children: []*Section{
				{Title: "Detect framework presence"},
				// Missing "Discover tests" and "Parse test output"
			},
		}
		
		missing := ValidateImplementation(impl, iface)
		assert.Len(t, missing, 2, "Should have 2 missing sections")
		// The comparison should be case-insensitive
		assert.Contains(t, missing, "discover tests")
		assert.Contains(t, missing, "parse test output")
	})
	
	t.Run("implementation with different casing should match", func(t *testing.T) {
		iface := &Section{
			Title: "API [INTERFACE]",
			Children: []*Section{
				{Title: "Get Data"},
			},
		}
		
		impl := &Section{
			Title: "REST API [IMPLEMENTS: API]",
			Children: []*Section{
				{Title: "GET DATA"}, // Different casing
			},
		}
		
		missing := ValidateImplementation(impl, iface)
		assert.Empty(t, missing, "Should match despite different casing")
	})
}

func TestAllowAdditionalSectionsInImplementation(t *testing.T) {
	// Interface with minimal requirements
	iface := &Section{
		Title: "Basic Interface [INTERFACE]",
		Children: []*Section{
			{Title: "Required Method"},
		},
	}
	
	// Implementation with extra sections
	impl := &Section{
		Title: "Extended Implementation [IMPLEMENTS: Basic Interface]",
		Children: []*Section{
			{Title: "Required Method"},
			{Title: "Additional Method"},  // Extra - this is allowed
			{Title: "Another Extra"},      // Another extra - also allowed
		},
	}
	
	missing := ValidateImplementation(impl, iface)
	assert.Empty(t, missing, "Extra sections should be allowed")
}