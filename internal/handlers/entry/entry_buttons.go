package entry

import (
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func (e *EntryModule) HandleLanguaheStart(msg *tgbotapi.Message) {
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

	text := "Choose your language"
	msg1 := tgbotapi.NewMessage(chatID, text)

	msg1.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)
	_, err := e.bot.Send(msg1)
	if err != nil {
		log.Error().Err(err).Msg("Error when sending the /start")
	}

}

func (e *EntryModule) ShowMainMenu(ctx *context.MsgContext) {

	buttonMyAccount := tgbotapi.NewKeyboardButton("👤My account")
	buttonTrack := tgbotapi.NewKeyboardButton("📈Track")
	buttonSupport := tgbotapi.NewKeyboardButton("🧠Learning")
	buttonSubscriptions := tgbotapi.NewKeyboardButton("💳 Subscription")

	row1 := tgbotapi.NewKeyboardButtonRow(buttonMyAccount, buttonTrack)
	row2 := tgbotapi.NewKeyboardButtonRow(buttonSupport, buttonSubscriptions)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)

	keyboard.ResizeKeyboard = true
	//Означает, что клавиатура не исчезнет после одного нажатия.
	keyboard.OneTimeKeyboard = false

	msg := tgbotapi.NewMessage(ctx.ChatID, "🏠")
	msg.ReplyMarkup = keyboard

	_, err := e.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing menu")
	}
}
