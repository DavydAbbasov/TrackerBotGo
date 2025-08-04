package inlinecommands

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/learning"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/profile"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/track"
)

type InlinneModule struct {
	bot      interfaces.BotAPI
	track    *track.TrackModule
	profile  *profile.ProfileModule
	learning *learning.LearningModule
}

func New(bot interfaces.BotAPI, track *track.TrackModule, profile *profile.ProfileModule, learning *learning.LearningModule) *InlinneModule {
	return &InlinneModule{
		bot:      bot,
		track:    track,
		profile:  profile,
		learning: learning,
	}
}
func (i *InlinneModule) HandlePrefixRoute(ctx *context.CallbackContext) bool {
	// route table
	prefixHandlers := map[string]func(*context.CallbackContext){ //?
		"activity_report_": i.handleActivityReport,
		"lang_":            i.handleLanguageChange,
		"edit_language_":   i.handleShowLanguageSelection,
	}
	// dynamic route
	for prefix, handler := range prefixHandlers {
		if strings.HasPrefix(ctx.Data, prefix) {
			ctx.ActivityName = strings.TrimPrefix(ctx.Data, prefix) //?
			handler(ctx)                                            //?
			i.callbackResponse(ctx)
			return true
		}

	}
	i.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Unknown prefixHandlers"))

	return false
}

func (i *InlinneModule) HandleExactRoute(ctx *context.CallbackContext) bool {
	log.Debug().Str("data", ctx.Data).Msg("Получен callback")

	exactHandlers := map[string]func(ctx *context.CallbackContext){ //?
		"refresh_profile":      i.handleShowEditProfileMenu,
		"summary_activity":     i.handleShowActivityList,
		"create_activity":      i.handleAddActivity,
		"selection_activity":   i.handleSelectionActivityPromt,
		"add_collection":       i.handleAddCollection,
		"switch_learning_actv": i.handleSowUserCollections,
	}

	if handler, ok := exactHandlers[ctx.Data]; ok { //?
		handler(ctx) //?
		i.callbackResponse(ctx)
		return true
	}
	i.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Unknown exactHandlers"))

	return false
}

// prefixHandlers
func (i *InlinneModule) handleActivityReport(ctx *context.CallbackContext) {
	i.track.ShowActivityReport(ctx)
}

func (i *InlinneModule) handleShowLanguageSelection(ctx *context.CallbackContext) {
	i.profile.ShowLanguageSelection(ctx)
}

func (i *InlinneModule) handleLanguageChange(ctx *context.CallbackContext) {
	i.profile.LanguageChange(ctx)
}

// exactHandlers
func (i *InlinneModule) handleShowEditProfileMenu(ctx *context.CallbackContext) {
	i.profile.ShowEditProfileMenu(ctx)
}
func (i *InlinneModule) handleAddActivity(ctx *context.CallbackContext) {
	i.track.AddActivity(ctx)
}

func (i *InlinneModule) handleSelectionActivityPromt(ctx *context.CallbackContext) {
	i.track.SelectionActivityPromt(ctx)
}
func (i *InlinneModule) handleShowActivityList(ctx *context.CallbackContext) {
	i.track.ShowActivityList(ctx)
}

func (i *InlinneModule) handleSowUserCollections(ctx *context.CallbackContext) {
	i.learning.SowUserCollections(ctx)
}

func (i *InlinneModule) handleAddCollection(ctx *context.CallbackContext) {
	i.learning.AddCollection(ctx)
}

// Confirming callback receipt
func (i *InlinneModule) callbackResponse(ctx *context.CallbackContext) {
	i.bot.Send(tgbotapi.NewCallback(ctx.Callback.ID, ""))
}

/*
если префиксы будут пересекаться или усложняться
(например, lang_en_menu)- перейти
на дерево префиксов или strings.SplitN.
*/
