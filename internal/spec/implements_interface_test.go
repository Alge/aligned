package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectImplementationMarker(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		isImpl       bool
		interfaceName string
	}{
		{
			name:         "basic implementation marker",
			title:        "Go Connector [IMPLEMENTS: Connector Interface]",
			isImpl:       true,
			interfaceName: "Connector Interface",
		},
		{
			name:         "implementation with extra spaces",
			title:        "Python Connector [IMPLEMENTS:   Framework Connector  ]",
			isImpl:       true,
			interfaceName: "Framework Connector",
		},
		{
			name:         "regular section",
			title:        "Regular Section",
			isImpl:       false,
			interfaceName: "",
		},
		{
			name:         "interface section not implementation",
			title:        "Base Interface [INTERFACE]",
			isImpl:       false,
			interfaceName: "",
		},
		{
			name:         "section with implements word but not marker",
			title:        "This implements the pattern",
			isImpl:       false,
			interfaceName: "",
		},
		{
			name:         "multi-word interface name",
			title:        "MyImpl [IMPLEMENTS: Complex Multi Word Interface]",
			isImpl:       true,
			interfaceName: "Complex Multi Word Interface",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := &Section{Title: tt.title}
			interfaceName := section.GetImplementedInterface()
			
			if tt.isImpl {
				assert.NotEmpty(t, interfaceName, "Should detect implementation marker")
				assert.Equal(t, tt.interfaceName, interfaceName, "Should extract correct interface name")
			} else {
				assert.Empty(t, interfaceName, "Should not detect implementation marker")
			}
		})
	}
}

func TestExtractInterfaceName(t *testing.T) {
	// This test focuses specifically on the extraction logic
	testCases := []struct {
		title         string
		expected      string
	}{
		{
			title:    "Go [IMPLEMENTS: Test Framework]",
			expected: "Test Framework",
		},
		{
			title:    "Implementation [IMPLEMENTS: API Contract]",
			expected: "API Contract",
		},
		{
			title:    "Service [IMPLEMENTS: Base Service Interface]",
			expected: "Base Service Interface",
		},
		{
			title:    "No marker here",
			expected: "",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			section := &Section{Title: tc.title}
			result := section.GetImplementedInterface()
			assert.Equal(t, tc.expected, result)
		})
	}
}