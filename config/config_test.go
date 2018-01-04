package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c, err := NewConfig("./example.yml")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "buildgraph", c.Storage.Postgres.User)
}
