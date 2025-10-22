package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Alge/aligned/internal/spec"
)

// ParseDirectory parses all .md files in a directory recursively and builds
// a unified specification tree based on the directory structure
func ParseDirectory(rootPath string) (*spec.Specification, error) {
	// Build the tree structure from filesystem
	rootSection, err := buildDirectoryTree(rootPath, rootPath, -1) // Start at -1 so root files are at level 0
	if err != nil {
		return nil, err
	}
	
	// The root section is just a container - return its children as the root sections
	if rootSection != nil && rootSection.Title == "" {
		return &spec.Specification{
			Sections: rootSection.Children,
		}, nil
	}
	
	// If we have content at root level, include it
	if rootSection != nil {
		return &spec.Specification{
			Sections: []*spec.Section{rootSection},
		}, nil
	}
	
	return &spec.Specification{
		Sections: []*spec.Section{},
	}, nil
}

// buildDirectoryTree recursively builds a section tree from a directory
func buildDirectoryTree(dirPath string, rootPath string, level int) (*spec.Section, error) {
	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	
	// Separate files and directories
	var mdFiles []os.DirEntry
	var subDirs []os.DirEntry
	var hasDirnameMd bool
	
	dirName := filepath.Base(dirPath)
	dirnameMdFile := dirName + ".md"
	
	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry)
		} else if strings.HasSuffix(entry.Name(), ".md") {
			mdFiles = append(mdFiles, entry)
			if entry.Name() == dirnameMdFile {
				hasDirnameMd = true
			}
		}
	}
	
	var parentSection *spec.Section
	var childSections []*spec.Section
	
	// Process dirname.md if it exists
	if hasDirnameMd {
		content, err := os.ReadFile(filepath.Join(dirPath, dirnameMdFile))
		if err != nil {
			return nil, err
		}
		
		parsed, err := ParseMarkdown(string(content))
		if err != nil {
			return nil, err
		}
		
		// Use the first section as parent
		if len(parsed.Sections) > 0 {
			parentSection = parsed.Sections[0]
			parentSection.Level = level + 1
			
			// Add remaining sections as children
			if len(parsed.Sections) > 1 {
				for _, s := range parsed.Sections[1:] {
					adjustSectionLevels(s, level+2)
					childSections = append(childSections, s)
				}
			}
			
			// Adjust children from the first section
			for _, child := range parentSection.Children {
				adjustSectionLevels(child, level+2)
			}
		}
	}
	
	// Process other .md files
	for _, file := range mdFiles {
		if hasDirnameMd && file.Name() == dirnameMdFile {
			continue // Already processed
		}
		
		content, err := os.ReadFile(filepath.Join(dirPath, file.Name()))
		if err != nil {
			return nil, err
		}
		
		parsed, err := ParseMarkdown(string(content))
		if err != nil {
			return nil, err
		}
		
		// Add all sections
		for _, section := range parsed.Sections {
			// For root directory files, keep level 1
			// For subdirectory files without parent, use level+1
			// For subdirectory files with parent, use level+2
			if dirPath == rootPath {
				// Root level files keep their original level
			} else if parentSection == nil {
				adjustSectionLevels(section, level+1)
			} else {
				adjustSectionLevels(section, level+2)
			}
			childSections = append(childSections, section)
		}
	}
	
	// Process subdirectories
	for _, subDir := range subDirs {
		subSection, err := buildDirectoryTree(filepath.Join(dirPath, subDir.Name()), rootPath, level+1)
		if err != nil {
			return nil, err
		}
		if subSection != nil {
			childSections = append(childSections, subSection)
		}
	}
	
	// If no parent section was created from dirname.md, create one from directory name
	if parentSection == nil && len(childSections) > 0 {
		// Don't create a section for the root directory
		if dirPath == rootPath {
			// Return a container section that will be unwrapped
			return &spec.Section{
				Children: childSections,
			}, nil
		}
		
		parentSection = &spec.Section{
			Title:    ConvertSnakeCaseToTitleCase(dirName),
			Level:    level + 1,
			Children: childSections,
		}
		
		// Set parent references
		for _, child := range childSections {
			child.Parent = parentSection
		}
	} else if parentSection != nil {
		// Append children to existing parent
		parentSection.Children = append(parentSection.Children, childSections...)
		
		// Set parent references
		for _, child := range childSections {
			child.Parent = parentSection
		}
	} else if dirPath == rootPath && len(childSections) == 0 {
		// Empty root directory
		return nil, nil
	}
	
	return parentSection, nil
}

// adjustSectionLevels recursively adjusts the levels of a section and its children
func adjustSectionLevels(section *spec.Section, newLevel int) {
	section.Level = newLevel
	for _, child := range section.Children {
		adjustSectionLevels(child, newLevel+1)
	}
}

// MergeSpecifications merges multiple specifications under a parent section
func MergeSpecifications(specs []*spec.Specification, parentSection *spec.Section) *spec.Specification {
	if parentSection == nil {
		// If no parent, just combine all sections
		var allSections []*spec.Section
		for _, s := range specs {
			allSections = append(allSections, s.Sections...)
		}
		return &spec.Specification{
			Sections: allSections,
		}
	}
	
	// Add all sections as children of parent
	for _, s := range specs {
		for _, section := range s.Sections {
			// Adjust levels
			adjustSectionLevels(section, parentSection.Level+1)
			section.Parent = parentSection
			parentSection.Children = append(parentSection.Children, section)
		}
	}
	
	return &spec.Specification{
		Sections: []*spec.Section{parentSection},
	}
}

// ConvertSnakeCaseToTitleCase converts snake_case strings to Title Case
func ConvertSnakeCaseToTitleCase(input string) string {
	if input == "" {
		return ""
	}
	
	// Split by underscore
	parts := strings.Split(input, "_")
	
	// Capitalize each part
	for i, part := range parts {
		if len(part) > 0 {
			// Handle acronyms that are already uppercase
			if strings.ToUpper(part) == part {
				parts[i] = strings.Title(strings.ToLower(part))
			} else {
				parts[i] = strings.Title(part)
			}
		}
	}
	
	return strings.Join(parts, " ")
}