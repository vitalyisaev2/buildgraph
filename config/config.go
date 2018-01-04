package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config top-level structure
type Config struct {
	Storage   *StorageConfig   `yaml:"storage"`
	Webserver *WebserverConfig `yaml:"webserver"`
	Projects  *ProjectsConfig  `yaml:"projects"`
}

func (c *Config) validate() error {

	if c.Storage == nil {
		return fmt.Errorf("Missing required section 'storage'")
	}
	if err := c.Storage.validate(); err != nil {
		return err
	}

	if c.Webserver == nil {
		return fmt.Errorf("Missing required section 'server'")
	}
	if err := c.Webserver.validate(); err != nil {
		return err
	}

	if c.Projects == nil {
		return fmt.Errorf("Missing required section 'projects'")
	}
	if err := c.Projects.validate(); err != nil {
		return err
	}

	return nil
}

func NewConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
