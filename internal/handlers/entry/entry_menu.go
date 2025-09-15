package entry

import (
	"context"

	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type EntryModule struct {
	bot        interfaces.BotAPI
	fsm        interfaces.FSMManager
	repository interfaces.UserRepository
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager, repository interfaces.UserRepository) *EntryModule {
	return &EntryModule{
		bot:        bot,
		fsm:        fsm,
		repository: repository,
	}
}
func (e *EntryModule) HandleCommand(message *tgbotapi.Message) {
	switch message.Command() { //?
	case model.CommandStart:
		e.handleStart(message)
		e.HandleLanguageStart(message)
	default:
		e.handleUnknown(message)
	}
}

func (e *EntryModule) handleStart(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, model.WELCOME)

	ctx := context.Background()

	// Save user into DB
	user := model.User{
		TgUserID: message.Chat.ID,
		// и т.д

	}

	err := e.repository.CreateUserByTelegramID(ctx, &user)
	if err != nil {
		log.Error().Err(err).Msgf("failed create new user: %v", user.ID)
	}

	_, err = e.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Start not initiated")
	}
}

func (e *EntryModule) handleUnknown(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, model.TextUnknownCmd)
	_, err := e.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("unknown command error")
	}
}
