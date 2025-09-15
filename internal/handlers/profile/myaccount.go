package profile

import (
	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/entry"
	"github.com/DavydAbbasov/trecker_bot/internal/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type ProfileModule struct {
	bot       interfaces.BotAPI
	entry     *entry.EntryModule
	repo      interfaces.Repo
	validator user.UserValidator
}

func New(bot interfaces.BotAPI, entry *entry.EntryModule, repo interfaces.Repo, validator user.UserValidator) *ProfileModule {
	return &ProfileModule{
		bot:       bot,
		entry:     entry,
		repo:      repo,
		validator: validator,
	}
}
func (d *ProfileModule) ShowProfileMock(ctx *context.MsgContext) {
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
	msg := tgbotapi.NewMessage(ctx.ChatID, text)

	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = BuildProfileKeyboard()
	_, err := d.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing profil")

	}
}

func (d *ProfileModule) ShowLanguageSelection(ctx *context.CallbackContext) {
	text := ("🌐 Выберите язык")

	keyboard := ShowLanguageSelectionKeyboard()

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	if _, err := d.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing language selection")
	}

}
func (d *ProfileModule) ShowEditProfileMenu(ctx *context.CallbackContext) {
	text := `
	👤🔁 *Обновления профиля*
	Выберите поле для изменения:
	`

	keyboard := ShowEditProfileMenuKeyboard()

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	if _, err := d.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing edit menu")
	}

}
