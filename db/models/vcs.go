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
	ProjectID int

	// fields
	State State
}

// Commit is an elementary change of the project in terms of VCS
type Commit struct {
	gorm.Model

	// relates to
	Project   Project
	ProjectID int
	Author    Author
	AuthorID  int
	Hook      Hook
	HookID    int

	// fields
	Hash    string
	Message string
	URL     string
	Author  Author
}

// Action helps to find out what exactly occured with a
// file/dir in a particular commit
type Action struct {
	gorm.Model

	Description string
}

// Path describes which file/dir was affected by a commit
type Path struct {
	gorm.Model

	// relates to
	Project   Project
	ProjectID int
	Commit    Commit
	CommitID  int
	Action    Action
	ActionID  int

	Name   string
	Action PathAction
}
