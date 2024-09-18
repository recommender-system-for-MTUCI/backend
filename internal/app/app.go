package app

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
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
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Skipper:      middleware.DefaultSkipper,
			ErrorMessage: "many time, try again after some time",
			OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
				zap.Any("path", c.Path())
			},
			Timeout: time.Second * 7,
		}),
		// when will be add registration< need to change id for token
		middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 30, ExpiresIn: 2 * time.Minute},
			),
			IdentifierExtractor: func(ctx echo.Context) (string, error) {
				id := ctx.RealIP()
				return id, nil
			},
			ErrorHandler: func(context echo.Context, err error) error {
				return context.JSON(http.StatusForbidden, nil)
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return context.JSON(http.StatusTooManyRequests, nil)
			},
		}),
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
