package dispatcher

import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func ShowProfileMock(bot interfaces.BotAPI, chatID int64) {
	text := `
👤 *My account*

— 👤 Имя: *Давид*
— 🧠 Активность: *Трекается*
— 🔥 Streak: *5 дней*
— 🌐 Язык: *Русский*
— 📍 Часовой пояс: *Europe/Berlin*
— 🗃 Подписка : *12 month*
— 📧 Контакт: @alaamov

Настроить профиль можно ниже:
`
	msg := tgbotapi.NewMessage(chatID, text)
	//Эта строка включает режим разметки Markdown
	msg.ParseMode = "Markdown"
	//Эта строка прикрепляет клавиатуру к сообщению,
	msg.ReplyMarkup = buildProfileKeyboard()
	//это вызов API Telegram,
	// который отправляет сообщение msg пользователю.
	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing profil")

	}
}
func buildProfileKeyboard() tgbotapi.InlineKeyboardMarkup {
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
func (d *Dispatcher)ShowLanguageSelection(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "🌐 Выберите язык")
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru"),
		tgbotapi.NewInlineKeyboardButtonData("🇺🇸 English", "lang_en"),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🇩🇪 Deutsch", "lang_de"),
		tgbotapi.NewInlineKeyboardButtonData("🇺🇦 Українська", "lang_uk"),
	)
	row3 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🇹🇷 Türkçe", "lang_tur"),
		tgbotapi.NewInlineKeyboardButtonData("🇸🇦 العربية", "lang_arab"),
	)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)

	if _, err := d.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing language selection")
	}

}
func (d *Dispatcher) ShowEditProfileMenu(chatID int64) {
	text := `
	👤🔁 *Обновления профиля*
	Выберите поле для изменения:
	`

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

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := d.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing edit menu")
	}

}
