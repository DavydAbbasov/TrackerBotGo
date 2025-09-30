package context

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MsgContext struct {
	Message  *tgbotapi.Message
	ChatID   int64
	UserID   int64
	Text     string
	DBUserID int64 // ВНУТРЕННИЙ users.id
	UserName string
}
type CallbackContext struct {
	Message      *tgbotapi.Message
	Ctx          *context.Context
	Callback     *tgbotapi.CallbackQuery
	ChatID       int64
	UserID       int64
	Data         string
	ActivityName string
	DBUserID     int64 // ВНУТРЕННИЙ users.id
	UserName     string
}
