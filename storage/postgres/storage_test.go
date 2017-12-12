package postgres

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
	"go.uber.org/zap"
)

var (
	stubLogger = logrus.New()
)

// integration tests for PostgreSQL-backed storage
type storageSuite struct {
	suite.Suite
	logger  *logrus.Logger
	storage storage.Storage
	ctx     context.Context
}

func (s *storageSuite) SetupSuite() {
	var err error

	s.ctx = context.Background()
	s.logger = stubLogger

	// Run storage abstraction (this will cause migrations as well)
	cfg, _ := config.NewConfig("../../config/example.yaml")

	// Database initialization
	s.storage, err = NewStorage(stubLogger, cfg.Storage.Postgres)
	if err != nil {
		s.logger.Error("Failed to initialize storage", zap.Error(err))
		s.T().Fail()
	}

}

func (s *storageSuite) TearDownSuite() {}

func TestIntegration_PostgreSQLStorage(t *testing.T) {
	suite.Run(t, &storageSuite{})
}
