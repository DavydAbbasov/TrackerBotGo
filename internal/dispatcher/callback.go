package dispatcher

import (
	"strings"

	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallbackQuery(bot interfaces.BotAPI, callback *tgbotapi.CallbackQuery) {
	if callback.Message == nil {
		return
	}
	chatID := callback.Message.Chat.ID
	data := callback.Data

	//Обработка спецкоманд до switch
	switch {
	case data == "edit_language":
		ShowLanguageSelection(bot, chatID)
		return

	case strings.HasPrefix(data, "activity_report_"):
		activityName := strings.TrimPrefix(data, "activity_report_")
		userID := callback.From.ID
		ShowActivityReport(bot, chatID, userID, activityName)
		return

	case strings.HasPrefix(data, "calendar_"):
		activity := strings.TrimPrefix(data, "calendar_")
		ShowCalendar(bot, chatID, activity)
		return
	}

	var replyText string

	switch data {
	case "lang_en":
		replyText = "Language set to English."
	case "lang_ru":
		replyText = "Язык установлен: русский."
	case "lang_dch":
		replyText = "Sprache eingestellt: Deutsch."
	case "lang_ukr":
		replyText = "Мову встановлено: українська."
	case "lang_arab":
		replyText = "تم تعيين اللغة: العربية."
	case "lang_tur":
		replyText = "Dil ayarlandı: Türkçe."

		//return только для тех case, где есть вызов функции
	case "refresh_profile":
		ShowEditProfileMenu(bot, callback.Message.Chat.ID)
		return // Останавливает выполнение, дальше код не идёт

	case "summary":
		userID := callback.From.ID
		ShowActivityList(bot, callback.Message.Chat.ID, userID)
		return

	case "create_activity":
		AddActivity(bot, callback.Message.Chat.ID)
		return

	case "selection_activity":
		userID := callback.From.ID
		SelectionActivityPromt(bot, callback.Message.Chat.ID, userID)
		return
	case "add_collection":
		AddCollection(bot, callback.Message.Chat.ID)
		return
	case "switch_learning_actv":
		userID := callback.From.ID
		SowUserCollections(bot, callback.Message.Chat.ID, userID)
		return

	default:
		replyText = "Unknown option selected."
	}

	callbackResponse := tgbotapi.NewCallback(callback.ID, "")
	bot.Send(callbackResponse)

	if strings.HasPrefix(callback.Data, "lang_") {
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
		bot.Send(deleteMsg)
	}

	msg := tgbotapi.NewMessage(chatID, replyText)
	bot.Send(msg)

	// ShowMainMenu(bot, chatID)

	if strings.HasPrefix(data, "lang_") {
		ShowMainMenu(bot, chatID)
	}
}
