package handlers

import (
	"github.com/DavydAbbasov/trecker_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func (d *Dispatcher) HandleCommand(message *tgbotapi.Message) {
	switch message.Command() { //?
	case model.CommandStart:
		d.handleStart(message)
	default:
		d.handleUnknown(message)
	}
}
func (d *Dispatcher) handleStart(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, model.WELCOME)
	_, err := d.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Start not initiated")
	}
}
func (d *Dispatcher) handleUnknown(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, model.TextUnknownCmd)
	_, err := d.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("unknown command error")
	}
}
