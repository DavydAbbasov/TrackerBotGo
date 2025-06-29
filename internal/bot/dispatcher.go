package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Она будет запускать цикл получения апдейтов и направлять их по маршрутам
func StartDispatcher(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message  == nil {
			continue
		}
		
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handleStart(bot, update.Message)
			}
		}
	}
}
