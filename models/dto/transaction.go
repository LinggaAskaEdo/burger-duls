package dto

import "time"

// UserResult struct
type UserResult struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

// MenuResult struct
type MenuResult struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Qty   int    `json:"qty"`
	Total int    `json:"total"`
}

// TransactionResult struct
type TransactionResult struct {
	ID                uint         `json:"id"`
	CreatedAt         time.Time    `json:"date"`
	TransactionNumber string       `json:"transactionNumber"`
	Invoice           string       `json:"invoice"`
	Confirm           bool         `json:"confirm"`
	User              UserResult   `json:"user"`
	Details           []MenuResult `json:"details"`
}
