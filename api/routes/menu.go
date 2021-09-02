package routes

import (
	"github.com/LinggaAskaEdo/burger-duls/api/controllers"
	"github.com/LinggaAskaEdo/burger-duls/lib"
)

// MenuRoutes struct
type MenuRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	menuController controllers.MenuController
}

// Setup menu routes
func (s MenuRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/burger-duls/menu")
	{
		api.POST("/add", s.menuController.AddMenu)
		api.GET("/all", s.menuController.AllMenu)
	}
}

// NewMenuRoutes creates new menu controller
func NewMenuRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	menuController controllers.MenuController,
) MenuRoutes {
	return MenuRoutes{
		handler:        handler,
		logger:         logger,
		menuController: menuController,
	}
}
