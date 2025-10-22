# Help command

Users need to discover available commands and understand basic usage without consulting external documentation.

## Display help when no command provided

When align is invoked without any arguments, it displays the help message and exits with code 0. This provides the same output as `align help`.

**Test:** `Alge/aligned/cmd/align.TestNoArgsShowsHelp`

## Display help with help command

The `align help` command displays help information.

**Test:** `Alge/aligned/cmd/align.TestHelpCommand`

## Include tool description

The help output includes a description: "Validate that specifications are covered by tests"

**Test:** `Alge/aligned/cmd/align.TestHelpIncludesDescription`

## Show usage pattern

The help output shows the usage pattern: `align <command> [options]`

**Test:** `Alge/aligned/cmd/align.TestHelpShowsUsagePattern`

## Document check command

The help output includes the check command with usage `check <path>` and description.

**Test:** `Alge/aligned/cmd/align.TestHelpDocumentsCheck`

## Document show command

The help output includes the show command with usage `show <path>` and description.

**Test:** `Alge/aligned/cmd/align.TestHelpDocumentsShow`

## Document init command

The help output includes the init command with usage `init <type> <path>` and description.

**Test:** `Alge/aligned/cmd/align.TestHelpDocumentsInit`

## Document list-tests command

The help output includes the list-tests command and description.

**Test:** `Alge/aligned/cmd/align.TestHelpDocumentsListTests`

## Document checkconf command

The help output includes the checkconf command and description.

**Test:** `Alge/aligned/cmd/align.TestHelpDocumentsCheckconf`

## Document version command

The help output includes the version command and description.

**Test:** `Alge/aligned/cmd/align.TestHelpDocumentsVersion`

## Document help command

The help output includes the help command itself and description.

**Test:** `Alge/aligned/cmd/align.TestHelpDocumentsHelp`

## Exit with success code

The help command exits with code 0.

**Test:** `Alge/aligned/cmd/align.TestHelpExitCode`
