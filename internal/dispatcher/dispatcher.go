// handles logic, routing, and handler delegation.
package dispatcher

import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Dispatcher struct {
	bot interfaces.BotAPI
}

func New(bot interfaces.BotAPI) *Dispatcher {
	return &Dispatcher{
		bot: bot,
	}
}
func (d *Dispatcher) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := d.bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			d.RunCallback(update.CallbackQuery) //log
			ctx := d.newCallbackContext(update.CallbackQuery)
			d.handlePrefixRoute(ctx)
			d.handleExactRoute(ctx)
		case update.Message != nil:
			d.handleMessage(update.Message)
		}
	}
}

func (d *Dispatcher) handleMessage(msg *tgbotapi.Message) {
	userID := msg.From.ID

	if state, ok := UserStates[userID]; ok && state.State == "waiting_for_collection_name" {
		d.ProcessCollectionCreation(msg)
		return
	}
	if state, ok := TrackingUserStates[userID]; ok && state.State == "waiting_for_activity_name" {
		d.ProcessAddActivity(msg)
		return
	}
	if msg.IsCommand() {
		d.handleCommand(msg)
		return
	}

	switch msg.Text {
	case "👤My account":
		ShowProfileMock(d.bot, msg.Chat.ID)
	case "📈Track":
		ShowTrackingMenu(d.bot, msg.Chat.ID)
	case "🧠Learning":
		ShowLearningMenu(d.bot, msg.Chat.ID)
	case "💳Subscription":
		ShowSubscriptionMenu(d.bot, msg.Chat.ID)
	case "↩ Назад Home":
		d.ShowMainMenu(msg.Chat.ID)
	case "📅 Период":
		d.ShowCalendar(msg.Chat.ID)
	}
}

func (d *Dispatcher) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "/start":
		HandleStart(d.bot, msg)
	}
}
