package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Controller struct {
	server *echo.Echo
	log    *zap.Logger
}

func New(log *zap.Logger, err error) (*Controller, error) {
	log.Info("Initialize serverr")
	ctrl := &Controller{
		server: echo.New(),
		log:    log,
	}
	ctrl.configure()
	return ctrl, nil
}

func (ctrl *Controller) configure() error {
	ctrl.configureMiddlewares()
	ctrl.configureRouters()
	return nil
}

func (ctrl *Controller) configureMiddlewares() {
	//need add some middlewares
	//need correct Cors
	var middlewares = []echo.MiddlewareFunc{
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodDelete, http.MethodGet, http.MethodPost, http.MethodPatch},
		}),
		middleware.Recover(),
		middleware.Logger(),
	}
	ctrl.server.Use(middlewares...)
}

func (ctrl *Controller) configureRouters() {
	app := ctrl.server.Group("/app")
	app.POST("/form", nil)
}
