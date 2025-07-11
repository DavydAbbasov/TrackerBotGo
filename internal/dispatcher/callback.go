package dispatcher

import (
	"strings"

	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallbackQuery(bot interfaces.StoppableBot, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	var replyText string
	switch data {
	case "lang_en":
		replyText = "Language set to English."
	case "lang_ru":
		replyText = "Язык установлен: русский."
	default:
		replyText = "Unknown option selected."
	}

	callbackResponse := tgbotapi.NewCallback(callback.ID, "")
	bot.Send(callbackResponse)

	//deletede old message (language)
	// deleteMsg := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
	// bot.Send(deleteMsg)

	if strings.HasPrefix(callback.Data, "lang_") {
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
		bot.Send(deleteMsg)
	}

	msg := tgbotapi.NewMessage(chatID, replyText)
	bot.Send(msg)

	ShowMainMenu(bot, chatID)
}
