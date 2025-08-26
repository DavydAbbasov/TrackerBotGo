package subscription

import "fmt"

const (
	SubscriptionMessage = "Для оформления подписки перейдите в раздел: 🗓 Тарифные планы"
)

func ShowSubscriptionMenuText(data SubscriptionReportData) string {
	return fmt.Sprintf(`
%s

%s *%s*
%s *%s*

%s
	`,

		MenuTitleMainSubscription,
		MenuLabelActivePlan, data.ActivePlanFree,
		MenuLabelDaysTheEnd, data.DaysEnd,

		SubscriptionMessage,
	)
}
