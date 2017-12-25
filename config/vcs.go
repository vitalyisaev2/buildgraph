package config

type VCSConfig struct {
	Gitlab *GitlabConfig
}

// GitlabConfig describes configuration of Github CI server,
// that will post notifications about CI events
type GitlabConfig struct {
	Endpoint string
	User     string
	Password string
}
