// handles logic, routing, and handler delegation.
package dispatcher

import (
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/entry"
	inlinecommands "github.com/DavydAbbasov/trecker_bot/internal/handlers/handler_buttons/handler_inline"
	replycommands "github.com/DavydAbbasov/trecker_bot/internal/handlers/handler_buttons/handler_reply"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/learning"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/profile"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/subscription"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/track"
	"github.com/DavydAbbasov/trecker_bot/storage"

	"github.com/DavydAbbasov/trecker_bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type Dispatcher struct {
	bot             interfaces.BotAPI
	fsm             interfaces.FSMManager
	track           *track.TrackModule
	entry           *entry.EntryModule
	activityStorage storage.ActivityStorage
	subscription    *subscription.SubscriptuonModule
	inlinecommands  *inlinecommands.InlinneModule
	replycommands   *replycommands.ReplyModule
	profile         *profile.ProfileModule
	learning        *learning.LearningModule
	learningStorage storage.LearningStorage
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager, activityStorage storage.ActivityStorage, learningStorage storage.LearningStorage) *Dispatcher {
	if bot == nil {
		log.Fatal().Msg("Dispatcher: nil bot interfaces.BotAPI")
	}

	d := &Dispatcher{
		bot:             bot,
		fsm:             fsm,
		activityStorage: activityStorage,
		learningStorage: learningStorage,
	}

	d.entry = entry.New(bot, fsm)
	d.subscription = subscription.New(bot)
	d.profile = profile.New(bot, d.entry)
	d.track = track.New(bot, fsm, d.entry, activityStorage)
	d.learning = learning.New(bot, fsm, d.entry, learningStorage)
	d.inlinecommands = inlinecommands.New(bot, d.track, d.profile, d.learning)
	d.replycommands = replycommands.New(bot, d.track, d.subscription, d.entry, d.profile, d.learning)

	return d

}

func (d *Dispatcher) Run() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := d.bot.GetUpdatesChan(u)

	for update := range updates {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error().Any("panic", r).Msg("Recovered from panic in update loop")
				}
			}()

			switch {
			case update.CallbackQuery != nil:
				d.RunCallback(update.CallbackQuery)
			case update.Message != nil && update.Message.IsCommand():
				d.entry.HandleCommand(update.Message)
			case update.Message != nil:
				d.handleMessage(update.Message)
			}
		}()
	}

}

func (d *Dispatcher) handleMessage(msg *tgbotapi.Message) {
	ctx := d.newMessageContext(msg)

	if d.handleUserState(ctx) {
		return
	}
	if d.replycommands.HandleReplyButtons(ctx) { //
		return
	}

}
func (d *Dispatcher) newMessageContext(msg *tgbotapi.Message) *context.MsgContext {
	return &context.MsgContext{
		Message: msg,
		ChatID:  msg.Chat.ID,
		UserID:  msg.From.ID,
		Text:    msg.Text,
	}

}
func (d *Dispatcher) handleUserState(ctx *context.MsgContext) bool {
	state := d.fsm.Get(ctx.UserID)
	if state == nil {
		return false
	}

	switch state.State {
	case "waiting_for_collection_name":
		d.learning.ProcessCollectionCreation(ctx)
		return true
	case "waiting_for_activity_name":
		d.track.ProcessAddActivity(ctx)
		return true
	}
	return false

}

// Flush реализует интерфейс Flushable
func (d *Dispatcher) Flush() error { //?
	log.Info().Msg("dispatcher: flush called")
	return nil
}

// Close реализует интерфейс Flushable
func (d *Dispatcher) Close() error { //?
	log.Info().Msg("dispatcher: close called")
	return nil
}
