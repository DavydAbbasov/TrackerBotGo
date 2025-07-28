package handlers

import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func HandleStart(bot interfaces.BotAPI, msg *tgbotapi.Message) {
	//We get the chat ID to know who to send the reply to.
	chatID := msg.Chat.ID

	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru"),
		tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang_en"),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🇩🇪 Deutsch", "lang_de"),
		tgbotapi.NewInlineKeyboardButtonData("🇺🇦 Українська", "lang_uk"),
	)
	row3 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🇹🇷 Türkçe", "lang_tur"),
		tgbotapi.NewInlineKeyboardButtonData("🇸🇦 العربية", "lang_arab"),
	)

	//Sending a message with a language selection
	text := "Choose your language"
	msg1 := tgbotapi.NewMessage(chatID, text)
	// message.ReplyMarkup = keyboard
	msg1.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)

	_, err := bot.Send(msg1) //Sends any text, emoji, command, button
	if err != nil {
		log.Error().Err(err).Msg("Error when sending the /start")
	}

}

func (d *Dispatcher)ShowMainMenu( chatID int64) {

	buttonMyAccount := tgbotapi.NewKeyboardButton("👤My account")
	buttonTrack := tgbotapi.NewKeyboardButton("📈Track")
	buttonSupport := tgbotapi.NewKeyboardButton("🧠Learning")
	buttonSubscriptions := tgbotapi.NewKeyboardButton("💳Subscription")

	row1 := tgbotapi.NewKeyboardButtonRow(buttonMyAccount, buttonTrack)
	row2 := tgbotapi.NewKeyboardButtonRow(buttonSupport, buttonSubscriptions)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)

	//Уменьшает клавиатуру под контент (не будет занимать весь экран).
	keyboard.ResizeKeyboard = true
	//Означает, что клавиатура не исчезнет после одного нажатия.
	keyboard.OneTimeKeyboard = false

	msg := tgbotapi.NewMessage(chatID, "🏠")
	msg.ReplyMarkup = keyboard

	_, err := d.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing menu")
	}
}
