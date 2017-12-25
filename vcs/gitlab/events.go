package gitlab

import (
	"fmt"
	"time"

	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/vcs"
)

var _ vcs.PushEvent = (*PushEvent)(nil)

type PushEvent struct {
	common.Object
	ObjectKind        string      `json:"object_kind"`
	Before            string      `json:"before"`
	After             string      `json:"after"`
	Ref               string      `json:"ref"`
	CheckoutSHA       string      `json:"checkout_sha"`
	UserName          string      `json:"user_name"`
	UserUsername      string      `json:"user_username"`
	UserEmail         string      `json:"user_email"`
	UserAvatar        string      `json:"user_avatar"`
	ProjectID         int         `json:"project_id"`
	Project           *Project    `json:"project"`
	Repository        *Repository `json:"repository"`
	Commits           []*Commit   `json:"commits"`
	TotalCommitsCount int         `json:"total_commits_count"`
}

func (p *PushEvent) GetProject() vcs.Project { return p.Project }

func (p *PushEvent) GetCommits() []vcs.Commit {
	results := make([]vcs.Commit, 0, len(p.Commits))
	for _, c := range p.Commits {
		results = append(results, c)
	}
	return results
}

// Project

var _ vcs.Project = (*Project)(nil)

type Project struct {
	common.Object
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

func (p *Project) GetName() string { return p.Name }

func (p *Project) GetNamespace() string { return p.Namespace }

func (p *Project) GetHTTPURL() string { return p.HTTPURL }

// Repository

type Repository struct {
	common.Object
	Name            string `json:"name"`
	URL             string `json:"url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	GitHTTPURL      string `json:"git_http_url"`
	GitSSHURL       string `json:"git_ssh_url"`
	VisibilityLevel int    `json:"visibility_level"`
}

// Commit

var _ vcs.Commit = (*Commit)(nil)

type Commit struct {
	common.Object
	Hash      string   `json:"id"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
	Author    *Author  `json:"author"`
	Added     []string `json:"added"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
}

func (c *Commit) GetHash() string { return c.Hash }

func (c *Commit) GetMessage() string { return c.Message }

func (c *Commit) GetTimestamp() time.Time {
	t, err := time.Parse(time.RFC3339, c.Timestamp)
	if err != nil {
		msg := fmt.Sprintf(
			"implementation error: timestamp layout '%s' doesn't match with valut '%s'",
			time.RFC3339, c.Timestamp,
		)
		panic(msg)
	}
	return t
}

func (c *Commit) GetAuthor() vcs.Author { return c.Author }

func (c *Commit) GetURL() string { return c.URL }

func (c *Commit) GetAdded() []string { return c.Added }

func (c *Commit) GetModified() []string { return c.Modified }

func (c *Commit) GetRemoved() []string { return c.Removed }

// Author

var _ vcs.Author = (*Author)(nil)

type Author struct {
	common.Object
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (a *Author) GetName() string  { return a.Name }
func (a *Author) GetEmail() string { return a.Email }
