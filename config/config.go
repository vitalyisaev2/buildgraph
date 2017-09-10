package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config top-level structure
type Config struct {
	Storage *StorageConfig
}

func (c *Config) validate() error {
	var err error
	if c.Storage != nil {
		err = c.Storage.validate()
	}
	return err
}

// StorageConfig describes various configurations of a storage layer
type StorageConfig struct {
	Postgres *PostgresConfig
}

func (c *StorageConfig) validate() error {
	var err error
	if c.Postgres != nil {
		err = c.Postgres.validate()
	}
	return err
}

// PostgresConfig describes configuration of PostgreSQL client
type PostgresConfig struct {
	Endpoint string
	User     string
	Password string
	Database string
}

func (pc *PostgresConfig) validate() error {
	if pc.Endpoint == "" {
		return fmt.Errorf("Wrong Storage.PostgresConfig.Endpoint value: %s", pc.Endpoint)
	}
	if pc.User == "" {
		return fmt.Errorf("Wrong Storage.PostgresConfig.User value: %s", pc.User)
	}
	if pc.Password == "" {
		return fmt.Errorf("Wrong Storage.PostgresConfig.Password value: %s", pc.Password)
	}
	if pc.Database == "" {
		return fmt.Errorf("Wrong Storage.PostgresConfig.Database value: %s", pc.Password)
	}
	return nil
}

func (pc *PostgresConfig) URL() string {
	// FIXME: TLS
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pc.User, pc.Password, pc.Endpoint, pc.Database)
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
