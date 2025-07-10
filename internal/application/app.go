package application

import (
	"github.com/DavydAbbasov/trecker_bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type App struct {
	bot *tgbotapi.BotAPI
	cfg *config.Config

	shutdown chan struct{}
}

func New(cfg *config.Config) *App {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal().Err(err).Msg("error to create telegram bot")
	}
	bot.Debug = cfg.TelegramBotDebug

	log.Info().Msgf("avtorisotion as: %s", bot.Self.UserName)

	return &App{
		bot: bot,
		cfg: cfg,
	}
}

func (a *App) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.bot.GetUpdatesChan(u)
	a.dispatch(updates)

	log.Info().Msg("telegram bot gracefully shutdown")
}

func (a *App) dispatch(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			a.handleCommand(update.Message)
		}
	}
}

func (a *App) Shutdown() {
}
