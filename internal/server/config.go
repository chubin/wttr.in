package server

import "fmt"

// Config holds the server configuration loaded from YAML.
type Config struct {
	// HTTP/HTTPS ports
	PortHTTP  int `yaml:"portHttp"`  // e.g. 80  - set to 0 to disable
	PortHTTPS int `yaml:"portHttps"` // e.g. 443 - set to 0 to disable

	// === Legacy single certificate (backward compatible) ===
	// Use these only for simple setups with one domain.
	TLSCertFile string `yaml:"tlsCertFile,omitempty"`
	TLSKeyFile  string `yaml:"tlsKeyFile,omitempty"`

	// === Modern multi-certificate configuration ===
	// Preferred way for multiple domains / wildcards.
	TLSCerts []TLSCertConfig `yaml:"tlsCerts,omitempty"`
}

// TLSCertConfig defines one certificate + its matching domain(s).
type TLSCertConfig struct {
	// Domain this certificate should be used for.
	// Examples:
	//   - "wttr.in"           (exact match)
	//   - "*.wttr.in"         (wildcard)
	Domain string `yaml:"domain"`

	// Paths to certificate and private key (PEM format)
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

// Validate performs basic validation of the configuration.
func (c *Config) Validate() error {
	if c.PortHTTP == 0 && c.PortHTTPS == 0 {
		return ErrNoServersConfigured
	}

	// Check legacy single cert
	if c.TLSCertFile != "" || c.TLSKeyFile != "" {
		if c.TLSCertFile == "" {
			return fmt.Errorf("tlsCertFile is required when tlsKeyFile is set")
		}
		if c.TLSKeyFile == "" {
			return fmt.Errorf("tlsKeyFile is required when tlsCertFile is set")
		}
	}

	// Validate multiple certs
	for i, tc := range c.TLSCerts {
		if tc.Domain == "" {
			return fmt.Errorf("tlsCerts[%d]: domain is required", i)
		}
		if tc.CertFile == "" {
			return fmt.Errorf("tlsCerts[%d]: certFile is required", i)
		}
		if tc.KeyFile == "" {
			return fmt.Errorf("tlsCerts[%d]: keyFile is required", i)
		}
	}

	return nil
}

// HasMultipleCerts returns true if the new multi-certificate config is used.
func (c *Config) HasMultipleCerts() bool {
	return len(c.TLSCerts) > 0
}

// HasLegacyCert returns true if the old single certificate fields are populated.
func (c *Config) HasLegacyCert() bool {
	return c.TLSCertFile != "" && c.TLSKeyFile != ""
}

// Example configurations

/*
# ================================================
# Example 1: Simple single domain (legacy style)
# ================================================
portHttp: 80
portHttps: 443
tlsCertFile: "/etc/letsencrypt/live/wttr.in/fullchain.pem"
tlsKeyFile:  "/etc/letsencrypt/live/wttr.in/privkey.pem"


# ================================================
# Example 2: Multiple domains + wildcard (recommended)
# ================================================
portHttp: 80
portHttps: 443

tlsCerts:
  - domain: "wttr.in"
    certFile: "/wttr.in/certs/wttr.in/fullchain.pem"
    keyFile:  "/wttr.in/certs/wttr.in/privkey.pem"

  - domain: "*.wttr.in"
    certFile: "/wttr.in/certs/wildcard.wttr.in/fullchain.pem"
    keyFile:  "/wttr.in/certs/wildcard.wttr.in/privkey.pem"
*/
