package subscription

import (
	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type SubscriptuonModule struct {
	bot interfaces.BotAPI
}

func New(bot interfaces.BotAPI) *SubscriptuonModule {
	return &SubscriptuonModule{
		bot: bot,
	}
}
func (s *SubscriptuonModule) ShowSubscriptionMenu(ctx *context.MsgContext) {
	data := SubscriptionReportData{
		ActivePlanFree: "Free",
		DaysEnd:        "23",
	}

	text := ShowSubscriptionMenuText(data)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = BuildSubscriptionKeyboardKeyboard()

	_, err := s.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing Subscription")

	}
}
