package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectInterfaceMarker(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected bool
	}{
		{
			name:     "section with INTERFACE marker",
			title:    "Connector Interface [INTERFACE]",
			expected: true,
		},
		{
			name:     "section with interface marker in different position",
			title:    "Some Section [INTERFACE]",
			expected: true,
		},
		{
			name:     "regular section without marker",
			title:    "Regular Section",
			expected: false,
		},
		{
			name:     "section with IMPLEMENTS marker should not be interface",
			title:    "Implementation [IMPLEMENTS: Something]",
			expected: false,
		},
		{
			name:     "section with interface word but not marker",
			title:    "Interface Design Patterns",
			expected: false,
		},
		{
			name:     "case sensitivity check",
			title:    "Section [interface]",
			expected: false, // Should be case-sensitive
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := &Section{Title: tt.title}
			result := section.IsInterface()
			assert.Equal(t, tt.expected, result, "IsInterface() for title '%s'", tt.title)
		})
	}
}