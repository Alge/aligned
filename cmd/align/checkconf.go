package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	
	"github.com/Alge/aligned/internal/config"
	"gopkg.in/yaml.v3"
)

func checkconf(args []string, stdout, stderr io.Writer) int {
	// Check for -v flag
	verbose := false
	if len(args) > 0 && args[0] == "-v" {
		verbose = true
	}
	
	// Try to load config
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
	
	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(stderr, "Error: Invalid configuration: %v\n", err)
		return 1
	}
	
	// If verbose, print config
	if verbose {
		data, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Fprintf(stderr, "Error formatting config: %v\n", err)
			return 1
		}
		fmt.Fprint(stdout, string(data))
	} else {
		fmt.Fprintln(stdout, "Configuration is valid")
	}
	
	return 0
}