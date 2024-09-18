package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"go.uber.org/zap"
)

type Controller struct {
	server *echo.Echo
	log    *zap.Logger
	cfg    *config.Config
}

func New(log *zap.Logger, cfg *config.Config) (*Controller, error) {
	log.Info("Initializing server")
	ctrl := &Controller{
		server: echo.New(),
		log:    log,
		cfg:    cfg,
	}
	err := ctrl.configure()
	if err != nil {
		return nil, err
	}
	return ctrl, nil
}

func (ctrl *Controller) configure() error {
	ctrl.configureMiddlewares()
	ctrl.configureRouters()
	return nil
}

func (ctrl *Controller) configureMiddlewares() {
	// CORS configuration
	var middlewares = []echo.MiddlewareFunc{
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodPatch,
			},
		}),
		middleware.Recover(),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogStatus:    true,
			LogMethod:    true,
			LogURIPath:   true,
			LogError:     true,
			LogRequestID: true,
			LogUserAgent: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {

				return nil

			},
		}),
		middleware.Logger(),
	}

	ctrl.server.Use(middlewares...)
}

func (ctrl *Controller) configureRouters() {
	// Define routes
	app := ctrl.server.Group("/app")
	app.GET("/hi", ctrl.handleFormSubmission)
}

func (ctrl *Controller) handleFormSubmission(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Form submitted successfully!",
	})
}

func (ctrl *Controller) Run() error {
	//ctx, cancel := context.WithCancel(c)
	//go func() {
	//	ctrl.log.Info("starting HTTP server on address", zap.String("", ctrl.cfg.Server.GetServerAddress()))
	//	err := ctrl.server.Start(ctrl.cfg.Server.GetServerAddress())
	//	if err != nil {
	//		cancel()
	//	}
	//}()
	err := ctrl.server.Start(ctrl.cfg.Server.GetServerAddress())
	ctrl.log.Info(ctrl.cfg.Server.GetServerAddress())
	return err
}
