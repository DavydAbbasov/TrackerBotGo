package profile

import (
	"strings"

	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func (d *ProfileModule) LanguageSwitch(ctx *context.CallbackContext) {
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
	// 1) из "lang_ru" → "ru"
	code := strings.TrimPrefix(ctx.Data, "lang_")

	// 2) валидируем только язык
	lang, err := d.validator.ValidateLanguage(code)
	if err != nil {
		// язык не из поддерживаемых — отвечаем и выходим
		_, _ = d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Unknown language selected"))
		return
	}
	// 3) апсерт по tg_id: записать язык
	u := &model.User{
		TgUserID: ctx.UserID, // или ctx.FromID — что у тебя в контексте
		Language: &lang,
	}
	if err := d.repo.CreateUserByTelegramID(ctx.Ctx, u); err != nil {
		log.Error().Err(err).Msg("update language failed")
		_, _ = d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Save error, try later"))
		return
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
