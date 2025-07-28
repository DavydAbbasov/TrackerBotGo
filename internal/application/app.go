// инициализирует, запускает, делегирует.
package application

import (
	"context"

	"github.com/DavydAbbasov/trecker_bot/internal/handlers"
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"

	"github.com/DavydAbbasov/trecker_bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type App struct {
	bot        interfaces.BotAPI
	cfg        *config.Config
	handlers   *handlers.Dispatcher
	flushables []interfaces.Flushable //?
}

func New(cfg *config.Config) *App {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal().Err(err).Msg("error to create telegram bot")
	}

	bot.Debug = cfg.TelegramBotDebug
	log.Info().Msgf("avtorisation as: %s", bot.Self.UserName)

	handlers := handlers.New(bot)
	flushables := []interfaces.Flushable{ //?
		handlers,
		// в будущем: db, cache, queue
	}

	return &App{
		bot:        bot,
		cfg:        cfg,
		handlers:   handlers,
		flushables: flushables,
	}
}
func (a *App) Run(ctx context.Context) error { //?
	go a.handlers.Run()

	<-ctx.Done()

	log.Info().Msg("shutdown initiated")
	a.bot.StopReceivingUpdates()

	for _, f := range a.flushables {
		if err := f.Flush(); err != nil {
			log.Error().Err(err).Msg("flush failed")
		}
		if err := f.Close(); err != nil {
			log.Error().Err(err).Msg("close failed")
		}
	}

	log.Info().Msg("Shutdown complete")
	return nil
}
