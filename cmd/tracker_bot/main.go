package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/DavydAbbasov/trecker_bot/application"
	"github.com/DavydAbbasov/trecker_bot/config"

	log "github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("application exited with error")
	}
	log.Info().Msg("shutdown complete")
}

func run() error { //?
	cfg := config.MustLoadConfig(".env")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM) //создает родительскиййй контекст который ловит сигналы от опирационной системы
	defer stop()

	app := application.New(cfg)

	return app.Run(ctx)
}
