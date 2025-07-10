package main

import (
	"sync"

	"github.com/DavydAbbasov/trecker_bot/config"
	app "github.com/DavydAbbasov/trecker_bot/internal/application"
	log "github.com/rs/zerolog/log"
)

var wg sync.WaitGroup

func main() {
	cfg := config.MustLoadConfig(".env")

	bot := app.New(cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot.Run()
	}()

	wg.Wait()

	log.Info().Msg("telegram bot is shutdown")
}
