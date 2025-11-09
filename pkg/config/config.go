package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	DefaultPort    = 9222
	DefaultTimeout = 30 * time.Second
	MinPort        = 1
	MaxPort        = 65535
)

// Config holds configuration for browser connection
type Config struct {
	// Port is the Chrome DevTools Protocol debugging port
	Port int

	// Timeout for browser operations (0 means no timeout)
	Timeout time.Duration
}

// Default returns a Config with default values
func Default() *Config {
	return &Config{
		Port:    ResolvePort(0),
		Timeout: DefaultTimeout,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Port < MinPort || c.Port > MaxPort {
		return fmt.Errorf("port must be between %d and %d, got %d", MinPort, MaxPort, c.Port)
	}
	if c.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}
	return nil
}

// ResolvePort determines which port to use based on flag, environment, or default
// Priority: flagPort (if > 0) > BROW_DEBUG_PORT env var > DefaultPort
func ResolvePort(flagPort int) int {
	if flagPort > 0 {
		return flagPort
	}

	if envPort := os.Getenv("BROW_DEBUG_PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil && port > 0 {
			return port
		}
	}

	return DefaultPort
}
