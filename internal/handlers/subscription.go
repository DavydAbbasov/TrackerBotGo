package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func (d *Dispatcher) ShowSubscriptionMenu(chatID int64) {
	text := `
💳*Subscription*

🔁 Активный план: *Free*  
📅 Дней до окончания: *35*

Для оформления подписки перейдите в раздел: 🗓 Тарифные планы

`
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = buildSubscriptionKeyboard()

	_, err := d.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing Subscription")

	}
}
func buildSubscriptionKeyboard() tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🗓 Тарифные планы", "tariff_plans"),
		tgbotapi.NewInlineKeyboardButtonData("🎁 Free", "free_version"),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🛫 Поддержка", "support"),
		tgbotapi.NewInlineKeyboardButtonData("💳 Изменить оплату", "change_payment"),
	)
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2)
}
