package config

import "fmt"

type WebserverConfig struct {
	Endpoint string `yaml:"endpoint"`
}

func (c *WebserverConfig) validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("Wrong ServerConfig.Endpoint")
	}
	return nil
}
