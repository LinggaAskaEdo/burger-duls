package repository

import (
	"github.com/LinggaAskaEdo/burger-duls/lib"
	"gorm.io/gorm"
)

// TransactionRepository database structure
type TransactionRepository struct {
	lib.Database
	logger lib.Logger
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db lib.Database, logger lib.Logger) TransactionRepository {
	return TransactionRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository with transaction
func (r TransactionRepository) WithTrx(trxHandle *gorm.DB) TransactionRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}

	r.Database.DB = trxHandle
	return r
}
