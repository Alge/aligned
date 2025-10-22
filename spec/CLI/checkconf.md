# Checkconf command

## Display success message when configuration is valid

The checkconf output includes a success message when configuration is valid.

**Test:** `Alge/aligned/cmd/align.TestCheckconfSuccessMessage`

## Exit with success code when configuration is valid

The checkconf command exits with code 0 when configuration is valid.

**Test:** `Alge/aligned/cmd/align.TestCheckconfExitCode`

## Display configuration details with verbose flag

The `align checkconf -v` command prints the loaded configuration details including connector type, executable, and path.

**Test:** `Alge/aligned/cmd/align.TestCheckconfVerbosePrintsDetails`

## Exit with success code when using verbose flag

The checkconf command exits with code 0 when using verbose flag with valid configuration.

**Test:** `Alge/aligned/cmd/align.TestCheckconfVerboseExitCode`

## Exit with error when config file missing

The checkconf command exits with code 1 when .align.yml is not found.

**Test:** `Alge/aligned/cmd/align.TestCheckconfMissingFile`

## Exit with error when config has invalid YAML

The checkconf command exits with code 1 when .align.yml contains malformed YAML.

**Test:** `Alge/aligned/cmd/align.TestCheckconfInvalidYAML`

## Exit with error when config has no connectors

The checkconf command exits with code 1 when .align.yml is empty or has no connectors defined.

**Test:** `Alge/aligned/cmd/align.TestCheckconfEmptyConfig`
