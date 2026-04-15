// Package config handles loading and validating portwatch configuration.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch daemon configuration.
type Config struct {
	Scan    ScanConfig  `yaml:"scan"`
	Alert   AlertConfig `yaml:"alert"`
}

// ScanConfig controls port scanning behaviour.
type ScanConfig struct {
	PortStart int           `yaml:"port_start"`
	PortEnd   int           `yaml:"port_end"`
	Interval  time.Duration `yaml:"interval"`
}

// AlertConfig controls how alerts are delivered.
type AlertConfig struct {
	Backend string `yaml:"backend"` // "log", "stdout"
	LogFile string `yaml:"log_file"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Scan: ScanConfig{
			PortStart: 1,
			PortEnd:   65535,
			Interval:  30 * time.Second,
		},
		Alert: AlertConfig{
			Backend: "log",
		},
	}
}

// Load reads a YAML config file from path and merges it over defaults.
func Load(path string) (Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Validate returns an error if the configuration is logically invalid.
func (c Config) Validate() error {
	if c.Scan.PortStart < 1 || c.Scan.PortStart > 65535 {
		return fmt.Errorf("config: scan.port_start must be between 1 and 65535")
	}
	if c.Scan.PortEnd < 1 || c.Scan.PortEnd > 65535 {
		return fmt.Errorf("config: scan.port_end must be between 1 and 65535")
	}
	if c.Scan.PortStart > c.Scan.PortEnd {
		return fmt.Errorf("config: scan.port_start must be <= scan.port_end")
	}
	if c.Scan.Interval <= 0 {
		return fmt.Errorf("config: scan.interval must be positive")
	}
	return nil
}
