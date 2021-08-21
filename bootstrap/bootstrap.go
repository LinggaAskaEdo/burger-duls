package bootstrap

import (
	"context"
	"time"

	"github.com/LinggaAskaEdo/burger-duls/api/controllers"
	"github.com/LinggaAskaEdo/burger-duls/api/middlewares"
	"github.com/LinggaAskaEdo/burger-duls/api/routes"
	"github.com/LinggaAskaEdo/burger-duls/lib"
	"github.com/LinggaAskaEdo/burger-duls/repository"
	"github.com/LinggaAskaEdo/burger-duls/services"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	controllers.Module,
	routes.Module,
	lib.Module,
	services.Module,
	middlewares.Module,
	repository.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler lib.RequestHandler,
	routes routes.Routes,
	env lib.Env,
	logger lib.Logger,
	middlewares middlewares.Middlewares,
	database lib.Database,
) {
	conn, _ := database.DB.DB() // trigger database connection

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting Application")
			logger.Info("---------------------")
			logger.Info("------- CLEAN -------")
			logger.Info("---------------------")

			conn.SetMaxOpenConns(env.DBMaxOpenConn)
			conn.SetMaxIdleConns(env.DBMaxIdleConn)
			conn.SetConnMaxLifetime(time.Duration(env.DBMaxLifeTime))
			conn.SetConnMaxIdleTime(time.Duration(env.DBMaxIdleTime))

			go func() {
				middlewares.Setup()
				routes.Setup()
				handler.Gin.Run(":" + env.ServerPort)
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Stopping Application")
			conn.Close()
			return nil
		},
	})
}
