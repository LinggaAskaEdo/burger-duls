package models

import (
	"gorm.io/gorm"
)

// User model
type Menu struct {
	gorm.Model
	Name        string `gorm:"NOT NULL"`
	Description string `gorm:"NOT NULL"`
	Price       int    `gorm:"NOT NULL"`
	Type        string `gorm:"NOT NULL"`
}

// TableName gives table name of model
func (m Menu) TableName() string {
	return "menu"
}
