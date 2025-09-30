package track

import (
	ctx2 "context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/domain"
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
	activities    domain.ActivityRepo
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager, entry *entry.EntryModule, activityStore storage.ActivityStorage, activities domain.ActivityRepo) *TrackModule {
	return &TrackModule{
		bot:           bot,
		fsm:           fsm,
		entry:         entry,
		activityStore: activityStore,
		activities:    activities,
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
func GetActivityMenuKeyboardAddAndArchive() tgbotapi.ReplyKeyboardMarkup {
	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📳 Activate"),
			tgbotapi.NewKeyboardButton("🛒 Archive"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🐘 Delete"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)
	replyMenu.ResizeKeyboard = true

	return replyMenu
}

func (t *TrackModule) ProcessAddActivity(ctx *context.MsgContext) {
	log.Debug().Str("user", strconv.FormatInt(ctx.UserID, 10)).Str("text", ctx.Text).Msg("ProcessAddActivity вызван")

	input := strings.TrimSpace(ctx.Text)

	if input == "ℹ️ Помощь" {
		t.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "времено не доступно"))
		return
	}

	if input == "" {
		t.fsm.Reset(ctx.UserID) //сбрасываем состояние

		t.entry.ShowMainMenu(ctx)
		return
	}

	log.Debug().
		Int64("tg_id", ctx.UserID).
		Int64("db_user_id", ctx.DBUserID).
		Str("name", input).
		Msg("Create activity")

	ctx3 := ctx2.Background()
	act, err := t.activities.Create(ctx3, ctx.DBUserID, input, "")

	switch {
	case errors.Is(err, domain.ErrActivityExists):
		t.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "😒"))
		t.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "Такая активность уже есть."))
		return

	case err != nil:
		t.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "😥"))
		t.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "*Ошибка сохранения*: "+err.Error()))
		return

	default:
		t.fsm.Reset(ctx.UserID)
		t.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "👍"))

		time.Sleep(700 * time.Microsecond)

		text := fmt.Sprintf("*Активность*: *%s*, *создана😊*", act.Name)
		confirmMsg := tgbotapi.NewMessage(ctx.ChatID, text)
		confirmMsg.ParseMode = "Markdown"

		repluMenu := GetActivityMenuKeyboardAddAndArchive()
		confirmMsg.ReplyMarkup = repluMenu
		if _, err := t.bot.Send(confirmMsg); err != nil {
			log.Error().Err(err).Msg("err showing add_activity")
		}

		time.Sleep(500 * time.Microsecond)

		followupMsg := tgbotapi.NewMessage(ctx.ChatID, "👇")
		t.bot.Send(followupMsg)
	}

}

func (t *TrackModule) SelectionActivityPromt(ctx *context.CallbackContext) {
	text := `
📂 *Выбрать активность*

📂 Текущие активности: *🦫Go*
📂 Архив активносте: *12*

*Выберите активность для трека:*
`
	ctx3, stop := ctx2.WithTimeout(ctx2.Background(), 5*time.Second)
	defer stop()

	// Сразу собираем текст и inline-клавиатуру одной функцией
	text, markup, _, total, err := t.buildActivityUI(ctx3, ctx.DBUserID, true)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Не удалось загрузить активности. Попробуйте позже.")
		_, _ = t.bot.Send(msg)
		return
	}
	if total == 0 {
		// Нет активностей — покажем подсказку + reply-клавиатуру и выйдем
		msg := tgbotapi.NewMessage(ctx.ChatID, "Пока нет активностей. Нажмите «Добавить» ниже.")
		msg.ReplyMarkup = GetActivityMenuKeyboardAddAndArchive()
		_, _ = t.bot.Send(msg)
		return
	}

	// 1) Отправляем reply-клавиатуру (меню)
	msgReply := tgbotapi.NewMessage(ctx.ChatID, "🗂")
	msgReply.ReplyMarkup = GetActivityMenuKeyboardAddAndArchive()
	if _, err := t.bot.Send(msgReply); err != nil {
		log.Error().Err(err).Msg("send reply keyboard failed")
	}

	// 2) Отправляем сообщение с inline-кнопками (галочки уже учтены в buildActivityUI)
	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "HTML" // единый режим
	msg.ReplyMarkup = markup
	if _, err := t.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("send inline list failed")
	}
}

func (t *TrackModule) ActivityToggle(ctx *context.CallbackContext) {
	ctx3, stop := ctx2.WithTimeout(ctx2.Background(), 5*time.Second)
	defer stop()

	payload := strings.TrimPrefix(ctx.Data, "act_toggle_:")
	activityID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		t.bot.Request(tgbotapi.NewCallback(ctx.Callback.ID, "НЕКОРЕКТНЫЕ ДАННЫЕ"))
		return
	}

	if err := t.activities.ToggleSelectedActive(ctx3, ctx.DBUserID, activityID); err != nil {
		log.Error().Err(err).Msg("taggle selected failed")
		t.bot.Request(tgbotapi.NewCallback(ctx.Callback.ID, "не удалось изменить отобразить выбор"))
		return
	}

	// text, markup, _, _, err := t.buildActivityUI(ctx3, ctx.DBUserID, true)
	// if err != nil {
	// 	// редактируем текущее сообщение сообщением об ошибке и убираем клавиатуру
	// 	empty := tgbotapi.NewInlineKeyboardMarkup()
	// 	edit := tgbotapi.NewEditMessageTextAndMarkup(ctx.ChatID, ctx.Message.MessageID,
	// 		"⚠️ <b>Не удалось обновить список.</b>", empty)
	// 	edit.ParseMode = "HTML"
	// 	t.bot.Send(edit)
	// 	return
	// }
	// edit := tgbotapi.NewEditMessageTextAndMarkup(ctx.ChatID, ctx.Message.MessageID, text, markup)
	// edit.ParseMode = "HTML"
	// t.bot.Send(edit)

	acts, err := t.activities.ListActive(ctx3, ctx.DBUserID)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.ChatID, ctx.Message.MessageID, "⚠️ удалось обновить список.")
		t.bot.Send(edit)
		return
	}

	ids, err := t.activities.SelectedListActive(ctx3, ctx.DBUserID)
	if err != nil {
		log.Error().Err(err).Msg("failed select sctivities from actitity selected activities ")
		return
	}
	
	selected := make(map[int64]bool, len(ids))
	for _, id := range ids {
		selected[id] = true
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(ids)+1)
	for _, a := range acts {
		if strings.TrimSpace(a.Name) == "" {
			continue
		}
		check := "⚪"
		if selected[a.ID] {
			check = "🟢"
		}
		title := check + " " + a.Name
		if a.Emoji != "" {
			title = check + " " + a.Emoji + " " + a.Name
		}

		cb := fmt.Sprintf("act_toggle_:%d", a.ID)
		btn := tgbotapi.NewInlineKeyboardButtonData(title, cb)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("↩️ Назад", "back_to_main"),
	))

	text := fmt.Sprintf("📂 <b>Выбрать активность</b>\n\nВыбрано: %d из %d", len(ids), len(acts))

	// редактируем текущее сообщение (текст + клавиатуру)
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	edit := tgbotapi.NewEditMessageTextAndMarkup(ctx.ChatID, ctx.Message.MessageID, text, markup)
	edit.ParseMode = "HTML"

	if _, err := t.bot.Send(edit); err != nil {
		log.Error().Err(err).Msg("edit activities list failed")
	}

}

// общий билдер UI для экрана выбора активностей
func (t *TrackModule) buildActivityUI(ctx2 ctx2.Context, userID int64, addBack bool) (text string, markup tgbotapi.InlineKeyboardMarkup, selectedCount, total int, err error) {

	acts, err := t.activities.ListActive(ctx2, userID)
	if err != nil {
		return "", tgbotapi.NewInlineKeyboardMarkup(), 0, 0, err
	}

	ids, err := t.activities.SelectedListActive(ctx2, userID)
	if err != nil {
		// не падаем — просто без галочек (selected пустой)
		log.Error().Err(err).Msg("SelectedListActive failed")
	}

	// множество выбранных
	selected := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		selected[id] = struct{}{}
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(acts)+1)
	sel := 0
	for _, a := range acts {
		if strings.TrimSpace(a.Name) == "" {
			continue
		}

		check := "⚪"
		if _, ok := selected[a.ID]; ok {
			check = "🟢"
			sel++
		}
		title := check + " " + a.Name
		if a.Emoji != "" {
			title = check + " " + a.Emoji + " " + a.Name
		}

		cb := fmt.Sprintf("act_toggle_:%d", a.ID)
		btn := tgbotapi.NewInlineKeyboardButtonData(title, cb)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	if addBack {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("↩️ Назад", "back_to_main"),
		))
	}
	text = fmt.Sprintf("📂 <b>Выбрать активность</b>\n\nВыбрано: %d из %d", sel, len(acts))
	markup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return text, markup, sel, len(acts), nil
}
