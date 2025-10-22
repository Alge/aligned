# Init command

## Display help when called without parameters

The `align init` command without parameters displays usage information and lists all supported connectors.

**Test:** `Alge/aligned/cmd/align.TestInitNoArgsShowsHelp`

## Display help when called with help parameter

The `align init help` command displays the same usage information and list of supported connectors as calling `align init` without parameters.

**Test:** `Alge/aligned/cmd/align.TestInitHelpShowsHelp`

## Create configuration file

The `align init <language-framework> <path>` command creates a new .align.yml file.

**Test:** `Alge/aligned/cmd/align.TestInitCreatesFile`

## Write connector type to configuration

The created .align.yml file contains the specified connector type.

**Test:** `Alge/aligned/cmd/align.TestInitWritesConnectorType`

## Write path to configuration

The created .align.yml file contains the specified path.

**Test:** `Alge/aligned/cmd/align.TestInitWritesPath`

## Display success message

The init command displays a success message when .align.yml is created.

**Test:** `Alge/aligned/cmd/align.TestInitSuccessMessage`

## Exit with success code

The init command exits with code 0 when successful.

**Test:** `Alge/aligned/cmd/align.TestInitExitCode`

## Exit with error if config file already exists

The init command exits with code 1 if .align.yml already exists to prevent overwriting.

**Test:** `Alge/aligned/cmd/align.TestInitConfigExists`

## Exit with error if arguments missing

The init command exits with code 1 if connector type or path is not provided.

**Test:** `Alge/aligned/cmd/align.TestInitMissingArguments`

## Exit with error if connector type not supported

The init command exits with code 1 if the specified connector type is not supported.

**Test:** `Alge/aligned/cmd/align.TestInitUnsupportedConnector`
