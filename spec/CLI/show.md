# Show command

## Parse specification file

The `align show <spec-file>` command parses the specification file.

**Test:** `Alge/aligned/cmd/align.TestShowParses`

## Display section titles

The show output includes all section titles from the specification.

**Test:** `Alge/aligned/cmd/align.TestShowDisplaysTitles`

## Display hierarchical structure

The show output displays the hierarchical structure of sections.

**Test:** `Alge/aligned/cmd/align.TestShowDisplaysHierarchy`

## Exit with success code

The show command exits with code 0 when successful.

**Test:** `Alge/aligned/cmd/align.TestShowExitCode`

## Exit with error when file not found

The show command exits with code 1 when the specified file does not exist.

**Test:** `Alge/aligned/cmd/align.TestShowFileNotFound`

## Handle empty specification files

The show command successfully processes empty specification files (exit code 0).

**Test:** `Alge/aligned/cmd/align.TestShowEmptyFile`
