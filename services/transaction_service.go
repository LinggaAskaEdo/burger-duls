package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"cloud.google.com/go/storage"
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

// GetAllTransaction call to get all transactions
func (t TransactionService) GetTransactions() (rows *sql.Rows, err error) {
	return t.repository.
		Table("transaction").
		Select("transaction.id, transaction.created_at, transaction.transaction_number, transaction.invoice, transaction.confirm, user.id, user.name, user.email, user.age, user.address, user.phone").
		Joins("inner join user on user.id = transaction.user_id").
		Rows()
}

// GetMenus call to get menus of transaction
func (t TransactionService) GetMenus(transactionId uint) (rows *sql.Rows, err error) {
	return t.repository.
		Table("detail_transaction").
		Select("detail_transaction.menu_id, menu.name, detail_transaction.qty, detail_transaction.total").
		Joins("inner join menu on menu.id = detail_transaction.menu_id").
		Where("detail_transaction.transaction_id = ?", transactionId).
		Rows()
}

// GetTransactionsByUserId call to get all transactions
func (t TransactionService) GetTransactionsByUserId(userId int) (rows *sql.Rows, err error) {
	return t.repository.
		Table("transaction").
		Select("transaction.id, transaction.created_at, transaction.transaction_number, transaction.invoice, transaction.confirm, user.id, user.name, user.email, user.age, user.address, user.phone").
		Joins("inner join user on user.id = transaction.user_id").
		Where("transaction.user_id = ?", userId).
		Rows()
}

// CheckTransactionNumber call to check transaction number is exist or not
func (t TransactionService) CheckTransactionNumber(transactionNumber string) (transaction entity.Transaction, err error) {
	return transaction, t.repository.Limit(1).Find(&transaction, "transaction_number = ?", transactionNumber).Error
}

// UploadFile
func (t TransactionService) UploadFile(bucketName string, projectId string, imageName string, invoiceImage string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.logger.Error("Unable initate storage")
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	// Upload an object with storage.Writer.
	bkt := client.Bucket(bucketName).UserProject(projectId)
	obj := bkt.Object(imageName)
	wc := obj.NewWriter(ctx)
	wc.ContentType = "image/jpeg"

	//Encode from image format to writer
	dec, err := base64.StdEncoding.DecodeString(invoiceImage)
	if err != nil {
		t.logger.Error("Unable read Base64 string")
		return errors.New("Unable read Base64 string: " + err.Error())
	}

	if _, err = wc.Write(dec); err != nil {
		t.logger.Error("Unable to write data to bucket")
		return errors.New("Unable to write data to bucket: " + err.Error())
	}

	if err := wc.Close(); err != nil {
		t.logger.Error("Unable to close bucket")
		return errors.New("Unable to close bucket: " + err.Error())
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		t.logger.Error("Unable to put ACL for the named file")
		return errors.New("Unable to put ACL for the named file: " + err.Error())
	}

	return nil
}

// UpdateInvoiceUrl call to update invoice image url
func (t TransactionService) UpdateInvoiceUrl(transaction *entity.Transaction) (err error) {
	return t.repository.Save(&transaction).Error
}
