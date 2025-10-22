package main

import (
	"io"
	"os"
)

const version = "0.1.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return printUsage(stderr)
	}
	
	command := args[0]
	
	switch command {
	case "help":
		return printHelp(stdout)
	case "version":
		return printVersion(stdout)
	case "checkconf":
		return checkconf(args[1:], stdout, stderr)
	case "init":
		return initConfig(args[1:], stdout, stderr)
	case "list-tests":
		return listTests(args[1:], stdout, stderr)
	case "show":
		return show(args[1:], stdout, stderr)
	case "check":
		return check(args[1:], stdout, stderr)
	default:
		return printUsage(stderr)
	}
}

func printUsage(stderr io.Writer) int {
	stderr.Write([]byte("align: command not implemented\n"))
	return 1
}

func printHelp(stdout io.Writer) int {
	help := `Aligned - Validate that specifications are covered by tests

Usage: align <command> [options]

Commands:
  check <path>        Validate specification coverage
  show <path>         Display specification structure
  init <type> <path>  Create .align.yml configuration
  list-tests          List all discovered tests
  checkconf           Verify configuration is valid
  version             Show version information
  help                Show this help message
`
	stdout.Write([]byte(help))
	return 0
}