package transport

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/backend/internal/app"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"github.com/recommender-system-for-MTUCI/backend/internal/pkg/jwt"
	"go.uber.org/zap"
)

func ReccomendSystem() {
	var (
		server *app.Controller
		log    *zap.Logger
		err    error
		cfg    *config.Config
		ctx    context.Context
	)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()
	log, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Falied to initialize zap logger", zap.Error(err))
	}
	log.Info("Initilaze logger")
	cfg, err = config.NewConfig()
	if err != nil {
		log.Fatal("Failed to initilize config", zap.Error(err))
	}
	log.Info("Initialize config", zap.Any("config", cfg))

	prov, err := jwt.NewProvider(cfg.JWT, log)
	if err != nil {
		log.Fatal("Failed to create jwt provider", zap.Error(err))
	}
	id, err := uuid.NewUUID()
	if err == nil {
		got, err := prov.CreateTokenForUser(id, false)
		log.Info("generated token", zap.String("token", got), zap.Error(err))
		parsed, err := prov.GetDataFromToken(got)
		log.Info("parsed from token", zap.String("token", parsed.String()), zap.Error(err))
	}

	server, err = app.New(log, cfg)
	if err != nil {
		log.Fatal("Failed to initialize server", zap.Error(err))
	}

	log.Info("Initilize server", zap.Any("server", server))
	defer func() {
		log.Error(
			"Shutting down server",
			zap.Error(server.ShutDown(ctx)),
		)
	}()
	err = server.Run(ctx)
	if err != nil {
		log.Fatal("failed to initialize server")
	}
	<-ctx.Done()
}
