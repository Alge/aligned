# List-tests command

## Discover tests using connectors

The `align list-tests` command uses configured connectors to discover tests.

**Test:** `Alge/aligned/cmd/align.TestListTestsDiscovery`

## Print discovered test names

The output includes all discovered test names.

**Test:** `Alge/aligned/cmd/align.TestListTestsPrintsNames`

## Exit with success code

The list-tests command exits with code 0 when successful.

**Test:** `Alge/aligned/cmd/align.TestListTestsExitCode`

## Exit with error when config file missing

The list-tests command exits with code 1 when .align.yml is not found.

**Test:** `Alge/aligned/cmd/align.TestListTestsConfigMissing`

## Exit with error when config file invalid

The list-tests command exits with code 1 when .align.yml is malformed or invalid.

**Test:** `Alge/aligned/cmd/align.TestListTestsConfigInvalid`
