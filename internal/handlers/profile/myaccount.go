package profile

import (
	ctx2 "context"
	"time"

	"github.com/DavydAbbasov/trecker_bot/interfaces"

	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/entry"
	helper "github.com/DavydAbbasov/trecker_bot/internal/lib/postgresql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type ProfileModule struct {
	bot       interfaces.BotAPI
	entry     *entry.EntryModule
	repo      interfaces.UserRepository
	validator *helper.Validator
}

func New(bot interfaces.BotAPI, entry *entry.EntryModule, repo interfaces.UserRepository, validator *helper.Validator) *ProfileModule {
	return &ProfileModule{
		bot:       bot,
		entry:     entry,
		repo:      repo,
		validator: validator,
	}
}
func (d *ProfileModule) ShowProfileMock(ctx *context.MsgContext) {
	ctx2, cancel := ctx2.WithTimeout(ctx2.Background(), 2*time.Second)
	defer cancel()

	u, err := d.repo.GetUserByTelegramID(ctx2, ctx.UserID)
	if err != nil {
		log.Error().Err(err).Msg("profile: get user failed")
	}

	view := BuildProfileView(u, ctx.Message.From)
	txt := view.Markdown()

	msg := tgbotapi.NewMessage(ctx.ChatID, txt)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = BuildProfileKeyboard()
	_, err = d.bot.Send(msg)
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
