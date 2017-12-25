package webserver

import (
	"context"
	"net/http"
	"time"

	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/service"
)

const (
	// connection termination timeout on server exit
	terminationTimeout = time.Second
)

var _ Webserver = (*server)(nil)

type server struct {
	httpServer *http.Server

	// server configuration
	cfg *config.WebserverConfig

	// long-running server subsystems
	services *service.Collection

	// channel to dump fatal error to
	errChan chan<- error
}

func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), terminationTimeout)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}

func NewWebServer(services *service.Collection, cfg *config.WebserverConfig, errChan chan<- error) common.Service {
	s := &server{
		httpServer: &http.Server{
			Addr: cfg.Endpoint,
		},
		services: services,
		errChan:  errChan,
	}

	// compose multiplexor from gorilla router and negroni middleware
	router := newRouter(s)
	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(services.Logger, "webserver"))
	n.UseHandler(router)
	s.httpServer.Handler = n

	go func() {
		services.Logger.WithField("endpoint", cfg.Endpoint).Debug("starting listener")
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.errChan <- err
		}
	}()

	return s
}
