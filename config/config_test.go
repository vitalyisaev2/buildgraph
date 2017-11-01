package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c, err := NewConfig("./config.yml")
	assert.NoError(t, err)
	assert.Equal(t, "buildgraph", c.Storage.Postgres.User)
}
