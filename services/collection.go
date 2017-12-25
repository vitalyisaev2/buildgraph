package services

import (
	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
	"github.com/vitalyisaev2/buildgraph/storage/postgres"
	"github.com/vitalyisaev2/buildgraph/webserver"
	"go.uber.org/zap"
)

type Collection struct {
	Logger    *zap.Logger // TODO: switch it to interface
	WebServer common.Service
	Storage   storage.Storage
}

func (c *Collection) Stop() {
	c.Logger.Debug("Stopping webserver")
	c.WebServer.Stop()
	c.Logger.Debug("Stopping storage")
	c.Storage.Stop()
}

func NewCollection(logger *zap.Logger, cfg *config.Config, errChan chan<- error) (*Collection, error) {
	var (
		c   Collection
		err error
	)

	c.Logger = logger

	c.Logger.Info("Starting web server")
	c.WebServer = webserver.NewWebServer(c.Logger, cfg.Webserver, errChan)

	c.Logger.Info("Starting storage")
	if c.Storage, err = postgres.NewStorage(c.Logger, cfg.Storage.Postgres); err != nil {
		return nil, err
	}

	return &c, nil
}
