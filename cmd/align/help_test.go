package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"help"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode)
	assert.NotEmpty(t, stdout.String(), "help output should not be empty")
}

func TestHelpIncludesDescription(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	assert.Contains(t, stdout.String(), "Validate that specifications are covered by tests")
}

func TestHelpShowsUsagePattern(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	assert.Contains(t, stdout.String(), "align <command> [options]")
}

func TestHelpDocumentsCheck(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	output := stdout.String()
	assert.Contains(t, output, "check")
	assert.Contains(t, output, "<path>")
}

func TestHelpDocumentsShow(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	output := stdout.String()
	assert.Contains(t, output, "show")
	assert.Contains(t, output, "<path>")
}

func TestHelpDocumentsInit(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	output := stdout.String()
	assert.Contains(t, output, "init")
	assert.Contains(t, output, "<type>")
	assert.Contains(t, output, "<path>")
}

func TestHelpDocumentsListTests(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	assert.Contains(t, stdout.String(), "list-tests")
}

func TestHelpDocumentsCheckconf(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	assert.Contains(t, stdout.String(), "checkconf")
}

func TestHelpDocumentsVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	assert.Contains(t, stdout.String(), "version")
}

func TestHelpDocumentsHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"help"}, &stdout, &stderr)

	assert.Contains(t, stdout.String(), "help")
}

func TestHelpExitCode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"help"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "help command should exit with code 0")
}
