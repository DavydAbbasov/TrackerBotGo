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
func Start(bot interfaces.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// 1. Инлайн-кнопки (callback)
		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery)
			continue
		}
		// 2. Сообщение отсутствует — пропускаем
		if update.Message == nil {
			continue
		}

		// 1. Проверяем, есть ли состояние
		userID := update.Message.From.ID

		if state, ok := UserStates[userID]; ok && state.State == "waiting_for_collection_name" {
			ProcessCollectionCreation(bot, update.Message)
			continue
		}
		if state, ok := TrackingUserStates[userID]; ok && state.State == "waiting_for_activity_name" {
			ProcessAddActivity(bot, update.Message)
			continue
		}

		// 3. Команды (начинаются с "/")
		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "start":
				HandleStart(bot, update.Message)

			}

			continue
		}
		// 4. Текстовые кнопки (обычные сообщения)
		switch update.Message.Text {
		case "👤My account":
			ShowProfileMock(bot, update.Message.Chat.ID)
		case "📈Track":
			ShowTrackingMenu(bot, update.Message.Chat.ID)
		case "🧠Learning":
			ShowLearningMenu(bot, update.Message.Chat.ID)
		case "💳Subscription":
			ShowSubscriptionMenu(bot, update.Message.Chat.ID)
		case "↩ Назад Home":
			ShowMainMenu(bot, update.Message.Chat.ID)
		case "📅 Период":
			ShowCalendar(bot, update.Message.Chat.ID, "🦫Go")
		}
	}
}
