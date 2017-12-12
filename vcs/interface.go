package vcs

import (
	"time"

	"github.com/vitalyisaev2/buildgraph/common"
)

type PushEvent interface {
	common.Model
	GetProject() Project
	GetCommits() []Commit
}

type Project interface {
	common.Model
	GetName() string
	GetNamespace() string
	GetHTTPURL() string
}

type Commit interface {
	common.Model
	GetHash() string
	GetMessage() string
	GetTimestamp() time.Time
	GetAuthor() Author
	GetURL() string
	GetAdded() []string
	GetModified() []string
	GetRemoved() []string
}

type Author interface {
	common.Model
	GetName() string
	GetEmail() string
}
