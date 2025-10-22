# Check command

## Core Validation

### Validate and report success when all specs covered

The check command validates that all leaf specifications have corresponding tests, displays a success message, and exits with code 0 when all specifications are covered.

**Test:** `Alge/aligned/cmd/align.TestCheckSuccessCase`

### Report missing test references and exit with error

The check command exits with code 1 and reports which specifications are missing test references.

**Test:** `Alge/aligned/cmd/align.TestCheckMissingTestReferences`

### Report tests not found and exit with error

The check command exits with code 1 and reports which test references in the spec cannot be found in discovered tests.

**Test:** `Alge/aligned/cmd/align.TestCheckTestsNotFound`

## File Loading

### Load and check single specification file

The `align check <file>` command loads and validates a single specification file.

**Test:** `Alge/aligned/cmd/align.TestCheckLoadsSingleFile`

### Load directory recursively

The `align check <directory>` command recursively loads all .md files in a directory.

**Test:** `Alge/aligned/cmd/align.TestCheckLoadsDirectory`

## Output Formatting

### Collapse successful sections with summary counts

The check command collapses fully successful sections into a single line displaying "(X/X passed)" to indicate passed specifications.

**Test:** `Alge/aligned/cmd/align.TestCheckCollapsedOutput`

### Expand failing sections automatically

Sections containing failures are automatically expanded to show failing specifications.

**Test:** `Alge/aligned/cmd/align.TestCheckExpandsFailures`

### Display full tree with verbose flag

The `align check -v` command displays the full specification tree with all sections expanded.

**Test:** `Alge/aligned/cmd/align.TestCheckVerboseOutput`

## Interface Validation Reporting

### Validate interfaces and report errors

The check command exits with code 1 when implementations are missing required interface sections, reports which sections are missing, and indicates interface validation status in the output.

**Test:** `Alge/aligned/cmd/align.TestCheckInterfaceValidation`
