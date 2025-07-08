package dispatcher

/*
дирижёр, маршрутизатор
Запускает цикл получения апдейтов (Start)
Определяет, что за команда пришла, и кому её передать
Не должен сам обрабатывать команды — он только направляет.
*/
import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Она будет запускать цикл получения апдейтов и направлять их по маршрутам
func Start(bot interfaces.StoppableBot) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				HandleStart(bot, update.Message)
			}
		}
	}
}
