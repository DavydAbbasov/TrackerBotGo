package track

import "fmt"

const (
	ActivityListTitle     = "📋 Выберите активность для отчёта:"
	ActivityListConfirmed = "🎯 Активированы активности:"
)

func TrackingMenuText(data ActivityReportData) string {
	return fmt.Sprintf(`
%s
 %s *%s*
 %s *%s*
 %s *%s дня*
 %s *%s активности*
`,
		MenuTitleMainTrack,
		MenuLabelCurrent, data.CurrentActivity,
		MenuLabelTodayTime, data.TodayTimeActivity,
		MenuLabelStreak, data.StreakActivity,
		MenuLabelTodayCount, data.TodayCountActivity)
}

func ShowActivityReportText(data ActivityReportData) string {
	return fmt.Sprintf(`

%s *%s*
%s *%s дня*
%s *%s*
%s *%s*
%s *%s*

Выберите, что вы хотите сделать:`,
		ReportTitleActivity, data.CurrentActivity,

		MenuLabelStreak, data.StreakCurrentActivity,
		ReportLabelStartDate, data.StartDate,
		ReportLabelConsecutive, data.ReportLabelConsecutive,
		ReportLabelTodayTimeAccumulated, data.ReportLabelTodayTimeAccumulated,
	)
}
