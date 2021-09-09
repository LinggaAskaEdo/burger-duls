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

// UserController data type
type UserController struct {
	service   services.UserService
	logger    lib.Logger
	validator lib.Validator
}

// NewUserController creates new user controller
func NewUserController(userService services.UserService, logger lib.Logger, validator lib.Validator) UserController {
	return UserController{
		service:   userService,
		logger:    logger,
		validator: validator,
	}
}

type RegisterValidation struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
	Age      uint8  `json:"age" validate:"required,min=17,max=45"`
	Address  string `json:"address" validate:"required"`
}

type LoginValidation struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

// Register user
func (u UserController) Register(c *gin.Context) {
	u.logger.Info("Register route called")

	request := RegisterValidation{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&request); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	err := u.validator.Struct(request)
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

// Login user
func (u UserController) Login(c *gin.Context) {
	u.logger.Info("Login route called")

	request := LoginValidation{}

	if err := c.ShouldBindJSON(&request); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	err := u.validator.Struct(request)
	if err != nil {
		u.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error()})
		return
	}

	user, err := u.service.GetUserByEmailAndPassword(request.Email, request.Password)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
