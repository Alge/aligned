package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"version"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode)
	assert.NotEmpty(t, stdout.String(), "version output should not be empty")
}

func TestVersionShowsApplicationVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"version"}, &stdout, &stderr)

	// Verify that the output contains the version constant
	assert.Contains(t, stdout.String(), version)
}

func TestVersionExitCode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"version"}, &stdout, &stderr)

	assert.Equal(t, 0, exitCode, "version command should exit with code 0")
}
