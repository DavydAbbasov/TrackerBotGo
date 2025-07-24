// инициализирует, запускает, делегирует.
package application

import (
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher"
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"

	"github.com/DavydAbbasov/trecker_bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type App struct {
	bot      interfaces.BotAPI
	cfg      *config.Config
	dispatcher *dispatcher.Dispatcher
	shutdown chan struct{}
}

func New(cfg *config.Config) *App {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal().Err(err).Msg("error to create telegram bot")
	}

	bot.Debug = cfg.TelegramBotDebug
	log.Info().Msgf("avtorisation as: %s", bot.Self.UserName)

	dispatcher := dispatcher.New(bot)

	return &App{
		bot: bot,
		cfg: cfg,
		dispatcher: dispatcher,
	}
}
func (a *App) Run() {
	a.dispatcher.Run()

	log.Info().Msg("telegram bot gracefully shutdown")
}

// func (a *App)dispatcher(updates tgbotapi.UpdatesChannel){
// 	for 
// }

// func(d *Dispatcher)Dispatcher(updates tgbotapi.UpdatesChannel){
	
// }