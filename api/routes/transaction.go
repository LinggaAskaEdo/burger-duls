package routes

import (
	"github.com/LinggaAskaEdo/burger-duls/api/controllers"
	"github.com/LinggaAskaEdo/burger-duls/lib"
)

// TransactionRoutes struct
type TransactionRoutes struct {
	logger                lib.Logger
	handler               lib.RequestHandler
	transactionController controllers.TransactionController
}

// Setup transaction routes
func (t TransactionRoutes) Setup() {
	t.logger.Info("Setting up routes")
	api := t.handler.Gin.Group("/burger-duls/order")
	{
		api.POST("/add", t.transactionController.Order)
	}
}

// NewTransactionRoutes creates new transaction controller
func NewTransactionRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	transactionController controllers.TransactionController,
) TransactionRoutes {
	return TransactionRoutes{
		handler:               handler,
		logger:                logger,
		transactionController: transactionController,
	}
}
