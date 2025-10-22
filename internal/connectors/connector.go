package connectors

import "github.com/Alge/aligned/internal/config"

type Connector interface {
	DetectFramework() (bool, error)
	GenerateConfig(path string) config.ConnectorConfig
	DiscoverTests(path string) ([]string, error)
}
