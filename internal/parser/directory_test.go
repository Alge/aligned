package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Alge/aligned/internal/spec"
	"github.com/stretchr/testify/assert"
)

func TestBuildSpecTreeFromDirectory(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	
	// Create spec files
	mainSpec := `# Main Project

## Overview
Main project overview.

**Test:** ` + "`testproject.TestOverview`"
	
	submoduleSpec := `# Submodule

## Feature
Submodule feature.

**Test:** ` + "`testproject.TestSubmoduleFeature`"
	
	// Write files
	err := os.WriteFile(filepath.Join(tempDir, "main.md"), []byte(mainSpec), 0644)
	assert.NoError(t, err)
	
	subDir := filepath.Join(tempDir, "submodule")
	err = os.MkdirAll(subDir, 0755)
	assert.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(subDir, "feature.md"), []byte(submoduleSpec), 0644)
	assert.NoError(t, err)
	
	// Parse directory
	specification, err := ParseDirectory(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, specification)
	
	// Debug: print what we actually got
	t.Logf("Number of root sections: %d", len(specification.Sections))
	for i, s := range specification.Sections {
		t.Logf("Section %d: %s", i, s.Title)
		for j, c := range s.Children {
			t.Logf("  Child %d: %s", j, c.Title)
		}
	}
	
	// The root directory doesn't create a section, so we should have:
	// - Main Project (from main.md)
	// - Submodule (from the directory name)
	assert.Len(t, specification.Sections, 2)
	
	// Find Main Project
	var mainSection *spec.Section
	for _, s := range specification.Sections {
		if s.Title == "Main Project" {
			mainSection = s
			break
		}
	}
	assert.NotNil(t, mainSection)
	if mainSection != nil {
		assert.Len(t, mainSection.Children, 1)
		assert.Equal(t, "Overview", mainSection.Children[0].Title)
	}
	
	// Find Submodule section (from directory)
	var submoduleSection *spec.Section
	for _, s := range specification.Sections {
		if s.Title == "Submodule" {
			submoduleSection = s
			break
		}
	}
	assert.NotNil(t, submoduleSection)
	if submoduleSection != nil {
		// The "Submodule" from feature.md should be a child
		assert.Len(t, submoduleSection.Children, 1)
		assert.Equal(t, "Submodule", submoduleSection.Children[0].Title)
	}
}

func TestDirectoryNameAsParent(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create connectors/ directory with two spec files
	connectorsDir := filepath.Join(tempDir, "connectors")
	err := os.MkdirAll(connectorsDir, 0755)
	assert.NoError(t, err)
	
	interfaceSpec := `# Interface

## Define Contract
The interface defines the contract.

**Test:** ` + "`testproject.TestInterface`"
	
	goSpec := `# Go Implementation

## Implement Contract
Go implements the contract.

**Test:** ` + "`testproject.TestGoImpl`"
	
	err = os.WriteFile(filepath.Join(connectorsDir, "interface.md"), []byte(interfaceSpec), 0644)
	assert.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(connectorsDir, "go.md"), []byte(goSpec), 0644)
	assert.NoError(t, err)
	
	// Parse directory
	specification, err := ParseDirectory(tempDir)
	assert.NoError(t, err)
	
	// Should have 1 root section: "Connectors" (from directory name)
	assert.Len(t, specification.Sections, 1)
	assert.Equal(t, "Connectors", specification.Sections[0].Title)
	
	// Should have 2 children: Interface and Go Implementation
	connectorsSection := specification.Sections[0]
	assert.Len(t, connectorsSection.Children, 2)
	
	var hasInterface, hasGo bool
	for _, child := range connectorsSection.Children {
		if child.Title == "Interface" {
			hasInterface = true
		}
		if child.Title == "Go Implementation" {
			hasGo = true
		}
	}
	assert.True(t, hasInterface)
	assert.True(t, hasGo)
}

func TestDirnameMdAsParent(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create cli/ directory with cli.md and subcommand files
	cliDir := filepath.Join(tempDir, "cli")
	err := os.MkdirAll(cliDir, 0755)
	assert.NoError(t, err)
	
	// cli.md acts as parent section
	cliParentSpec := `# Command Line Interface

The CLI provides commands for specification validation.`
	
	checkSpec := `# Check Command

## Validate Coverage
Check that all specs have tests.

**Test:** ` + "`testproject.TestCheckCommand`"
	
	showSpec := `# Show Command  

## Display Structure
Show the spec tree.

**Test:** ` + "`testproject.TestShowCommand`"
	
	err = os.WriteFile(filepath.Join(cliDir, "cli.md"), []byte(cliParentSpec), 0644)
	assert.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(cliDir, "check.md"), []byte(checkSpec), 0644)
	assert.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(cliDir, "show.md"), []byte(showSpec), 0644)
	assert.NoError(t, err)
	
	// Parse directory
	specification, err := ParseDirectory(tempDir)
	assert.NoError(t, err)
	
	// Should have 1 root section: "Command Line Interface" from cli/cli.md
	assert.Len(t, specification.Sections, 1)
	assert.Equal(t, "Command Line Interface", specification.Sections[0].Title)
	
	// Should have 2 children: Check Command and Show Command
	cliSection := specification.Sections[0]
	assert.Len(t, cliSection.Children, 2)
	
	var hasCheck, hasShow bool
	for _, child := range cliSection.Children {
		if child.Title == "Check Command" {
			hasCheck = true
		}
		if child.Title == "Show Command" {
			hasShow = true
		}
	}
	assert.True(t, hasCheck)
	assert.True(t, hasShow)
}

func TestMergeSpecifications(t *testing.T) {
	// Test merging multiple specs into one tree
	spec1 := &spec.Specification{
		Sections: []*spec.Section{
			{
				Title: "Feature A",
				Level: 1,
				Children: []*spec.Section{},
			},
		},
	}
	
	spec2 := &spec.Specification{
		Sections: []*spec.Section{
			{
				Title: "Feature B", 
				Level: 1,
				Children: []*spec.Section{},
			},
		},
	}
	
	parentSection := &spec.Section{
		Title: "Parent",
		Level: 1,
		Children: []*spec.Section{},
	}
	
	merged := MergeSpecifications([]*spec.Specification{spec1, spec2}, parentSection)
	
	assert.NotNil(t, merged)
	assert.Len(t, merged.Sections, 1)
	assert.Equal(t, "Parent", merged.Sections[0].Title)
	assert.Len(t, merged.Sections[0].Children, 2)
	
	// Check children
	var hasFeatureA, hasFeatureB bool
	for _, child := range merged.Sections[0].Children {
		if child.Title == "Feature A" {
			hasFeatureA = true
			assert.Equal(t, 2, child.Level) // Should be nested under parent
		}
		if child.Title == "Feature B" {
			hasFeatureB = true
			assert.Equal(t, 2, child.Level)
		}
	}
	assert.True(t, hasFeatureA)
	assert.True(t, hasFeatureB)
}

func TestConvertSnakeCaseToTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "test_framework_integrations",
			expected: "Test Framework Integrations",
		},
		{
			input:    "go_connector",
			expected: "Go Connector",
		},
		{
			input:    "simple",
			expected: "Simple",
		},
		{
			input:    "test_with_many_words_here",
			expected: "Test With Many Words Here",
		},
		{
			input:    "CLI_commands",
			expected: "Cli Commands",
		},
		{
			input:    "",
			expected: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ConvertSnakeCaseToTitleCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}