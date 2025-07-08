package dispatcher

import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

/*
handlers - исполнители (обработчики команд)
Реализуют, что делать по команде (/start, /help, /track)
Отвечают за бизнес-логику команды
(например, "Привет, я бот" или "добавить активность")
*/
func HandleStart(bot interfaces.StoppableBot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := "Finaly"

	message := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("Error when sending the /start command")
	}
}
