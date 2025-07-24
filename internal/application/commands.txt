package application

import (
	"github.com/DavydAbbasov/trecker_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

const (
	commandStart   = "start"
	commandButtons = "buttons"
	action         = "Choose a action"
)

func (a *App) startCommand(message *tgbotapi.Message) error {
	// buttonTrack := tgbotapi.NewKeyboardButton("/track")
	// buttonHelp := tgbotapi.NewKeyboardButton("/help")
	// row := tgbotapi.NewKeyboardButtonRow(buttonTrack, buttonHelp)
	// msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(row)
	msg := tgbotapi.NewMessage(message.Chat.ID, action)
	_, err := a.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Error when sending the /start")
		return err
	}
	return nil
}

func (a *App) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return a.startCommand(message)

	case commandButtons:
		return a.getButtonsCommand(message)
	default:
		return a.unknownCommand(message)
	}

}

func (a *App) unknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, model.UnknownCommandString)
	_, err := a.bot.Send(msg)
	return err
}

func (a *App) getButtonsCommand(message *tgbotapi.Message) error {
	var row1 = []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("track", "track"),
		tgbotapi.NewInlineKeyboardButtonData("help", "help"),
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)
	msg := tgbotapi.NewMessage(message.Chat.ID, model.SettingCommandString)
	msg.ReplyMarkup = keyboard
	_, err := a.bot.Send(msg)

	return err
}
