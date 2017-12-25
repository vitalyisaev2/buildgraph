package gitlab

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

const (
	pathAPI     string = "/api/v3"
	pathSession string = pathAPI + "/session"
)

// gitlabConfig provides parameters to access Gitlab API
type gitlabConfig struct {
	endpoint string
	login    string
	password string
	timeout  time.Duration
	projects map[string]string
}

// this interface contains a limited set of methods used only in integration tests
type gitlabAPI interface {
	CreateGroup(name string) error
	GetGroup(name string) error
}

// gitlabAPIImpl implements methods required in tests
type gitlabAPIImpl struct {
	client *gitlab.Client
	token  string
	logger *logrus.Logger
	cfg    *gitlabConfig
}

// Create
func (api *gitlabAPIImpl) CreateGroup(name string) error {
	opt := &gitlab.CreateGroupOptions{
		Name: &name,
		Path: &name,
	}
	_, _, err := api.client.Groups.CreateGroup(opt)
	if err != nil {
		return err
	}
	return nil
}

func (api *gitlabAPIImpl) GetGroup(name string) error {
	_, _, err := api.client.Groups.GetGroup(name)
	if err != nil {
		return err
	}
	return nil
}

func newGitlabAPI(cfg *gitlabConfig) (gitlabAPI, error) {

	// prepare http client
	httpClient := &http.Client{
		Timeout: cfg.timeout,
	}

	// prepare logger
	logger := logrus.New()

	// obtain API token
	tokenName := "test" + time.Now().Format("2006-01-02_15.04.02")
	tokenExpirationDate := time.Now().Add(48 * time.Hour)
	token, err := NewPersonalAccessToken(
		cfg.endpoint,
		cfg.login,
		cfg.password,
		tokenName,
		tokenExpirationDate,
	)
	if err != nil {
		return nil, err
	}
	// prepare Gitlab API client
	gitlabClient := gitlab.NewClient(httpClient, token)
	gitlabClient.SetBaseURL(cfg.endpoint + pathAPI)

	api := &gitlabAPIImpl{
		cfg:    cfg,
		client: gitlabClient,
		logger: logger,
	}
	return api, nil
}
