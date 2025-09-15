package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/DavydAbbasov/trecker_bot/config"
	"github.com/DavydAbbasov/trecker_bot/internal/application"

	log "github.com/rs/zerolog/log"
)

func main() {
	cfg := config.MustLoadConfig(".env")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := application.New(cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("application exited with error")
	}
	log.Info().Msg("shutdown complete")
}
