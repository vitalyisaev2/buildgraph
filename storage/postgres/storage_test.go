package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
	"go.uber.org/zap"
)

var (
	stubLogger, _ = zap.NewDevelopment()
)

// integration tests for PostgreSQL-backed storage
type storageSuite struct {
	suite.Suite
	logger  *zap.Logger
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

func (s *storageSuite) Test_Author() {
	inputAuthor := &author{
		name:  "Vitaly Isaev",
		email: "admin@vitalya.ru",
	}
	err := s.storage.SaveAuthor(s.ctx, inputAuthor)
	s.Assert().NoError(err)
	s.Assert().NotZero(inputAuthor.GetID())

	outputAuthor, err := s.storage.GetAuthor(s.ctx, inputAuthor.GetName(), inputAuthor.GetEmail())
	s.Assert().NoError(err)
	if s.Assert().NotNil(outputAuthor) {
		s.Assert().Equal(inputAuthor.GetName(), outputAuthor.GetName())
		s.Assert().Equal(inputAuthor.GetEmail(), outputAuthor.GetEmail())
		s.Assert().Equal(inputAuthor.GetID(), outputAuthor.GetID())
	}
}

func (s *storageSuite) TearDownSuite() { s.logger.Sync() }

func TestIntegration_PostgreSQLStorage(t *testing.T) {
	suite.Run(t, &storageSuite{})
}
