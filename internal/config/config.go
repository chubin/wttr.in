package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/chubin/wttr.go/internal/ip"
	"github.com/chubin/wttr.go/internal/location"
	"github.com/chubin/wttr.go/internal/weather"
)

type Config struct {
	Geo     *location.Config
	IP      *ip.Config
	Weather struct {
		WWO *weather.WWOConfig
	}
}

// LoadFromYAML loads configuration from a YAML file and returns a pointer to Config
func LoadFromYAML(filePath string) (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	// Create a new Config instance
	config := &Config{}

	// Unmarshal YAML data into the Config struct
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	return config, nil
}

func (c *Config) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
