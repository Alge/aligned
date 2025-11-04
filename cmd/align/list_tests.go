package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Alge/aligned/internal/config"
	"github.com/Alge/aligned/internal/connectors"
)

func listTests(args []string, stdout, stderr io.Writer) int {
	// Load configuration
	configPath := filepath.Join(".", ".align.yml")
	cfg, err := config.LoadConfiguration(configPath)
	
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(stderr, "Error: .align.yml not found")
			return 1
		}
		fmt.Fprintf(stderr, "Error: Invalid configuration: %v\n", err)
		return 1
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(stderr, "Error: Invalid configuration: %v\n", err)
		return 1
	}
	
	// Discover tests from each connector
	for _, connectorCfg := range cfg.Connectors {
		var connector connectors.Connector

		switch connectorCfg.Type {
		case "go":
			executable := connectorCfg.Executable
			if executable == "" {
				executable = "go"
			}
			connector = connectors.NewGoConnector(executable)
		case "pytest":
			executable := connectorCfg.Executable
			if executable == "" {
				executable = "pytest"
			}
			connector = connectors.NewPytestConnector(executable)
		case "elixir":
			executable := connectorCfg.Executable
			if executable == "" {
				executable = "mix"
			}
			connector = connectors.NewElixirConnector(executable)
		case "gleam":
			executable := connectorCfg.Executable
			if executable == "" {
				executable = "gleam"
			}
			connector = connectors.NewGleamConnector(executable)
		case "vitest":
			executable := connectorCfg.Executable
			if executable == "" {
				executable = "vitest"
			}
			connector = connectors.NewVitestConnector(executable)
		default:
			fmt.Fprintf(stderr, "Error: Unsupported connector type: %s\n", connectorCfg.Type)
			return 1
		}
		
		// Discover tests
		tests, err := connector.DiscoverTests(connectorCfg.Path)
		if err != nil {
			fmt.Fprintf(stderr, "Error discovering tests: %v\n", err)
			return 1
		}
		
		// Print tests
		for _, test := range tests {
			fmt.Fprintln(stdout, test)
		}
	}
	
	return 0
}