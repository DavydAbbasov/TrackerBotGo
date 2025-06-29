package bot

import (
	"github.com/DavydAbbasov/trecker_bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func New(config *config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = config.TelegramBotDebug

	log.Info().Msgf("avtorisotion as: %s", bot.Self.UserName)

	return bot, nil
}
