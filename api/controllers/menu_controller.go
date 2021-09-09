package controllers

import (
	"net/http"

	"github.com/LinggaAskaEdo/burger-duls/constants"
	"github.com/LinggaAskaEdo/burger-duls/lib"
	entity "github.com/LinggaAskaEdo/burger-duls/models/entity"
	"github.com/LinggaAskaEdo/burger-duls/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuController data type
type MenuController struct {
	service   services.MenuService
	logger    lib.Logger
	validator lib.Validator
}

// NewMenuController creates new menu controller
func NewMenuController(menuService services.MenuService, logger lib.Logger, validator lib.Validator) MenuController {
	return MenuController{
		service:   menuService,
		logger:    logger,
		validator: validator,
	}
}

type AddMenuValidation struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required"`
	Price       int    `json:"price" validate:"required,min=5000"`
	Type        string `json:"type" validate:"required"`
}

// AddMenu user
func (m MenuController) AddMenu(c *gin.Context) {
	m.logger.Info("AddMenu route called")

	request := AddMenuValidation{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&request); err != nil {
		m.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	err := m.validator.Struct(request)
	if err != nil {
		m.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error()})
		return
	}

	menu := entity.Menu{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Type:        request.Type}

	result, err := m.service.WithTrx(trxHandle).AddMenu(menu)
	if err != nil {
		m.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"menu": result})
}

// AllMenu user
func (m MenuController) AllMenu(c *gin.Context) {
	m.logger.Info("AllMenu route called")

	menu, err := m.service.GetAllMenu()
	if err != nil {
		m.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"menus": menu})
}
