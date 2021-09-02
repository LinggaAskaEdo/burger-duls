package services

import (
	"github.com/LinggaAskaEdo/burger-duls/lib"
	entity "github.com/LinggaAskaEdo/burger-duls/models/entity"
	"github.com/LinggaAskaEdo/burger-duls/repository"
	"gorm.io/gorm"
)

// MenuService service layer
type MenuService struct {
	logger     lib.Logger
	repository repository.MenuRepository
}

// NewMenuService creates a new userservice
func NewMenuService(logger lib.Logger, repository repository.MenuRepository) MenuService {
	return MenuService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s MenuService) WithTrx(trxHandle *gorm.DB) MenuService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// AddMenu call to add the menu
func (s MenuService) AddMenu(menu entity.Menu) (result entity.Menu, err error) {
	return menu, s.repository.Create(&menu).Error
}

// GetAllMenu call to get all menu
func (s MenuService) GetAllMenu() (menus []entity.Menu, err error) {
	return menus, s.repository.Find(&menus).Error
}
