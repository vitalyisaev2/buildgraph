package postgres

import (
	"time"

	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/vcs"
)

// PushEvent

var _ (vcs.PushEvent) = (*pushEvent)(nil)

type pushEvent struct {
	common.Object
	project *project
	commits []*commit
}

func (e *pushEvent) GetProject() vcs.Project { return e.project }
func (e *pushEvent) GetCommits() []vcs.Commit {
	result := make([]vcs.Commit, 0, len(e.commits))
	for _, c := range e.commits {
		result = append(result, c)
	}
	return result
}

// Project

var _ (vcs.Project) = (*project)(nil)

type project struct {
	common.Object
	namespace string
	name      string
	httpURL   string
}

func (m *project) GetNamespace() string { return m.namespace }
func (m *project) GetName() string      { return m.name }
func (m *project) GetHTTPURL() string   { return m.httpURL }

// Commit

var _ (vcs.Commit) = (*commit)(nil)

type commit struct {
	common.Object
	hash      string
	message   string
	timestamp time.Time
	author    *author
	url       string
	added     []string
	modified  []string
	removed   []string
}

func (c *commit) GetHash() string         { return c.hash }
func (c *commit) GetMessage() string      { return c.message }
func (c *commit) GetTimestamp() time.Time { return c.timestamp }
func (c *commit) GetAuthor() vcs.Author   { return c.author }
func (c *commit) GetURL() string          { return c.url }
func (c *commit) GetAdded() []string      { return c.added }
func (c *commit) GetModified() []string   { return c.modified }
func (c *commit) GetRemoved() []string    { return c.removed }

// Author

var _ (vcs.Author) = (*author)(nil)

type author struct {
	common.Object
	name  string
	email string
}

func (m *author) GetName() string  { return m.name }
func (m *author) GetEmail() string { return m.email }
