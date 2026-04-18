// Package config handles loading and validation of portwatch configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// PortRange defines the range of ports to scan, e.g. "1-1024".
	PortRange string `json:"port_range"`

	// StateFile is the path to the file used to persist port state.
	StateFile string `json:"state_file"`

	// AlertEmail is an optional email address to send alerts to.
	AlertEmail string `json:"alert_email,omitempty"`

	// Concurrency controls how many ports are scanned in parallel.
	Concurrency int `json:"concurrency"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		PortRange:   "1-1024",
		StateFile:   ".portwatch_state.json",
		Concurrency: 100,
	}
}

// Load reads a JSON config file from the given path and merges it with defaults.
// Fields absent in the file retain their default values.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are acceptable.
func (c *Config) Validate() error {
	if c.PortRange == "" {
		return fmt.Errorf("config: port_range must not be empty")
	}
	if c.StateFile == "" {
		return fmt.Errorf("config: state_file must not be empty")
	}
	if c.Concurrency <= 0 {
		return fmt.Errorf("config: concurrency must be greater than zero")
	}
	return nil
}
