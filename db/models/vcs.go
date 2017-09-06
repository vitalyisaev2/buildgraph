package models

import (
	"github.com/jinzhu/gorm"
)

// Author represents a person who owns the commit
type Author struct {
	gorm.Model

	Name  string
	Email string
}

// Hook is an event emitted by VCS
type Hook struct {
	gorm.Model

	// relates to
	Project   Project
	ProjectID uint

	// fields
	Commits []Commit
}

// Commit is an elementary change of the project in terms of Git VCS
type Commit struct {
	gorm.Model

	// relates to
	Author   Author
	AuthorID uint
	HookID   uint
	Changes  []PathChange

	// fields
	Hash    string
	Message string
	URL     string
}

// Path describes which file/dir was affected by a commit
type Path struct {
	gorm.Model

	// relates to
	Project   Project
	ProjectID uint
	// fields
	Name string
}

// PathEvent shows what exactly changed on path
type PathEvent struct {
	gorm.Model
	// fields
	Name string
}

// PathChange struct
type PathChange struct {
	gorm.Model

	// relates to
	Commit   Commit
	CommitID uint
	Path     Path
	PathID   uint
	Event    PathEvent
	EventID  uint
}
