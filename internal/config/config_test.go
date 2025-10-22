package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	t.Run("loads valid configuration with single connector", func(t *testing.T) {
		tempDir := t.TempDir()

		configContent := `connectors:
  - type: go
    executable: go
    path: ./
`
		configPath := filepath.Join(tempDir, ".align.yml")
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		assert.NoError(t, err)

		config, err := LoadConfiguration(configPath)

		assert.NoError(t, err)
		assert.Len(t, config.Connectors, 1)
		assert.Equal(t, "go", config.Connectors[0].Type)
		assert.Equal(t, "go", config.Connectors[0].Executable)
		assert.Equal(t, "./", config.Connectors[0].Path)
	})

	t.Run("loads configuration with multiple connectors", func(t *testing.T) {
		tempDir := t.TempDir()

		configContent := `connectors:
  - type: go
    executable: go
    path: ./
  - type: pytest
    executable: pytest
    path: ./tests
`
		configPath := filepath.Join(tempDir, ".align.yml")
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		assert.NoError(t, err)

		config, err := LoadConfiguration(configPath)

		assert.NoError(t, err)
		assert.Len(t, config.Connectors, 2)
		assert.Equal(t, "go", config.Connectors[0].Type)
		assert.Equal(t, "pytest", config.Connectors[1].Type)
	})

	t.Run("returns error for nonexistent file", func(t *testing.T) {
		_, err := LoadConfiguration("/nonexistent/path/.align.yml")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file")
	})

	t.Run("returns error for invalid YAML", func(t *testing.T) {
		tempDir := t.TempDir()

		invalidYAML := `connectors:
  - type: go
    executable: go
    path: ./
    invalid indentation
`
		configPath := filepath.Join(tempDir, ".align.yml")
		err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
		assert.NoError(t, err)

		_, err = LoadConfiguration(configPath)

		assert.Error(t, err)
		// YAML unmarshal errors typically mention "yaml" or parsing
		errMsg := err.Error()
		assert.True(t,
			len(errMsg) > 0,
			"error should provide details about YAML parsing failure")
	})

	t.Run("validates empty connectors", func(t *testing.T) {
		tempDir := t.TempDir()

		configContent := `connectors: []`
		configPath := filepath.Join(tempDir, ".align.yml")
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		assert.NoError(t, err)

		config, err := LoadConfiguration(configPath)
		assert.NoError(t, err, "loading should succeed")

		// But validation should fail
		err = config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no connectors")
	})

	t.Run("validates missing required fields", func(t *testing.T) {
		tempDir := t.TempDir()

		testCases := []struct {
			name      string
			content   string
			errorText string
		}{
			{
				name: "missing type",
				content: `connectors:
  - path: ./
`,
				errorText: "type is required",
			},
			{
				name: "missing path",
				content: `connectors:
  - type: go
    executable: go
`,
				errorText: "path is required",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				configPath := filepath.Join(tempDir, tc.name+".yml")
				err := os.WriteFile(configPath, []byte(tc.content), 0644)
				assert.NoError(t, err)

				config, err := LoadConfiguration(configPath)
				assert.NoError(t, err, "loading should succeed")

				// But validation should fail
				err = config.Validate()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorText)
			})
		}
	})
}