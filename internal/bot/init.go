package bot

import (
	"github.com/DavydAbbasov/trecker_bot/config"
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func New(cfg *config.Config) (interfaces.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = cfg.TelegramBotDebug

	log.Info().Msgf("avtorisotion as: %s", bot.Self.UserName)

	return bot, nil
}
