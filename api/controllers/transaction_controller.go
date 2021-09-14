package controllers

import (
	"net/http"
	"strconv"

	"github.com/LinggaAskaEdo/burger-duls/constants"
	"github.com/LinggaAskaEdo/burger-duls/lib"
	"github.com/LinggaAskaEdo/burger-duls/models/dto"
	entity "github.com/LinggaAskaEdo/burger-duls/models/entity"
	"github.com/LinggaAskaEdo/burger-duls/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

// TransactionController data type
type TransactionController struct {
	service   services.TransactionService
	logger    lib.Logger
	validator lib.Validator
	env       lib.Env
}

// NewTransactionController creates new transaction controller
func NewTransactionController(transactionService services.TransactionService, logger lib.Logger, validator lib.Validator, env lib.Env) TransactionController {
	return TransactionController{
		service:   transactionService,
		logger:    logger,
		validator: validator,
		env:       env,
	}
}

type DetailValidation struct {
	MenuId uint `json:"menuId" validate:"required"`
	Qty    int  `json:"qty" validate:"required"`
	Price  int  `json:"price" validate:"required"`
}

type OrderValidation struct {
	UserId  uint               `json:"userId" validate:"required"`
	Details []DetailValidation `json:"details" validate:"required,dive,required"`
}

type InvoiceValidation struct {
	TransactionNumber string `json:"transactionNumber" validate:"required"`
	InvoiceImage      string `json:"invoiceImage" validate:"required"`
}

// Order
func (t TransactionController) Order(c *gin.Context) {
	t.logger.Info("Order route called")

	request := OrderValidation{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	// unmarshal incoming input into pre-defined structure
	if err := c.ShouldBindJSON(&request); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Just terminate the request if the input is not valid
	err := t.validator.Struct(request)
	if err != nil {
		t.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error()})
		return
	}

	// Map to domain entity
	transaction := entity.Transaction{
		TransactionNumber: xid.New().String(),
		UserID:            request.UserId,
		Confirm:           false}

	result, err := t.service.WithTrx(trxHandle).StoreTransaction(transaction)
	if err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	var detailList []entity.DetailTransaction

	for _, dt := range request.Details {
		detailList = append(detailList, entity.DetailTransaction{
			TransactionID: result.ID,
			MenuID:        dt.MenuId,
			Qty:           dt.Qty,
			Total:         dt.Qty * dt.Price})
	}

	nil, err := t.service.WithTrx(trxHandle).StoreDetailTransaction(detailList)
	if err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusCreated, "message": "Request Was Successfully"})
}

// AllTransaction
func (t TransactionController) AllTransaction(c *gin.Context) {
	t.logger.Info("AllTransaction route called")

	transactions := make([]*dto.TransactionResult, 0)

	results, err := t.service.GetTransactions()
	if err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	defer results.Close()

	// construct general transaction
	for results.Next() {
		transaction := new(dto.TransactionResult)
		results.Scan(
			&transaction.ID, &transaction.CreatedAt, &transaction.TransactionNumber, &transaction.Invoice, &transaction.Confirm,
			&transaction.User.ID, &transaction.User.Name, &transaction.User.Email, &transaction.User.Age, &transaction.User.Address, &transaction.User.Phone)
		transactions = append(transactions, transaction)
	}

	// construct transaction detail menu
	for _, dt := range transactions {
		results, err := t.service.GetMenus(dt.ID)
		if err != nil {
			t.logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  err.Error(),
			})
			return
		}

		for results.Next() {
			menu := new(dto.MenuResult)
			results.Scan(&menu.ID, &menu.Name, &menu.Qty, &menu.Total)
			dt.Details = append(dt.Details, *menu)
		}
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// TransactionByUserId
func (t TransactionController) TransactionByUserId(c *gin.Context) {
	t.logger.Info("TransactionByUserId route called")

	transactions := make([]*dto.TransactionResult, 0)

	userId, _ := strconv.Atoi(c.Param("id"))
	t.logger.Debug("userId:", userId)

	results, err := t.service.GetTransactionsByUserId(userId)
	if err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	defer results.Close()

	// construct general transaction
	for results.Next() {
		transaction := new(dto.TransactionResult)
		results.Scan(
			&transaction.ID, &transaction.CreatedAt, &transaction.TransactionNumber, &transaction.Invoice, &transaction.Confirm,
			&transaction.User.ID, &transaction.User.Name, &transaction.User.Email, &transaction.User.Age, &transaction.User.Address, &transaction.User.Phone)
		transactions = append(transactions, transaction)
	}

	// construct transaction detail menu
	for _, dt := range transactions {
		results, err := t.service.GetMenus(dt.ID)
		if err != nil {
			t.logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  err.Error(),
			})
			return
		}

		for results.Next() {
			menu := new(dto.MenuResult)
			results.Scan(&menu.ID, &menu.Name, &menu.Qty, &menu.Total)
			dt.Details = append(dt.Details, *menu)
		}
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// Invoice
func (t TransactionController) Invoice(c *gin.Context) {
	t.logger.Info("Invoice route called")

	request := InvoiceValidation{}

	// unmarshal incoming input into pre-defined structure
	if err := c.ShouldBindJSON(&request); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Just terminate the request if the input is not valid
	err := t.validator.Struct(request)
	if err != nil {
		t.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error()})
		return
	}

	// Check Valid Transaction Number

	transExist, err := t.service.CheckTransactionNumber(request.TransactionNumber)
	if err != nil {
		t.logger.Error("Error when CheckTransactionNumber")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	t.logger.Debug("TransactionNumber: ", transExist.TransactionNumber)

	if transExist.TransactionNumber != request.TransactionNumber {
		t.logger.Error("Invalid transaction number")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusNotFound,
			"error":  "Invalid transaction number",
		})
		return
	}

	bucketName := t.env.BucketName
	projectId := t.env.ProjectId
	imageName := request.TransactionNumber

	t.logger.Debug(bucketName + ", " + projectId + ", " + imageName)

	if err := t.service.UploadFile(bucketName, projectId, imageName, request.InvoiceImage); err != nil {
		t.logger.Error("Error when UploadFile")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	imageUrl := t.env.URLCloudStorage + imageName
	transExist.Invoice = imageUrl

	if err := t.service.UpdateInvoiceUrl(&transExist); err != nil {
		t.logger.Error("Error when UpdateInvoiceUrl")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"message":  "success",
		"imageUrl": imageUrl,
	})
}
