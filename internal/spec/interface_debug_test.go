package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceParentChainDebug(t *testing.T) {
	// This mimics exactly what should happen with the Connector Interface
	parent := &Section{
		Title: "Connector Interface [INTERFACE]",
		Children: []*Section{},
	}
	
	child := &Section{
		Title:    "Detect framework presence",
		Children: []*Section{}, // This is a leaf
		Parent:   parent,
	}
	
	parent.Children = append(parent.Children, child)
	
	// Debug checks
	assert.True(t, parent.IsInterface(), "Parent should be detected as interface")
	assert.True(t, child.IsLeaf(), "Child should be a leaf")
	assert.False(t, child.RequiresTest(), "Child of interface should NOT require test")
}