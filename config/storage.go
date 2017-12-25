package config

import "fmt"

// StorageConfig describes various configurations of a storage layer
type StorageConfig struct {
	Postgres *PostgresConfig `yaml:"postgres"`
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
	Endpoint string `yaml:"endpoint"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (pc *PostgresConfig) validate() error {
	if pc.Endpoint == "" {
		return fmt.Errorf("Wrong PostgresConfig.Endpoint value: %s", pc.Endpoint)
	}
	if pc.User == "" {
		return fmt.Errorf("Wrong PostgresConfig.User value: %s", pc.User)
	}
	if pc.Password == "" {
		return fmt.Errorf("Wrong PostgresConfig.Password value: %s", pc.Password)
	}
	if pc.Database == "" {
		return fmt.Errorf("Wrong PostgresConfig.Database value: %s", pc.Password)
	}
	return nil
}

func (pc *PostgresConfig) URL() string {
	// TODO: TLS
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pc.User, pc.Password, pc.Endpoint, pc.Database)
}
