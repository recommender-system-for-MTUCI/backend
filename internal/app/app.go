package app

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"github.com/recommender-system-for-MTUCI/backend/internal/pkg"
	"github.com/recommender-system-for-MTUCI/backend/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Controller struct {
	server        *echo.Echo
	log           *zap.Logger
	cfg           *config.Config
	tokenProvider pkg.Provider
	repo          storage.Storage
}

func New(log *zap.Logger, cfg *config.Config, tokenProvider pkg.Provider, repo storage.Storage) (*Controller, error) {
	log.Info("Initializing server")
	ctrl := &Controller{
		server:        echo.New(),
		log:           log,
		cfg:           cfg,
		tokenProvider: tokenProvider,
		repo:          repo,
	}
	err := ctrl.configure()
	if err != nil {
		ctrl.log.Error("got error when configure server")
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
	// CORS configurations need change when frontend will be ready
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
		// when will be add registration need to change id for token
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
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Skipper:      middleware.DefaultSkipper,
			Generator:    uuid.NewString,
			TargetHeader: echo.HeaderXRequestID,
		}),
	}

	ctrl.server.Use(middlewares...)
}

func (ctrl *Controller) configureRouters() {
	// Define routes
	app := ctrl.server.Group("/app")
	app.POST("/registration", ctrl.handleRegistration)
}

func (ctrl *Controller) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		ctrl.log.Info("starting HTTP server on address", zap.String("", ctrl.cfg.Server.GetServerAddress()))
		err := ctrl.server.Start(ctrl.cfg.Server.GetServerAddress())
		if err != nil {
			cancel()
		}
	}()
	return ctx.Err()
}

func (ctrl *Controller) ShutDown(ctx context.Context) error {
	return ctrl.server.Shutdown(ctx)
}
