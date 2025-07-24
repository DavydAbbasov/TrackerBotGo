package interfaces

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI defines minimal contract for Telegram bot operations used in the app.
type BotAPI interface {
	// Send sends a message or other Chattable content to Telegram.
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)

	// GetUpdatesChan starts receiving updates from Telegram using polling.
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel

	// StopReceivingUpdates gracefully shuts down the bot's update loop.
	StopReceivingUpdates()
}
