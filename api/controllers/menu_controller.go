package controllers

import (
	"net/http"

	"github.com/LinggaAskaEdo/burger-duls/constants"
	"github.com/LinggaAskaEdo/burger-duls/lib"
	dto "github.com/LinggaAskaEdo/burger-duls/models/dto"
	entity "github.com/LinggaAskaEdo/burger-duls/models/entity"
	"github.com/LinggaAskaEdo/burger-duls/services"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

// MenuController data type
type MenuController struct {
	service services.MenuService
	logger  lib.Logger
}

// NewMenuController creates new menu controller
func NewMenuController(menuService services.MenuService, logger lib.Logger) MenuController {
	return MenuController{
		service: menuService,
		logger:  logger,
	}
}

// AddMenu user
func (m MenuController) AddMenu(c *gin.Context) {
	m.logger.Info("AddMenu route called")

	request := dto.Request{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&request); err != nil {
		m.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	type AddMenuValidation struct {
		Name        string `validate:"required,min=3,max=50"`
		Description string `validate:"required"`
		Price       int    `validate:"required,min=5000"`
		Type        string `validate:"required"`
	}

	addMenuValidation := &AddMenuValidation{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Type:        request.Type}

	err := validator.New().Struct(addMenuValidation)
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
