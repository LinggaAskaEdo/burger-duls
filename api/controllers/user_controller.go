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

// UserController data type
type UserController struct {
	service services.UserService
	logger  lib.Logger
}

// NewUserController creates new user controller
func NewUserController(userService services.UserService, logger lib.Logger) UserController {
	return UserController{
		service: userService,
		logger:  logger,
	}
}

// Logout user
func (u UserController) Register(c *gin.Context) {
	u.logger.Info("Register route called")

	request := dto.Request{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&request); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	type RegisterValidation struct {
		Name     string `validate:"required,min=5,max=50"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=5"`
		Age      uint8  `validate:"required,min=17,max=45"`
		Address  string `validate:"required"`
	}

	registerValidation := &RegisterValidation{Name: request.Name, Email: request.Email, Password: request.Password, Age: request.Age, Address: request.Address}

	err := validator.New().Struct(registerValidation)
	if err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error()})
		return
	}

	user := entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Age:      request.Age,
		Address:  request.Address}

	result, err := u.service.WithTrx(trxHandle).CreateUser(user)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": result})
}
