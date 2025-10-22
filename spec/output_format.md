# Output Format

The output formatting system provides visual feedback about specification coverage through colors, symbols, and hierarchical display.

## Visual Elements

### Display status with symbols

Use ✓ symbol for passing specifications and ✗ symbol for failing specifications to provide quick visual status indicators.

**Test:** `Alge/aligned/cmd/align.TestCheckExpandsFailures`

### Display summary counts for collapsed sections

Show coverage summary in format "(X/X passed)" where X indicates the number of passing specifications out of total specifications in a collapsed section.

**Test:** `Alge/aligned/cmd/align.TestCheckCollapsedOutput`

### Indent sections hierarchically

Display section hierarchy using middle dot (·) indentation, with one middle dot per nesting level to show the specification tree structure.

**Test:** `Alge/aligned/cmd/align.TestShowDisplaysHierarchy`

## Check Command Output Modes

### Display collapsed output by default

The check command collapses fully passing sections into summary lines and automatically expands sections containing failures to show details.

**Test:** `Alge/aligned/cmd/align.TestCheckCollapsedOutput`

### Expand all sections in verbose mode

The `align check -v` command displays all sections in expanded form regardless of pass/fail status, showing complete specification tree with test details.

**Test:** `Alge/aligned/cmd/align.TestCheckVerboseOutput`

### Display specification coverage report header

The check command outputs "Specification coverage report:" as a header before displaying the specification tree.

**Test:** `Alge/aligned/cmd/align.TestCheckSuccessCase`

### Display final success or error message

The check command displays "All specifications covered ✓" on success or error summaries (count + descriptions) on failure after the specification tree.

**Test:** `Alge/aligned/cmd/align.TestCheckSuccessCase`

## Show Command Output

### Display section titles in hierarchical structure

The show command displays all section titles in a tree structure with proper nesting levels.

**Test:** `Alge/aligned/cmd/align.TestShowDisplaysTitles`

### Display test references

The show command displays test references in the format "Test: test_name" beneath the section that references them.

**Test:** `Alge/aligned/cmd/align.TestShowDisplaysTitles`

## Error Reporting

### Report missing test references

The check command displays which specifications are missing test references in failing sections, showing "Missing test reference" indicators.

**Test:** `Alge/aligned/cmd/align.TestCheckMissingTestReferences`

### Report tests not found

The check command displays which test references cannot be found in discovered tests, showing "Test not found: test_name" indicators.

**Test:** `Alge/aligned/cmd/align.TestCheckTestsNotFound`

### Report interface validation errors

The check command displays interface implementation errors with the implementation name and list of missing required sections.

**Test:** `Alge/aligned/cmd/align.TestCheckInterfaceValidation`
