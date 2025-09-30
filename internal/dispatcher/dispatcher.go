// handles logic, routing, and handler delegation.
package dispatcher

import (
	stdctx "context"
	"strings"
	"time"

	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/domain"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/entry"
	inlinecommands "github.com/DavydAbbasov/trecker_bot/internal/handlers/handler_buttons/handler_inline"
	replycommands "github.com/DavydAbbasov/trecker_bot/internal/handlers/handler_buttons/handler_reply"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/learning"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/profile"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/subscription"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/track"
	helper "github.com/DavydAbbasov/trecker_bot/internal/lib/postgresql"
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
	repo            interfaces.UserRepository
	validator       *helper.Validator
	activities      domain.ActivityRepo
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager,
	activityStorage storage.ActivityStorage, learningStorage storage.LearningStorage,
	repo interfaces.UserRepository, validator *helper.Validator,
	activities domain.ActivityRepo) *Dispatcher {
	if bot == nil {
		log.Fatal().Msg("Dispatcher: nil bot interfaces.BotAPI")
	}

	d := &Dispatcher{
		bot:             bot,
		fsm:             fsm,
		activityStorage: activityStorage,
		learningStorage: learningStorage,
		repo:            repo,
		validator:       validator,
		activities:      activities,
	}

	d.entry = entry.New(bot, fsm, repo, validator)
	d.subscription = subscription.New(bot)
	d.profile = profile.New(bot, d.entry, repo, validator)
	d.track = track.New(bot, fsm, d.entry, activityStorage, activities)
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
		switch {
		case update.CallbackQuery != nil:
			d.RunCallback(update.CallbackQuery)
		case update.Message != nil && update.Message.IsCommand():
			d.entry.HandleCommand(update.Message)
		case update.Message != nil:
			d.handleMessage(update.Message)
		}
	}
}

func (d *Dispatcher) handleMessage(msg *tgbotapi.Message) {
	ctx := d.newMessageContext(msg)

	if d.handleUserState(ctx) {
		return
	}
	if d.replycommands.HandleReplyButtons(ctx) {
		return
	}

}

// Событие – обычное сообщение: пользователь написал текст или нажал reply-кнопку
func (d *Dispatcher) newMessageContext(msg *tgbotapi.Message) *context.MsgContext {

	ctxMsg := &context.MsgContext{
		Message:  msg,
		ChatID:   msg.Chat.ID,
		UserID:   int64(msg.From.ID), // телеграм-ID (почему берем под int64)
		UserName: strings.TrimSpace(msg.From.UserName),
		Text:     msg.Text,
	}
	reqCxt, cancel := stdctx.WithTimeout(stdctx.Background(), 3*time.Second)
	defer cancel()

	id, err := d.repo.EnsureIDByTelegram(reqCxt, ctxMsg.UserID, ctxMsg.UserName)
	if err != nil {
		log.Error().
			Err(err).
			Int64("tg_id", ctxMsg.UserID).
			Str("data", msg.ForwardSenderName).
			Msg("ensure user failed (MsgContext)")

		d.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Не удалось инициализировать пользователя. Попробуйте ещё раз."))
		return ctxMsg
	}
	ctxMsg.DBUserID = id //users.id
	return ctxMsg
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
