package track

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/entry"
	"github.com/DavydAbbasov/trecker_bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type TrackModule struct {
	bot           interfaces.BotAPI
	fsm           interfaces.FSMManager
	entry         *entry.EntryModule
	activityStore storage.ActivityStorage
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager, entry *entry.EntryModule, activityStore storage.ActivityStorage) *TrackModule {
	return &TrackModule{
		bot:           bot,
		fsm:           fsm,
		entry:         entry,
		activityStore: activityStore,
	}
}

func (t *TrackModule) ShowTrackingMenu(ctx *context.MsgContext) {
	data := ActivityReportData{
		CurrentActivity:    "Go",
		TodayTimeActivity:  "4 ч 52 мин",
		StreakActivity:     "104",
		TodayCountActivity: "4",
	}

	text := TrackingMenuText(data)

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = BuildTrackKeyboard()

	_, err := t.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing profil")

	}
}

func (t *TrackModule) ShowActivityList(ctx *context.CallbackContext) {
	activities := t.activityStore.List(ctx.UserID)

	inlineMenu := BuildActivityInlineKeyboard(activities)
	replyMenu := BuildActivityReplyKeyboard()

	msg := tgbotapi.NewMessage(ctx.ChatID, ActivityListTitle)

	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("error sending activity list")
	}

	msg2 := tgbotapi.NewMessage(ctx.ChatID, ActivityListConfirmed)
	msg2.ReplyMarkup = inlineMenu
	t.bot.Send(msg2)

}

func (t *TrackModule) ShowActivityReport(ctx *context.CallbackContext) {
	data := ActivityReportData{
		CurrentActivity:                 "Go",
		StreakCurrentActivity:           "4 ч 52 мин",
		StartDate:                       "104",
		ReportLabelConsecutive:          "4",
		ReportLabelTodayTimeAccumulated: "2",
	}

	text := ShowActivityReportText(data)

	activities := t.activityStore.List(ctx.UserID)
	if len(activities) == 0 {
		msgError := tgbotapi.NewMessage(ctx.ChatID, "empty")
		if _, err := t.bot.Send(msgError); err != nil {
			log.Error().Err(err).Msg("ошибка при отправке сообщения")
			return
		}
	}

	replyMenu := ShowActivityReportKeyboard(activities)
	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu
	if _, err := t.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error in displaying the report menu")
	}

}
func (t *TrackModule) ShowCalendar(ctx *context.MsgContext) {
	text := `
📊 *Недельный отчёт по активности:* 
📅 *Статистика за желаемый период:*

📈 Ср. кл. ч. в день: *2ч 32мин*
📑 Сегодняшняя дата: *23.09.2026*

🗂 Выберите дату начала периода:
`
	weekDays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

	dates := []string{"8.7", "9.7", "10.7", "11.7", "12.7", "13.7", "14.7"}

	data := []string{"2:10", "1:30", "0:00", "3", "2:45", "1", "0:15"}

	var rows [][]tgbotapi.InlineKeyboardButton

	var dayRow []tgbotapi.InlineKeyboardButton
	for _, day := range weekDays {
		dayRow = append(dayRow, tgbotapi.NewInlineKeyboardButtonData(day, "noop"))

	}
	rows = append(rows, dayRow)

	var dateRow []tgbotapi.InlineKeyboardButton
	for _, date := range dates {
		dateRow = append(dateRow, tgbotapi.NewInlineKeyboardButtonData(date, "noop"))
	}
	rows = append(rows, dateRow)

	var dataRow []tgbotapi.InlineKeyboardButton
	for _, entry := range data {
		dataRow = append(dataRow, tgbotapi.NewInlineKeyboardButtonData(entry, "noop"))
	}
	rows = append(rows, dataRow)

	navRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" ⏪ ", "prev_week"),
		tgbotapi.NewInlineKeyboardButtonData(" ⏩ ", "next_week"),
	)
	rows = append(rows, navRow)

	inlineMenu := tgbotapi.NewInlineKeyboardMarkup(rows...)

	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)
	replyMenu.ResizeKeyboard = true

	msgReply := tgbotapi.NewMessage(ctx.ChatID, "📅")
	msgReply.ReplyMarkup = replyMenu
	if _, err := t.bot.Send(msgReply); err != nil {
		log.Error().Err(err).Msg("error showing calendar reply")
	}

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = inlineMenu

	if _, err := t.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing calendar inlain")
	}
}
func (t *TrackModule) AddActivity(ctx *context.CallbackContext) {
	log.Debug().Str("user", strconv.FormatInt(ctx.UserID, 10)).Msg("AddActivity вызван")

	text := `
📌 *Создание новой активности*

Активности нужны для трекинга того, чем вы занимаетесь. Примеры:  
- 🧠 Go  
- 📚 English  
- 🏋️ Workout

Введите *название вашей активности* 
`
	// t.fsm.Reset(ctx.UserID)

	t.fsm.Set(ctx.UserID, "waiting_for_activity_name")

	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)

	replyMenu.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu
	if _, err := t.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing create activity prompt")
	}
}
func GetActivityMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⏱️ Таймер:15"),
			tgbotapi.NewKeyboardButton("⏱️ Таймер:60"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛠 Создать тамер"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)
	replyMenu.ResizeKeyboard = true

	return replyMenu
}

func (t *TrackModule) ProcessAddActivity(ctx *context.MsgContext) {
	log.Debug().Str("user", strconv.FormatInt(ctx.UserID, 10)).Str("text", ctx.Text).Msg("ProcessAddActivity вызван")

	// userID := msg.From.ID
	// chatID := msg.Chat.ID
	input := strings.TrimSpace(ctx.Text)

	if input == "ℹ️ Помощь" {
		t.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "времено не доступно"))
		return
	}

	if input == "" {
		t.fsm.Reset(ctx.UserID) //сбрасываем состояние
		// delete(TrackingUserStates, ctx.UserID) //Удаляем пользователя из карты состояний

		t.entry.ShowMainMenu(ctx)
		return
	}

	// Сохраняем данные во FSM
	t.activityStore.Add(ctx.UserID, storage.Activity{
		NameActivity: input,
		TimeEntry:    []storage.TimeEntry{},
	})

	// t.fsm.Set(ctx.UserID, "create_activity")
	// t.fsm.SetData(ctx.UserID, "activity_name", input)
	t.fsm.Reset(ctx.UserID)
	log.Debug().Str("user", fmt.Sprint(ctx.UserID)).Msg("FSM очищен после создания активности")
	// state := TrackingUserStates[ctx.UserID]
	// state.CurrentName = input
	// state.State = "activity_created"

	// Создаём активность
	text := fmt.Sprintf("Ваша активность:*%s*,создана", input)
	confirmMsg := tgbotapi.NewMessage(ctx.ChatID, text)
	confirmMsg.ParseMode = "Markdown"

	repluMenu := GetActivityMenuKeyboard()
	confirmMsg.ReplyMarkup = repluMenu
	if _, err := t.bot.Send(confirmMsg); err != nil {
		log.Error().Err(err).Msg("err showing add_activity")
	}

	// ActivityCollections[ctx.UserID] = append(ActivityCollections[ctx.UserID], Activity{
	// 	NameActivity: input,
	// 	TimeEntry:    []TimeEntry{},
	// })

	followupMsg := tgbotapi.NewMessage(ctx.ChatID, "➕ Теперь вы можете добавить таймер для трекинга.")
	t.bot.Send(followupMsg)

}
func (t *TrackModule) SelectionActivityPromt(ctx *context.CallbackContext) {
	text := `
📂 *Выбрать активность*

📂 Текущие активности: *🦫Go*
📂 Архив активносте: *12*

*Выберите активность для трека:*
`
	activities := t.activityStore.List(ctx.UserID)

	if len(activities) == 0 {
		msg := tgbotapi.NewMessage(ctx.ChatID, "нет активностей")
		if _, err := t.bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("ошибка при отправке сообщения")
			return
		}
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, activity := range activities {
		if activity.NameActivity == "" {
			log.Warn().Msg("Обнаружена подборка без названия, пропускаем")
			continue
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(activity.NameActivity, "activity_selection_"+activity.NameActivity)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	inlineMenu := tgbotapi.NewInlineKeyboardMarkup(rows...)

	replyMenu := GetActivityMenuKeyboard()

	msgReply := tgbotapi.NewMessage(ctx.ChatID, "🗂")
	msgReply.ReplyMarkup = replyMenu
	if _, err := t.bot.Send(msgReply); err != nil {
		log.Error().Err(err).Msg("error showing calendar reply")
	}

	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = inlineMenu
	if _, err := t.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing calendar reply")
	}
}
