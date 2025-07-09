package dispatcher

import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

/*
handlers - исполнители (обработчики команд)
Реализуют, что делать по команде (/start, /help, /track)
Отвечают за бизнес-логику команды
(например, "Привет, я бот" или "добавить активность")
*/
func HandleStart(bot interfaces.StoppableBot, msg *tgbotapi.Message) {
	//We get the chat ID to know who to send the reply to.
	chatID := msg.Chat.ID
	//creating button text
	buttonTrack := tgbotapi.NewKeyboardButton("/track")//a visual button in the keyboard.
	buttonHelp := tgbotapi.NewKeyboardButton("/help")
	//merge two button in row for display ux
	row := tgbotapi.NewKeyboardButtonRow(buttonTrack, buttonHelp)
	//creating object keyboard
	keyboard := tgbotapi.NewReplyKeyboard(row)
	// creating message text
	text := "Choose a action"//gets the user 
	message := tgbotapi.NewMessage(chatID, text)//chatID-to whom to send,text-What to say
	//attach the keyboard
	message.ReplyMarkup = keyboard//"Show this keyboard along with the message."
	//send
	_, err := bot.Send(message)//Sends any text, emoji, command, button
	if err != nil {
		log.Error().Err(err).Msg("Error when sending the /start")
	}

}
