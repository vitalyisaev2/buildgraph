package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/config"
)

var (
	_ (http.Handler)   = (*handler)(nil)
	_ (common.Service) = (*handler)(nil)
)

type handler struct {
	mux *http.ServeMux
}

func (h *handler) ServeHTTP(http.ResponseWriter, *http.Request) { h.mux.ServeHTTP(w, r) }

type server struct {
	logger     common.Logger
	handler    *handler
	httpServer *http.Server

	errorChan chan<- error

	cfg *config.ServerConfig
}

func (s *server) Start() {
	s.logger.Debug("Starting HTTP server")
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.errorChan <- err
		}
	}()
}

const terminationTimeout = time.Second

func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), terminationTimeout)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}

func NewWebServer(logger common.Logger, cfg *config.ServerConfig, errChan chan<- error) common.Service {
	s := &server{
		httpServer: &http.Server{
			Addr:    s.cfg.Endpoint,
			Handler: s.handler,
		},
		handler: &handler{
			mux: http.NewServeMux(),
		},
		logger:  logger,
		errChan: errChan,
	}
}
