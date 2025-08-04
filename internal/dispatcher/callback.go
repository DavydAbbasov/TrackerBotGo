package dispatcher

import (
	"fmt"

	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

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

	ctx := d.NewCallbackContext(callback)

	if d.inlinecommands.HandlePrefixRoute(ctx) {
		return
	}

	if d.inlinecommands.HandleExactRoute(ctx) {
		return
	}
}
func (d *Dispatcher) NewCallbackContext(callback *tgbotapi.CallbackQuery) *context.CallbackContext {
	return &context.CallbackContext{
		Callback: callback,
		ChatID:   callback.Message.Chat.ID,
		UserID:   callback.From.ID,
		Data:     callback.Data,
	}
}
