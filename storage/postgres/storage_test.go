package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
	"go.uber.org/zap"
)

const (
	dockerTimeout    = time.Minute
	postgresImage    = "postgres" // official PostgreSQL Docker image
	postgresUser     = "buildgraph"
	postgresPassword = "password"
	postgresDatabase = "buildgraph"
)

// integration tests for PostgreSQL-backed storage
type storageSuite struct {
	suite.Suite

	logger *zap.Logger

	dockerClient *client.Client
	containerID  string // during test suite, a container with PostgreSQL will be created

	storage storage.Storage

	ctx context.Context
}

func (s *storageSuite) SetupSuite() {
	var (
		err error
	)

	// root context
	s.ctx = context.Background()

	ctx, cancel := context.WithTimeout(s.ctx, dockerTimeout)
	defer cancel()

	// Create logger
	if s.logger, err = zap.NewDevelopment(); err != nil {
		s.T().Fatalf("Failed to initialize logger: %v", err)
	}

	// Create Docker client
	if s.dockerClient, err = client.NewEnvClient(); err != nil {
		s.logger.Error("Failed to init Docker client", zap.Error(err))
		s.T().FailNow()
	}

	// Pull image
	s.logger.Debug("Pulling PostgreSQL image")
	_, err = s.dockerClient.ImagePull(ctx, postgresImage, types.ImagePullOptions{})
	if err != nil {
		s.logger.Error("Failed to pull Docker image", zap.Error(err))
		s.T().FailNow()
	}

	var (
		containerConfig = &container.Config{
			Image: postgresImage,
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", postgresUser),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", postgresPassword),
				fmt.Sprintf("POSTGRES_DB=%s", postgresDatabase),
			},
		}
		hostConfig = &container.HostConfig{
			PortBindings: nat.PortMap{
				"5432/tcp": []nat.PortBinding{
					nat.PortBinding{
						HostIP:   "0.0.0.0",
						HostPort: "5432",
					},
				},
			},
		}
		resp container.ContainerCreateCreatedBody
	)

	// Run container with PostgreSQL
	s.logger.Debug("Creating PostgreSQL container")
	resp, err = s.dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, "postgres-test")
	if err != nil {
		s.logger.Error("Failed to create container", zap.Error(err))
		s.T().FailNow()
	}
	s.containerID = resp.ID

	s.logger.Debug("Starting PostgreSQL container")
	err = s.dockerClient.ContainerStart(ctx, s.containerID, types.ContainerStartOptions{})
	if err != nil {
		s.logger.Error("Failed to start container", zap.Error(err))
		s.T().FailNow()
	}

	// Run storage abstraction (this will cause migrations as well)
	storageConfig := &config.PostgresConfig{
		Endpoint: "localhost:5432",
		User:     postgresUser,
		Password: postgresPassword,
		Database: postgresDatabase,
	}

	// Database initialization
	time.Sleep(5 * time.Second)

	s.storage, err = NewStorage(storageConfig)
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

func (s *storageSuite) TearDownSuite() {
	var (
		err         error
		ctx, cancel = context.WithTimeout(context.Background(), dockerTimeout)
	)
	defer cancel()

	s.logger.Debug("Removing PostgreSQL container")

	err = s.dockerClient.ContainerRemove(
		ctx,
		s.containerID,
		types.ContainerRemoveOptions{Force: true},
	)
	if err != nil {
		s.T().Fatalf("Failed to remove container: %v", err)
	}

	s.dockerClient.Close()
	s.logger.Sync()
}

func TestPostgreSQLStorage(t *testing.T) {
	suite.Run(t, &storageSuite{})
}
