package profile

import (
	context2 "context"
	"time"

	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (d *ProfileModule) LanguageSwitch(ctx *context.CallbackContext) {
	// Преобразуем CallbackContext → MsgContext
	msgCtx := &context.MsgContext{
		ChatID: ctx.ChatID,
		UserID: ctx.UserID,
	}

	CodeLang := map[string]string{
		"lang_ru":   "ru",
		"lang_en":   "en",
		"lang_de":   "de",
		"lang_uk":   "uk",
		"lang_arab": "ar",
	}

	replyCodeLang, ok := CodeLang[ctx.Data]
	if !ok {
		replyCodeLang = "Unknown language (mock) selected"
		d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, replyCodeLang))
		return
	}
	// валидация кодов через валидатор
	if _, err := d.validator.ValidateLanguage(replyCodeLang); err != nil {
		log.Error().Err(err).Msg("error validation language ")
		return
	}

	// сохранить в БД
	ctx2, cancel := context2.WithTimeout(context2.Background(), 5*time.Second)
	defer cancel()
	if _, err := d.repo.UpdateLanguage(ctx2, msgCtx.UserID, replyCodeLang); err != nil {
		log.Error().Err(err).Msg("error updare language and save to bd")
		return
	}

	langMap := map[string]string{
		"ru": "Язык установлен: русский!",
		"en": "Language set to English!",
		"de": "Sprache eingestellt: Deutsch!",
		"uk": "Мову встановлено: українська!",
		"ar": "تم تعيين اللغة: العربية!",
	}

	replyText, ok := langMap[replyCodeLang]

	// Уведомление Telegram, что мы обработали callback
	d.bot.Send(tgbotapi.NewCallback(ctx.Callback.ID, ""))

	// Удаляем предыдущее сообщение с кнопками
	d.bot.Send(tgbotapi.NewDeleteMessage(ctx.ChatID, ctx.Callback.Message.MessageID))

	// Отправляем сообщение об успешной смене языка
	d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, replyText))

	// Показываем главное меню
	d.entry.ShowMainMenu(msgCtx)
}
