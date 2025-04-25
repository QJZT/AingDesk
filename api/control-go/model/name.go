package model

import (
	"gorm.io/gorm"
)

type Name struct {
	gorm.Model
	Name string `json:"name"`
}
