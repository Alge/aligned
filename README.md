# Aligned

A tool for validating that specifications are covered by tests

Aligned does not automatically solve all your problems, but it:

 * Tries to help you write actionable, testable technical specifications
 * Tries to make the gap between non-technical product owners and development teams smaller by providing a single source of truth about what the system should do that can be read and understood by both sides.

**Note:** This project was developed using AI-assisted development, where I served as product owner and architect—writing specifications, designing the system architecture, and validating implementations.

## Installation

```bash
go install github.com/Alge/aligned/cmd/align@latest
```

Or build from source:
```bash
git clone https://github.com/Alge/aligned.git
cd aligned
make build
# Binary will be in bin/align
```

## Specification files

Aligned uses markdown files for its specifications, and uses headings to organize the files.

Example:

```
# Hello world application

## Application prints hello world
When run the application should print the string 'Hello World!'

**Test:** `test_prints_hello_world`

## Application return status 0 on success

**Test:** `application_returns_status_0`
```

You can then use the `$ align show` command to show a summary of the specification.
```
$ align show examples/hello_world_spec.md 
Hello world application
· Application prints hello world
·   Test: test_prints_hello_world
· Application return status 0 on success
·   Test: application_returns_status_0
```

## Commands

### init

`$ align init <connector-type> <path>`

Creates the `.align.yml` file pre-configured with default configuration for the connector type and a path to your test suite.

### show

`$ align show <path>`
prints a summary representation of the specification

### check

`$ align check <path> [-v]`

This is the main command for aligned. It parses a specification file, and:
* Makes sure all leaf nodes has a reference to a test that exists
* Makes sure all sections implementing interfaces includes all required sections
* Prints a summary of the current state of the specification. By default passing sections are collapsed and only shows the root level. The whole tree can be shown using the `-v` (verbose) flag.


## Supported test frameworks

* **Go** - Uses `go test` for discovery
* **Pytest** - Python testing via `pytest --collect-only`
* **Elixir** - ExUnit via `mix test --trace`

### Adding new frameworks

Feel free to send a PR to add support for more test frameworks!

Trying to live what we preach, new code should be added in these steps:

1. Create the specification for the new connector. It should implement the `Test Framework Connector Interface` interface.

2. Write the tests and implementation for the feature

3. Add the test references to the specification as they are written

4. Verify that all tests passes using `make test`

5. Verify that the specification passes using `make align`

Once all tests in the interface exists and passes, you are done!

## Development

Run tests:
```bash
make test              # Run all tests
go test ./... -v       # Verbose output (shows skipped tests)
```

Note: Full test suite requires Go, pytest, and mix installed. Tests for unavailable frameworks will be skipped automatically.