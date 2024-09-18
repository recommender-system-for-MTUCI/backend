package transport

import (
	"github.com/recommender-system-for-MTUCI/backend/internal/app"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"go.uber.org/zap"
)

func ReccomendSystem() {
	var (
		server *app.Controller
		log    *zap.Logger
		err    error
		cfg    *config.Config
		//ctx    context.Context
	)
	log, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Falied to initialize zap logger", zap.Error(err))
	}
	log.Info("Initilaze logger")
	cfg, err = config.New()
	if err != nil {
		log.Fatal("Failed to initilize config", zap.Error(err))
	}
	log.Info("Initialize config", zap.Any("config", cfg))

	server, err = app.New(log, cfg)
	if err != nil {
		log.Fatal("Failed to initialize server", zap.Error(err))
	}

	log.Info("Initilize server", zap.Any("server", server))
	err = server.Run()
	if err != nil {
		log.Fatal("failed to initialize server")
	}

}
