package services

import (
	"github.com/LinggaAskaEdo/burger-duls/lib"
	entity "github.com/LinggaAskaEdo/burger-duls/models/entity"
	"github.com/LinggaAskaEdo/burger-duls/repository"
	"gorm.io/gorm"
)

// UserService service layer
type UserService struct {
	logger     lib.Logger
	repository repository.UserRepository
}

// NewUserService creates a new userservice
func NewUserService(logger lib.Logger, repository repository.UserRepository) UserService {
	return UserService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s UserService) WithTrx(trxHandle *gorm.DB) UserService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// CreateUser call to create the user
func (s UserService) CreateUser(user entity.User) (result entity.User, err error) {
	return user, s.repository.Create(&user).Error
}

func (s UserService) GetUserByEmailAndPassword(email string, password string) (user entity.User, err error) {
	return user, s.repository.First(&user, "email = ? AND password = ?", email, password).Error
}
