package dispatcher

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

// uses pointers ?
type CallbackContext struct {
	Callback     *tgbotapi.CallbackQuery
	ChatID       int64
	UserID       int64
	Data         string
	ActivityName string
}

// pre-routing
func (d *Dispatcher) RunCallback(callback *tgbotapi.CallbackQuery) {
	if callback == nil || callback.Message == nil {
		log.Warn().Msg("CallbackQuery: nil callback tgbotapi.CallbackQuery")
		return

	} else {
		log.Debug().
			Str("user", fmt.Sprint(callback.From.ID)). //remake me
			Str("data", callback.Data).
			Msg("Callback context initialized")
	}

	ctx := d.newCallbackContext(callback)

	if d.handlePrefixRoute(ctx) {
		return
	}

	if d.handleExactRoute(ctx) {
		return
	}
}

// pre-routing
func (d *Dispatcher) newCallbackContext(callback *tgbotapi.CallbackQuery) *CallbackContext {
	return &CallbackContext{
		Callback: callback,
		ChatID:   callback.Message.Chat.ID,
		UserID:   callback.From.ID,
		Data:     callback.Data,
	}
}
func (d *Dispatcher) handlePrefixRoute(ctx *CallbackContext) bool {
	// route table
	prefixHandlers := map[string]func(*CallbackContext){ //?
		"calendar_":        d.handleCalendar,
		"activity_report_": d.handleActivityReport,
		"lang_":            d.handleLanguageChange,
		"edit_language_":   d.handleShowLanguageSelection,
	}
	// dynamic route
	for prefix, handler := range prefixHandlers {
		if strings.HasPrefix(ctx.Data, prefix) {
			ctx.ActivityName = strings.TrimPrefix(ctx.Data, prefix) //?
			handler(ctx)                                            //?
			d.callbackResponse(ctx)
			return true
		}

	}
	d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Unknown prefixHandlers"))

	return false
}

func (d *Dispatcher) handleExactRoute(ctx *CallbackContext) bool {

	exactHandlers := map[string]func(ctx *CallbackContext){ //?
		"refresh_profile":      d.handleShowEditProfileMenu,
		"summary_activity":     d.handleShowActivityList,
		"create_activity":      d.handleAddActivity,
		"selection_activity":   d.handleSelectionActivityPromt,
		"add_collection":       d.handleAddCollection,
		"switch_learning_actv": d.handleSowUserCollections,
	}

	if handler, ok := exactHandlers[ctx.Data]; ok { //?
		handler(ctx) //?
		d.callbackResponse(ctx)
		return true
	}
	d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Unknown exactHandlers"))

	return false
}

// prefixHandlers
func (d *Dispatcher) handleActivityReport(ctx *CallbackContext) {
	d.ShowActivityReport(ctx.ChatID, ctx.UserID, ctx.ActivityName)
}
func (d *Dispatcher) handleCalendar(ctx *CallbackContext) {
	d.ShowCalendar(ctx.ChatID)
}
func (d *Dispatcher) handleLanguageChange(ctx *CallbackContext) {
	d.LanguageChange(ctx.ChatID, ctx)
}
func (d *Dispatcher) handleShowLanguageSelection(ctx *CallbackContext) {
	d.ShowLanguageSelection(ctx.ChatID)
}

// exactHandlers
func (d *Dispatcher) handleShowEditProfileMenu(ctx *CallbackContext) {
	d.ShowEditProfileMenu(ctx.ChatID)
}
func (d *Dispatcher) handleAddActivity(ctx *CallbackContext) {
	d.AddActivity(ctx.ChatID)
}
func (d *Dispatcher) handleAddCollection(ctx *CallbackContext) {
	d.AddCollection(ctx.ChatID)
}
func (d *Dispatcher) handleSelectionActivityPromt(ctx *CallbackContext) {
	d.SelectionActivityPromt(ctx.ChatID, ctx.UserID)
}
func (d *Dispatcher) handleShowActivityList(ctx *CallbackContext) {
	d.ShowActivityList(ctx.ChatID, ctx.UserID)
}
func (d *Dispatcher) handleSowUserCollections(ctx *CallbackContext) {
	d.SowUserCollections(ctx.ChatID, ctx.UserID)
}

// Confirming callback receipt
func (d *Dispatcher) callbackResponse(ctx *CallbackContext) {
	d.bot.Send(tgbotapi.NewCallback(ctx.Callback.ID, ""))
}

/*
если префиксы будут пересекаться или усложняться
(например, lang_en_menu)- перейти
на дерево префиксов или strings.SplitN.
*/
