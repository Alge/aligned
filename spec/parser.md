# Parser

The parser converts Markdown specification files into the in-memory data model. It handles both single files and directory hierarchies, extracting heading structure, test references, and building the specification tree.

## Markdown Parsing

### Parse Markdown headings

Read a Markdown file and extract heading text and level (# = 1, ## = 2, ### = 3).

**Test:** `Alge/aligned/internal/parser.TestParseMarkdownHeadings`

### Extract test reference from specification

Find lines matching "**Test:** `test_name`" and extract the test name. Supports both backtick-wrapped and plain text formats, and handles fully qualified test names with package paths.

**Test:** `Alge/aligned/internal/parser.TestExtractTestReference`

## Directory-Based Hierarchy

### Build specification tree from directory structure

Load all .md files from a directory recursively and build a unified specification tree where directory structure determines the hierarchy.

**Test:** `Alge/aligned/internal/parser.TestBuildSpecTreeFromDirectory`

### Use directory name as parent section

When a directory has multiple .md files but no `dirname.md` file, use the directory name as the parent section title.

**Test:** `Alge/aligned/internal/parser.TestDirectoryNameAsParent`

### Use dirname.md as parent section

When `dirname/dirname.md` exists, use its content as the parent section for all other files in that directory.

**Test:** `Alge/aligned/internal/parser.TestDirnameMdAsParent`

### Merge multiple specifications into unified tree

Combine specifications from multiple files in a directory into a single tree structure with proper parent-child relationships.

**Test:** `Alge/aligned/internal/parser.TestMergeSpecifications`

### Convert snake_case to Title Case

Convert snake_case directory and file names to Title Case for section titles when using directory/file names as sections.

**Test:** `Alge/aligned/internal/parser.TestConvertSnakeCaseToTitleCase`
