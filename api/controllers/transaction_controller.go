package controllers

import (
	"net/http"

	"github.com/LinggaAskaEdo/burger-duls/constants"
	"github.com/LinggaAskaEdo/burger-duls/lib"
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
}

// NewTransactionController creates new transaction controller
func NewTransactionController(transactionService services.TransactionService, logger lib.Logger, validator lib.Validator) TransactionController {
	return TransactionController{
		service:   transactionService,
		logger:    logger,
		validator: validator,
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
