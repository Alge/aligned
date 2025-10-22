// internal/spec/spec_test.go
package spec

import (
	"testing"
)

func TestSection_IsLeaf(t *testing.T) {
	tests := []struct {
		name     string
		section  *Section
		expected bool
	}{
		{
			name: "section with no children is leaf",
			section: &Section{
				Title:    "Test Section",
				Children: nil,
			},
			expected: true,
		},
		{
			name: "section with empty children slice is leaf",
			section: &Section{
				Title:    "Test Section",
				Children: []*Section{},
			},
			expected: true,
		},
		{
			name: "section with children is not leaf",
			section: &Section{
				Title: "Parent Section",
				Children: []*Section{
					{Title: "Child"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.section.IsLeaf()
			if result != tt.expected {
				t.Errorf("IsLeaf() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSection_HasTest(t *testing.T) {
	tests := []struct {
		name     string
		section  *Section
		expected bool
	}{
		{
			name: "section with test name has test",
			section: &Section{
				TestName: "TestSomething",
			},
			expected: true,
		},
		{
			name: "section with empty test name has no test",
			section: &Section{
				TestName: "",
			},
			expected: false,
		},
		{
			name:     "section with nil test name has no test",
			section:  &Section{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.section.HasTest()
			if result != tt.expected {
				t.Errorf("HasTest() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSpecification_AllLeaves(t *testing.T) {
	tests := []struct {
		name           string
		spec           *Specification
		expectedCount  int
		expectedTitles []string
	}{
		{
			name: "flat structure with all leaves",
			spec: &Specification{
				Sections: []*Section{
					{Title: "Section 1", Children: nil},
					{Title: "Section 2", Children: nil},
					{Title: "Section 3", Children: nil},
				},
			},
			expectedCount:  3,
			expectedTitles: []string{"Section 1", "Section 2", "Section 3"},
		},
		{
			name: "nested structure",
			spec: &Specification{
				Sections: []*Section{
					{
						Title: "Parent 1",
						Children: []*Section{
							{Title: "Child 1.1", Children: nil},
							{Title: "Child 1.2", Children: nil},
						},
					},
					{
						Title: "Parent 2",
						Children: []*Section{
							{Title: "Child 2.1", Children: nil},
						},
					},
				},
			},
			expectedCount:  3,
			expectedTitles: []string{"Child 1.1", "Child 1.2", "Child 2.1"},
		},
		{
			name: "deeply nested structure",
			spec: &Specification{
				Sections: []*Section{
					{
						Title: "Level 1",
						Children: []*Section{
							{
								Title: "Level 2",
								Children: []*Section{
									{Title: "Level 3 Leaf", Children: nil},
								},
							},
						},
					},
				},
			},
			expectedCount:  1,
			expectedTitles: []string{"Level 3 Leaf"},
		},
		{
			name: "empty specification",
			spec: &Specification{
				Sections: []*Section{},
			},
			expectedCount:  0,
			expectedTitles: []string{},
		},
		{
			name: "parent with no children still counts as leaf",
			spec: &Specification{
				Sections: []*Section{
					{Title: "Lone Section", Children: []*Section{}},
				},
			},
			expectedCount:  1,
			expectedTitles: []string{"Lone Section"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaves := tt.spec.AllLeaves()

			if len(leaves) != tt.expectedCount {
				t.Errorf("AllLeaves() returned %d leaves, want %d", len(leaves), tt.expectedCount)
			}

			// Check titles match
			for i, expectedTitle := range tt.expectedTitles {
				if i >= len(leaves) {
					t.Errorf("Missing leaf at index %d", i)
					continue
				}
				if leaves[i].Title != expectedTitle {
					t.Errorf("leaves[%d].Title = %q, want %q", i, leaves[i].Title, expectedTitle)
				}
			}
		})
	}
}

func TestSpecification_RequiredTests(t *testing.T) {
	tests := []struct {
		name          string
		spec          *Specification
		expectedTests []string
	}{
		{
			name: "all leaves have tests",
			spec: &Specification{
				Sections: []*Section{
					{Title: "Spec 1", TestName: "TestOne", Children: nil},
					{Title: "Spec 2", TestName: "TestTwo", Children: nil},
					{Title: "Spec 3", TestName: "TestThree", Children: nil},
				},
			},
			expectedTests: []string{"TestOne", "TestTwo", "TestThree"},
		},
		{
			name: "some leaves missing tests",
			spec: &Specification{
				Sections: []*Section{
					{Title: "Spec 1", TestName: "TestOne", Children: nil},
					{Title: "Spec 2", TestName: "", Children: nil}, // No test
					{Title: "Spec 3", TestName: "TestThree", Children: nil},
				},
			},
			expectedTests: []string{"TestOne", "TestThree"},
		},
		{
			name: "nested structure with tests only on leaves",
			spec: &Specification{
				Sections: []*Section{
					{
						Title:    "Parent",
						TestName: "", // Parent has no test (shouldn't be included anyway)
						Children: []*Section{
							{Title: "Child 1", TestName: "TestChild1", Children: nil},
							{Title: "Child 2", TestName: "TestChild2", Children: nil},
						},
					},
				},
			},
			expectedTests: []string{"TestChild1", "TestChild2"},
		},
		{
			name: "no tests specified",
			spec: &Specification{
				Sections: []*Section{
					{Title: "Spec 1", TestName: "", Children: nil},
					{Title: "Spec 2", TestName: "", Children: nil},
				},
			},
			expectedTests: []string{},
		},
		{
			name: "empty specification",
			spec: &Specification{
				Sections: []*Section{},
			},
			expectedTests: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tests := tt.spec.RequiredTests()

			if len(tests) != len(tt.expectedTests) {
				t.Errorf("RequiredTests() returned %d tests, want %d", len(tests), len(tt.expectedTests))
			}

			// Convert to map for easy lookup
			testMap := make(map[string]bool)
			for _, test := range tests {
				testMap[test] = true
			}

			// Check all expected tests are present
			for _, expectedTest := range tt.expectedTests {
				if !testMap[expectedTest] {
					t.Errorf("RequiredTests() missing expected test %q", expectedTest)
				}
			}
		})
	}
}

func TestSection_ParentChildRelationships(t *testing.T) {
	parent := &Section{
		Title: "Parent",
	}

	child1 := &Section{
		Title:  "Child 1",
		Parent: parent,
	}

	child2 := &Section{
		Title:  "Child 2",
		Parent: parent,
	}

	parent.Children = []*Section{child1, child2}

	// Test parent has correct children
	if len(parent.Children) != 2 {
		t.Errorf("parent has %d children, want 2", len(parent.Children))
	}

	// Test children have correct parent
	if child1.Parent != parent {
		t.Errorf("child1.Parent is not parent")
	}
	if child2.Parent != parent {
		t.Errorf("child2.Parent is not parent")
	}

	// Test parent is not a leaf
	if parent.IsLeaf() {
		t.Errorf("parent.IsLeaf() = true, want false")
	}

	// Test children are leaves
	if !child1.IsLeaf() {
		t.Errorf("child1.IsLeaf() = false, want true")
	}
	if !child2.IsLeaf() {
		t.Errorf("child2.IsLeaf() = false, want true")
	}
}

func TestSection_RequiresTest(t *testing.T) {
	tests := []struct {
		name     string
		section  *Section
		expected bool
	}{
		{
			name: "leaf section requires test",
			section: &Section{
				Title:    "Feature",
				Children: nil,
			},
			expected: true,
		},
		{
			name: "parent section does not require test",
			section: &Section{
				Title: "Parent",
				Children: []*Section{
					{Title: "Child"},
				},
			},
			expected: false,
		},
		{
			name: "interface section does not require test",
			section: &Section{
				Title:    "API [INTERFACE]",
				Children: nil,
			},
			expected: false,
		},
		{
			name: "child of interface does not require test",
			section: func() *Section {
				parent := &Section{
					Title:    "API [INTERFACE]",
					Children: nil,
				}
				child := &Section{
					Title:  "Method",
					Parent: parent,
				}
				parent.Children = []*Section{child}
				return child
			}(),
			expected: false,
		},
		{
			name: "implementation section requires test",
			section: &Section{
				Title:    "Go Implementation [IMPLEMENTS: API]",
				Children: nil,
			},
			expected: true,
		},
		{
			name: "deeply nested child of interface does not require test",
			section: func() *Section {
				grandparent := &Section{
					Title: "API [INTERFACE]",
				}
				parent := &Section{
					Title:  "Group",
					Parent: grandparent,
				}
				child := &Section{
					Title:  "Method",
					Parent: parent,
				}
				grandparent.Children = []*Section{parent}
				parent.Children = []*Section{child}
				return child
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.section.RequiresTest()
			if result != tt.expected {
				t.Errorf("RequiresTest() = %v, want %v", result, tt.expected)
			}
		})
	}
}
