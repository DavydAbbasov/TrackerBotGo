package replycommands

import (
	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/entry"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/learning"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/profile"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/subscription"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/track"
	log "github.com/rs/zerolog/log"
)

type ReplyModule struct {
	bot          interfaces.BotAPI
	track        *track.TrackModule
	subscription *subscription.SubscriptuonModule
	entry        *entry.EntryModule
	profile      *profile.ProfileModule
	learning     *learning.LearningModule
}

func New(bot interfaces.BotAPI, track *track.TrackModule, subscription *subscription.SubscriptuonModule, entry *entry.EntryModule, profile *profile.ProfileModule, learning *learning.LearningModule) *ReplyModule {
	return &ReplyModule{
		bot:          bot,
		track:        track,
		subscription: subscription,
		entry:        entry,
		profile:      profile,
		learning:     learning,
	}
}
func (r *ReplyModule) HandleReplyButtons(ctx *context.MsgContext) bool {
	replyButtons := map[string]func(*context.MsgContext){
		"👤My account":    r.handleShowProfileMock,
		"📈Track":         r.handleShowTrackingMenu,
		"🧠Learning":      r.handleShowLearningMenu,
		"💳 Subscription": r.handleShowSubscriptionMenu,
		"↩ Назад Home":   r.handleShowMainMenu,
		"📅 Период":       r.handleShowCalendar,
	}
	if handler, ok := replyButtons[ctx.Text]; ok {
		handler(ctx) //

		return true
	}
	log.Warn().Msgf("Unknown reply button: %s", ctx.Text) //?
	return false
}

// replu button
func (d *ReplyModule) handleShowProfileMock(ctx *context.MsgContext) {
	d.profile.ShowProfileMock(ctx)
}
func (r *ReplyModule) handleShowTrackingMenu(ctx *context.MsgContext) {
	r.track.ShowTrackingMenu(ctx)
}

func (r *ReplyModule) handleShowSubscriptionMenu(ctx *context.MsgContext) {
	r.subscription.ShowSubscriptionMenu(ctx)
}
func (r *ReplyModule) handleShowMainMenu(ctx *context.MsgContext) {
	r.entry.ShowMainMenu(ctx)
}
func (r *ReplyModule) handleShowCalendar(ctx *context.MsgContext) {
	r.track.ShowCalendar(ctx)
}

func (d *ReplyModule) handleShowLearningMenu(ctx *context.MsgContext) {
	d.learning.ShowLearningMenu(ctx)
}
