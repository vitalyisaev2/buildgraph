package models

import "github.com/jinzhu/gorm"

type Build struct {
	gorm.Model

	// relates to
	ProjectContext   ProjectContext
	ProjectContextID int

	// fields
	CIQueueID int // Assigned by CI server
	CIBuildID int // Assigned by CI server
	Result    string
	URL       string
}
