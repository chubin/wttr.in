package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/chubin/wttr.in/internal/cache"
	"github.com/chubin/wttr.in/internal/ip"
	"github.com/chubin/wttr.in/internal/location"
	"github.com/chubin/wttr.in/internal/logging"
	"github.com/chubin/wttr.in/internal/server"
	"github.com/chubin/wttr.in/internal/uplink"
	"github.com/chubin/wttr.in/internal/weather"
)

type Config struct {
	Geo     *location.Config
	IP      *ip.Config
	Weather struct {
		WWO *weather.WWOConfig
	}
	Cache   cache.Config
	Logging logging.Config
	Uplink  uplink.Config
	Server  server.Config
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

	// Unmarshal YAML data into the Config struct with strict checking for unknown fields
	err = yaml.UnmarshalStrict(data, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	return config, nil
}

func (c *Config) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
