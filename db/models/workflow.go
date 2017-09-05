package models

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model

	Group string
	Name  string
}
