package webserver

import (
	"context"
	"net/http"
	"time"

	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/config"
	"go.uber.org/zap"
)

var (
	_ (http.Handler)   = (*handler)(nil)
	_ (common.Service) = (*server)(nil)
)

type handler struct {
	mux *http.ServeMux
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.mux.ServeHTTP(w, r) }

type server struct {
	logger     *zap.Logger
	handler    *handler
	httpServer *http.Server

	errChan chan<- error

	cfg *config.WebserverConfig
}

const terminationTimeout = time.Second

func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), terminationTimeout)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}

func NewWebServer(logger *zap.Logger, cfg *config.WebserverConfig, errChan chan<- error) common.Service {
	h := &handler{
		mux: http.NewServeMux(),
	}

	s := &server{
		httpServer: &http.Server{
			Addr:    cfg.Endpoint,
			Handler: h,
		},
		logger:  logger,
		errChan: errChan,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.errChan <- err
		}
	}()

	return s
}
