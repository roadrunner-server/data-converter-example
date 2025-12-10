package dataconverter

import "errors"

type Config struct {
	// yaml, sample mapping to your `key` in the .rr.yaml
	key string `mapstructure:"key"`
}

// InitDefaults used to initialize default configuration values
func (c *Config) InitDefaults() error {
	if c.key == "" {
		return errors.New("key is required")
	}

	// other default values initialization can go here

	return nil
}
