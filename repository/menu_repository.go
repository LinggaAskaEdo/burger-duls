package repository

import (
	"github.com/LinggaAskaEdo/burger-duls/lib"
	"gorm.io/gorm"
)

// MenuRepository database structure
type MenuRepository struct {
	lib.Database
	logger lib.Logger
}

// NewMenuRepository creates a new user repository
func NewMenuRepository(db lib.Database, logger lib.Logger) MenuRepository {
	return MenuRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository with transaction
func (r MenuRepository) WithTrx(trxHandle *gorm.DB) MenuRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}
