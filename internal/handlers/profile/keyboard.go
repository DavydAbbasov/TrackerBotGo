package profile

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildProfileKeyboard() tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🌐 Язык", "edit_language_"),
		tgbotapi.NewInlineKeyboardButtonData("📍 Часовой пояс", "edit_timezone"),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📧 Контакт", "edit_contact"),
		tgbotapi.NewInlineKeyboardButtonData("🔁 Обновить", "refresh_profile"),
	)
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2)
}
func ShowEditProfileMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Имя", "no_action"),
			tgbotapi.NewInlineKeyboardButtonData(" - ", "edit_firstname"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Соц.сеть", "no_action"),
			tgbotapi.NewInlineKeyboardButtonData(" - ", "edit_social"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Часовой пояс", "no_action"),
			tgbotapi.NewInlineKeyboardButtonData(" - ", "edit_timezone"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("↩ Назад Home", "go_back_profile"),
		),
	)
	return keyboard
}
func ShowLanguageSelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇺🇸 English", "lang_en"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇩🇪 Deutsch", "lang_de"),
			tgbotapi.NewInlineKeyboardButtonData("🇺🇦 Українська", "lang_uk"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇹🇷 Türkçe", "lang_tur"),
			tgbotapi.NewInlineKeyboardButtonData("🇸🇦 العربية", "lang_arab"),
		),
	)
	return keyboard
}
