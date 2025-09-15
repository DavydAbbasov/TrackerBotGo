package entry

import (
	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type EntryModule struct {
	bot interfaces.BotAPI
	fsm interfaces.FSMManager
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager) *EntryModule {
	return &EntryModule{
		bot: bot,
		fsm: fsm,
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
	_, err := e.bot.Send(msg)
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
