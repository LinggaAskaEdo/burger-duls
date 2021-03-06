package routes

import (
	"github.com/LinggaAskaEdo/burger-duls/api/controllers"
	"github.com/LinggaAskaEdo/burger-duls/lib"
)

// UserRoutes struct
type UserRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	userController controllers.UserController
}

// Setup user routes
func (s UserRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/burger-duls/user")
	{
		api.POST("/register", s.userController.Register)
		api.POST("/login", s.userController.Login)
	}
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	userController controllers.UserController,
) UserRoutes {
	return UserRoutes{
		handler:        handler,
		logger:         logger,
		userController: userController,
	}
}
