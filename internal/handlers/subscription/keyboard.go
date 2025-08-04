package subscription

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func BuildSubscriptionKeyboardKeyboard() tgbotapi.InlineKeyboardMarkup {
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
