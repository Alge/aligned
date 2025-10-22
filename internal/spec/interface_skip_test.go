package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceSkipsTestRequirements(t *testing.T) {
	t.Run("interface leaf section should not require test", func(t *testing.T) {
		interfaceSection := &Section{
			Title: "Connector Interface [INTERFACE]",
			Children: []*Section{
				{
					Title:    "Detect framework",
					Children: []*Section{}, // Leaf section
				},
			},
		}
		
		// Set parent relationship
		interfaceSection.Children[0].Parent = interfaceSection
		
		// Interface child should not require test
		assert.False(t, interfaceSection.Children[0].RequiresTest())
	})
	
	t.Run("regular leaf section should require test", func(t *testing.T) {
		regularSection := &Section{
			Title:    "Regular Feature",
			Children: []*Section{}, // Leaf section
		}
		
		assert.True(t, regularSection.RequiresTest())
	})
	
	t.Run("nested interface children should not require tests", func(t *testing.T) {
		interfaceSection := &Section{
			Title: "API Interface [INTERFACE]",
			Children: []*Section{
				{
					Title: "Authentication",
					Children: []*Section{
						{
							Title:    "Login method",
							Children: []*Section{}, // Deeply nested leaf
						},
					},
				},
			},
		}
		
		// Set parent relationships
		interfaceSection.Children[0].Parent = interfaceSection
		interfaceSection.Children[0].Children[0].Parent = interfaceSection.Children[0]
		
		// Even deeply nested children of interfaces should not require tests
		assert.False(t, interfaceSection.Children[0].Children[0].RequiresTest())
	})
	
	t.Run("non-leaf sections never require tests", func(t *testing.T) {
		section := &Section{
			Title: "Parent Section",
			Children: []*Section{
				{Title: "Child"},
			},
		}
		
		// Non-leaf sections don't require tests regardless
		assert.False(t, section.RequiresTest())
	})
}