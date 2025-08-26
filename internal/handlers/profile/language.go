package profile

import (
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (d *ProfileModule) LanguageChange(ctx *context.CallbackContext) {
	langMap := map[string]string{
		"lang_en":   "Language set to English.",
		"lang_ru":   "Язык установлен: русский.",
		"lang_dch":  "Sprache eingestellt: Deutsch.",
		"lang_ukr":  "Мову встановлено: українська.",
		"lang_arab": "تم تعيين اللغة: العربية.",
		"lang_tur":  "Dil ayarlandı: Türkçe.",
	}

	replyText, ok := langMap[ctx.Data]
	if !ok {
		replyText = "Unknown language selected"
	}

	// Уведомление Telegram, что мы обработали callback
	d.bot.Send(tgbotapi.NewCallback(ctx.Callback.ID, ""))

	// Удаляем предыдущее сообщение с кнопками
	d.bot.Send(tgbotapi.NewDeleteMessage(ctx.ChatID, ctx.Callback.Message.MessageID))

	// Отправляем сообщение об успешной смене языка
	d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, replyText))

	// Преобразуем CallbackContext → MsgContext
	msgCtx := &context.MsgContext{
		ChatID: ctx.ChatID,
		UserID: ctx.UserID,
	}
	// Показываем главное меню
	d.entry.ShowMainMenu(msgCtx)
}
