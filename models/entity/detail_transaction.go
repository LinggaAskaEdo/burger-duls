package models

// Deail Transaction model
type DetailTransaction struct {
	ID            uint `gorm:"primarykey"`
	TransactionID uint
	Transaction   Transaction
	MenuID        uint
	Menu          Menu
	Qty           int `gorm:"NOT NULL"`
	Total         int `gorm:"NOT NULL"`
}

// TableName gives table name of model
func (dt DetailTransaction) TableName() string {
	return "detail_transaction"
}
