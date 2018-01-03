package postgres

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
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
	cfg, err := config.NewConfig("../../config/example.yml")
	if err != nil {
		s.logger.WithError(err).Error("failed to read config")
		s.T().FailNow()
	}

	// Database initialization
	s.storage, err = NewStorage(stubLogger, cfg.Storage.Postgres)
	if err != nil {
		s.logger.WithError(err).Error("failed to initialize storage")
		s.T().FailNow()
	}
}

func (s *storageSuite) TearDownSuite() {}

func TestIntegration_PostgreSQLStorage(t *testing.T) {
	suite.Run(t, &storageSuite{})
}
