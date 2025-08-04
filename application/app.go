// инициализирует, запускает, делегирует.
package application

import (
	"context"

	"github.com/DavydAbbasov/trecker_bot/fsm"
	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher"
	handlers "github.com/DavydAbbasov/trecker_bot/internal/dispatcher"
	"github.com/DavydAbbasov/trecker_bot/storage"

	"github.com/DavydAbbasov/trecker_bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type App struct {
	bot             interfaces.BotAPI
	cfg             *config.Config
	dispatcher      *dispatcher.Dispatcher
	flushables      []interfaces.Flushable //?
	activityStorage storage.ActivityStorage
	learningStorage storage.LearningStorage
}

func New(cfg *config.Config) *App {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal().Err(err).Msg("error to create telegram bot")
	}

	bot.Debug = cfg.TelegramBotDebug
	log.Info().Msgf("avtorisation as: %s", bot.Self.UserName)

	activityStorage := storage.NewMemoryActivityStorage()
	learningStorage := storage.NewMemoryLearningStorage()
	fsmManager := fsm.NewFSM()                                                    //
	dispatcher := handlers.New(bot, fsmManager, activityStorage, learningStorage) //

	flushables := []interfaces.Flushable{ //?
		dispatcher,
		// to do
	}

	return &App{
		bot:             bot,
		cfg:             cfg,
		dispatcher:      dispatcher,
		flushables:      flushables,
		activityStorage: activityStorage,
	}
}
func (a *App) Run(ctx context.Context) error { //?
	go a.dispatcher.Run() //

	<-ctx.Done() //

	log.Info().Msg("shutdown initiated")
	a.bot.StopReceivingUpdates()

	for _, f := range a.flushables { //?
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
