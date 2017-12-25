package service

import (
	"github.com/sirupsen/logrus"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
	"github.com/vitalyisaev2/buildgraph/storage/postgres"
)

type Collection struct {
	Logger  *logrus.Logger // TODO: turn into interface
	Storage storage.Storage
}

func (c *Collection) Stop() {
	c.Logger.Debug("stopping storage")
	c.Storage.Stop()
}

func NewCollection(logger *logrus.Logger, cfg *config.Config) (*Collection, error) {
	var (
		c   Collection
		err error
	)

	c.Logger = logger

	c.Logger.Info("starting storage")
	if c.Storage, err = postgres.NewStorage(c.Logger, cfg.Storage.Postgres); err != nil {
		return nil, err
	}

	return &c, nil
}
