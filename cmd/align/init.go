package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Alge/aligned/internal/config"
	"github.com/Alge/aligned/internal/connectors"
	"gopkg.in/yaml.v3"
)

// Map of supported connectors to their factory functions
// Format: [language]-[framework]
var connectorFactories = map[string]func() connectors.Connector{
	"go-test":        func() connectors.Connector { return connectors.DefaultGoConnector() },
	"python-pytest":  func() connectors.Connector { return connectors.DefaultPytestConnector() },
	"elixir-exunit":  func() connectors.Connector { return connectors.DefaultElixirConnector() },
}

func displayInitHelp(w io.Writer) {
	fmt.Fprintln(w, "Usage: align init <language-framework> <path>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Supported connectors:")
	fmt.Fprintln(w, "  go-test        - Go with built-in testing")
	fmt.Fprintln(w, "  python-pytest  - Python with pytest")
	fmt.Fprintln(w, "  elixir-exunit  - Elixir with ExUnit")
}

func initConfig(args []string, stdout, stderr io.Writer) int {
	// Handle help command or no arguments
	if len(args) == 0 || (len(args) == 1 && args[0] == "help") {
		displayInitHelp(stdout)
		return 0
	}

	// Check arguments
	if len(args) < 2 {
		displayInitHelp(stderr)
		return 1
	}

	connectorType := args[0]
	path := args[1]

	// Check if connector type is supported
	factory, supported := connectorFactories[connectorType]
	if !supported {
		fmt.Fprintf(stderr, "Error: Unsupported connector type: %s\n", connectorType)
		fmt.Fprintln(stderr, "")
		displayInitHelp(stderr)
		return 1
	}

	// Check if config already exists
	configPath := filepath.Join(".", ".align.yml")
	if _, err := os.Stat(configPath); err == nil {
		fmt.Fprintln(stderr, "Error: .align.yml already exists")
		return 1
	}

	// Create connector and generate config
	connector := factory()
	connectorConfig := connector.GenerateConfig(path)

	// Create config using connector's GenerateConfig method
	cfg := config.Configuration{
		Connectors: []config.ConnectorConfig{
			connectorConfig,
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "Error creating config: %v\n", err)
		return 1
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Fprintf(stderr, "Error writing config file: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Created .align.yml with %s connector at %s\n", connectorType, path)
	return 0
}