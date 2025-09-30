package track

import (
	"github.com/DavydAbbasov/trecker_bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func BuildTrackKeyboard() tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ButtonSelectActivity, "selection_activity"),
		tgbotapi.NewInlineKeyboardButtonData(ButtonCreateActivity, "create_activity"),
	)

	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ButtonViewReports, "summary_activity"),
		tgbotapi.NewInlineKeyboardButtonData(ButtonViewArchive, "archive_activity"),
	)

	return tgbotapi.NewInlineKeyboardMarkup(row1, row2)
}

func BuildActivityInlineKeyboard(activities []storage.Activity) tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, activity := range activities {
		if activity.NameActivity == "" {
			log.Warn().Msg("Обнаружена подборка без названия, пропускаем")
			continue
		}

		btn := tgbotapi.NewInlineKeyboardButtonData(activity.NameActivity,
			"activity_report_"+activity.NameActivity)

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func BuildActivityReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {

	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ButtonReplyToday),
			tgbotapi.NewKeyboardButton(ButtonReplyBack),
		),
	)
	replyMenu.ResizeKeyboard = true
	return replyMenu
}
func ShowActivityReportKeyboard(activities []storage.Activity) tgbotapi.ReplyKeyboardMarkup {

	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Период"),
			tgbotapi.NewKeyboardButton("📊 Неделя"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📤 Экспорт"),
			tgbotapi.NewKeyboardButton(ButtonReplyToday),
		),

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗑 Удалить"),
			tgbotapi.NewKeyboardButton(ButtonReplyBack),
		),
	)
	replyMenu.ResizeKeyboard = true
	return replyMenu
}
