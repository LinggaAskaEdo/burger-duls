package models

import "gorm.io/gorm"

// Transaction model
type Transaction struct {
	gorm.Model
	TransactionNumber string `gorm:"NOT NULL"`
	UserID            uint
	User              User
	Confirm           bool `gorm:"NOT NULL"`
}

// TableName gives table name of model
func (t Transaction) TableName() string {
	return "transaction"
}
