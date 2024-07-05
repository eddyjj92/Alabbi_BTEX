package models

import (
	"gorm.io/gorm"
)

type Process struct {
	gorm.Model
	File      string
	Extension string
	Folder    string
	Route     string
	Status    int8
}
