package context

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MsgContext struct {
	Message *tgbotapi.Message
	ChatID  int64
	UserID  int64
	Text    string
}
type CallbackContext struct {
	Ctx          context.Context
	Callback     *tgbotapi.CallbackQuery
	ChatID       int64
	UserID       int64
	Data         string
	ActivityName string
}
