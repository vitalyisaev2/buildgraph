package config

import "fmt"

type ServerConfig struct {
	Endpoint string `yaml:"endpoint"`
}

func (c *ServerConfig) validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("Wrong ServerConfig.Endpoint")
	}
	return nil
}
