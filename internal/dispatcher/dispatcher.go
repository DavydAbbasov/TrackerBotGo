// handles logic, routing, and handler delegation.
package dispatcher

import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type Dispatcher struct {
	bot interfaces.BotAPI
}
type MsgContext struct {
	Message *tgbotapi.Message
	ChatID  int64
	UserID  int64
	Text    string
}

func New(bot interfaces.BotAPI) *Dispatcher {
	if bot == nil {
		log.Warn().Msg("Dispatcher: nil bot interfaces.BotAPI")
	}

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
			d.RunCallback(update.CallbackQuery)
		case update.Message != nil:
			d.handleMessage(update.Message)
		}
	}

}
func (d *Dispatcher) newMessageContext(msg *tgbotapi.Message) *MsgContext {
	return &MsgContext{
		Message: msg,
		ChatID:  msg.Chat.ID,
		UserID:  msg.From.ID,
		Text:    msg.Text,
	}
}
func (d *Dispatcher) handleMessage(msg *tgbotapi.Message) {//почему я тут передаю tgbotapi.Message
	ctx := d.newMessageContext(msg)

	if msg.IsCommand() { //?
		d.handleCommand(msg)
		return
	}
	if d.handleUserState(ctx) {
		return
	}
	if d.handleReplyButtons(ctx) {
		return
	}

}
func (d *Dispatcher) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "/start":
		HandleStart(d.bot, msg)//
	}
}
func (d *Dispatcher) handleUserState(ctx *MsgContext) bool {

	if state, ok := UserStates[ctx.UserID]; ok && state.State == "waiting_for_collection_name" {
		log.Printf("User %d said: %s", ctx.UserID, ctx.Text)
		d.ProcessCollectionCreation(ctx)
		return true
	}
	if state, ok := TrackingUserStates[ctx.UserID]; ok && state.State == "waiting_for_activity_name" {
		d.ProcessAddActivity(ctx)
		return true
	}
	return false
}
func (d *Dispatcher) handleReplyButtons(ctx *MsgContext) bool {
	replyButtons := map[string]func(*MsgContext){
		"👤My account":   d.handleShowProfileMock,
		"📈Track":        d.handleShowTrackingMenu,
		"🧠Learning":     d.handleShowLearningMenu,
		"💳Subscription": d.handleShowSubscriptionMenu,
		"↩ Назад Home":  d.handleShowMainMenu,
		"📅 Период":      d.handleShowCalendar,
	}
	if handler, ok := replyButtons[ctx.Text]; ok {
		handler(ctx)

		return true
	}
	log.Warn().Msgf("Unknown reply button: %s", ctx.Text) //?
	return false
}

// replu button
func (d *Dispatcher) handleShowProfileMock(ctx *MsgContext) {
	d.ShowProfileMock(ctx.ChatID)
}
func (d *Dispatcher) handleShowTrackingMenu(ctx *MsgContext) {
	d.ShowTrackingMenu(ctx.ChatID)
}
func (d *Dispatcher) handleShowLearningMenu(ctx *MsgContext) {
	d.ShowLearningMenu(ctx.ChatID)
}
func (d *Dispatcher) handleShowSubscriptionMenu(ctx *MsgContext) {
	d.ShowSubscriptionMenu(ctx.ChatID)
}
func (d *Dispatcher) handleShowMainMenu(ctx *MsgContext) {
	d.ShowMainMenu(ctx.ChatID)
}
func (d *Dispatcher) handleShowCalendar(ctx *MsgContext) {
	d.ShowCalendar(ctx.ChatID)
}
