package services

import (
	"github.com/LinggaAskaEdo/burger-duls/lib"
	entity "github.com/LinggaAskaEdo/burger-duls/models/entity"
	"github.com/LinggaAskaEdo/burger-duls/repository"
	"gorm.io/gorm"
)

// TransactionService service layer
type TransactionService struct {
	logger     lib.Logger
	repository repository.TransactionRepository
}

// NewTransactionService creates a new transactionservice
func NewTransactionService(logger lib.Logger, repository repository.TransactionRepository) TransactionService {
	return TransactionService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (t TransactionService) WithTrx(trxHandle *gorm.DB) TransactionService {
	t.repository = t.repository.WithTrx(trxHandle)
	return t
}

// StoreTransaction call to store the transaction
func (t TransactionService) StoreTransaction(transaction entity.Transaction) (result entity.Transaction, err error) {
	return transaction, t.repository.Create(&transaction).Error
}

// StoreDetailTransaction call to store the detail transaction
func (t TransactionService) StoreDetailTransaction(detailTransaction []entity.DetailTransaction) (nil, err error) {
	return nil, t.repository.Create(&detailTransaction).Error
}
