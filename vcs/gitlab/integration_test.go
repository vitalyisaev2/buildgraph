package gitlab

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	projectGraphFile = "../../test/data/simple2.yml"
)

type gitlabSuite struct {
	suite.Suite
	api gitlabAPI
}

func (s *gitlabSuite) SetupSuite() {
	config := &gitlabConfig{
		credentials: &credentials{
			login:    "root",
			password: "password",
		},
		timeout: 10 * time.Second,
	}

	config.endpoint = "http://localhost:10080"

	var err error
	s.api, err = newGitlabAPI(config)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *gitlabSuite) TestIntegration_CreateGroup() {
	err := s.api.CreateGroup("buildgraph")
	s.Assert().NoError(err)
	err = s.api.GetGroup("buildgraph")
	s.Assert().NoError(err)
}

func (s *gitlabSuite) TearDownSuite() {}

func TestIntegration_Gitlab(t *testing.T) {
	suite.Run(t, &gitlabSuite{})
}
