package dispatcher

import (
	"fmt"
	"strings"
	"time"

	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type TrackingUserState struct {
	State       string
	CurrentName string
}

var TrackingUserStates = map[int64]*TrackingUserState{}

type Activity struct {
	NameActivity string
	TimeEntry    []TimeEntry
}

type TimeEntry struct {
	Timestamp time.Time
	Start     time.Time
	End       time.Time
	Duration  time.Duration
}

var ActivityCollections = map[int64][]Activity{}

func ShowTrackingMenu(bot interfaces.BotAPI, chatID int64) {
	text := `
📈 *Track*

📊 Текущая Активность: *Go*  
⏱  Сегодняшний трек: *4 ч 52 мин* 
🔥 Стрик: *104 дня*  
📅 Сегодня: *4 активности*
`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = buildTrackKeyboard()

	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing profil")

	}
}
func buildTrackKeyboard() tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Выбрать активность", "selection_activity"),
		tgbotapi.NewInlineKeyboardButtonData("➕ Создать активность", "create_activity"),
	)

	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📈 Мои отчёты", "summary"),
		tgbotapi.NewInlineKeyboardButtonData("🛑 Завершить", "exit"),
	)

	return tgbotapi.NewInlineKeyboardMarkup(row1, row2)

}
func ShowActivityList(bot interfaces.BotAPI, chatID int64, userID int64) {
	text := "📋 Выберите активность для отчёта:"

	// заглушка (mock), потом реальные активности из хранения
	activities := ActivityCollections[userID]
	// activities := []string{"🦫Go", "📘English", "🏋️‍♀️Workout"}
	/*
					   Объявляется двумерный срез кнопок,
					   использоваться для создания инлайн-клавиатуры
					   [][] нужно чтобы Telegram понял, как отображать кнопки — по строкам.

					для каждой активности:
				Мы создаём кнопку (btn := ...)
				Помещаем её в ряд: row := [btn]
				Добавляем этот ряд в общую клавиатуру: rows = append(rows, row)

				Вот и получается двумерный срез: [][]InlineKeyboardButton

		[][]InlineKeyboardButton = массив рядов
		[]InlineKeyboardButton = один ряд
		InlineKeyboardButton = одна кнопка

		каждый РЯД кнопок — это срез,
		и все ряды вместе — это двумерный срез.


	*/
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, activity := range activities {
		if activity.NameActivity == "" {
			log.Warn().Msg("Обнаружена подборка без названия, пропускаем")
			continue
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(activity.NameActivity,
			"activity_report_"+activity.NameActivity)
		//Мы создаём строку (ряд) с этой одной кнопкой
		//И добавляем этот ряд в rows, чтобы собрать всю клавиатуру
		//Пользователь ещё ничего не выбрал — мы готовим клавиатуру для выбора.
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	/*
	   Создаёт объект клавиатуры NewInlineKeyboardMarkup

	   создаёт финальный объект InlineKeyboardMarkup,
	   который мы можем передать в msg.ReplyMarkup.
	*/

	//"вынь все элементы из среза rows и передай их по одному в функцию".
	inlineMenu := tgbotapi.NewInlineKeyboardMarkup(rows...)

	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📊 Сегодня"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)
	replyMenu.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu

	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("error sending activity list")
	}

	// Отдельно отправляем inline
	msg2 := tgbotapi.NewMessage(chatID, "🎯 Активированы активности:")
	msg2.ReplyMarkup = inlineMenu
	bot.Send(msg2)

}
func ShowActivityReport(bot interfaces.BotAPI, chatID int64, userID int64, activityName string) {

	activities := ActivityCollections[userID]

	if len(activities) == 0 {
		msgError := tgbotapi.NewMessage(chatID, "empty")
		if _, err := bot.Send(msgError); err != nil {
			log.Error().Err(err).Msg("ошибка при отправке сообщения")
			return
		}
	}

	text := fmt.Sprintf(`

📌 *Отчёт по активности:* _%s_ 
	
🔥 Стрик: *104 дня* 
📅 Начало: *12 мая 2024*  
📈 Дней подряд: *31*  
⏱ Время сегодня: *2 ч 40 мин*

Выберите, что вы хотите сделать:`, activityName)

	// var rows [][]tgbotapi.InlineKeyboardButton

	// for _, activity := range activities {
	// 	if activity.NameActivity == "" {
	// 		log.Warn().Msg("Обнаружена activity без названия, пропускаем")
	// 		continue
	// 	}

	// 	if activity.NameActivity == activityName {

	// 	}
	// 	btn := tgbotapi.NewInlineKeyboardButtonData(activity.NameActivity, "activities_report_menu"+activity.NameActivity)
	// 	rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	// }

	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Период"),
			tgbotapi.NewKeyboardButton("📊 Неделя"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📤 Экспорт"),
			tgbotapi.NewKeyboardButton("📊 Сегодня"),
		),

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗑 Удалить"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)
	replyMenu.ResizeKeyboard = true

	// inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// msgInline := tgbotapi.NewMessage(chatID, "Активности для статистики")
	// msgInline.ReplyMarkup = inlineKeyboard
	// if _, err := bot.Send(msgInline); err != nil {
	// 	log.Error().Err(err).Msg("error in displaying the inline report mrnu")
	// }

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu
	if _, err := bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error in displaying the report menu")
	}

}
func ShowCalendar(bot interfaces.BotAPI, chatID int64, activity string) {
	text := fmt.Sprintf(`
📊 *Недельный отчёт по активности:* _%s_
📅 *Статистика за желаемый период:*

📈 Ср. кл. ч. в день: *2ч 32мин*
📑 Сегодняшняя дата: *23.09.2026*

🗂 Выберите дату начала периода:
`, activity)

	// Названия дней недели
	weekDays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	// Примерные даты (можно потом генерировать на основе текущей недели)
	dates := []string{"8.7", "9.7", "10.7", "11.7", "12.7", "13.7", "14.7"}
	// Фиктивные данные активности (в дальнейшем подставим реальные)
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

	// 1. Reply-кнопки (внизу)
	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)
	replyMenu.ResizeKeyboard = true

	msgReply := tgbotapi.NewMessage(chatID, "📅")
	msgReply.ReplyMarkup = replyMenu
	if _, err := bot.Send(msgReply); err != nil {
		log.Error().Err(err).Msg("error showing calendar reply")
	}

	// 2. Инлайн-календарь
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = inlineMenu

	if _, err := bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing calendar inlain")
	}
}
func AddActivity(bot interfaces.BotAPI, chatID int64) {
	text := `
📌 *Создание новой активности*

Активности нужны для трекинга того, чем вы занимаетесь. Примеры:  
- 🧠 Go  
- 📚 English  
- 🏋️ Workout

Введите *название вашей активности* 
`
	TrackingUserStates[chatID] = &TrackingUserState{
		State: "waiting_for_activity_name",
	}

	replyMenu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)

	replyMenu.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu
	if _, err := bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing create activity prompt")
	}
}
func GetActivityMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⏱️ Таймер:15"),
			tgbotapi.NewKeyboardButton("⏱️ Таймер:60"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛠 Создать тамер"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)

}
func ProcessAddActivity(bot interfaces.BotAPI, msg *tgbotapi.Message) {
	userID := msg.From.ID
	chatID := msg.Chat.ID
	input := strings.TrimSpace(msg.Text)

	if input == "ℹ️ Помощь" {
		bot.Send(tgbotapi.NewMessage(chatID, "времено не доступно"))
		return
	}

	if input == "" {
		delete(TrackingUserStates, userID) //Удаляем пользователя из карты состояний
		ShowMainMenu(bot, chatID)
		return
	}

	state := TrackingUserStates[userID]
	state.CurrentName = input
	state.State = "activity_created" //обновляешь состояние пользователя - метка нового шага в логике.

	text := fmt.Sprintf("Ваша активность:%s,создана", input)
	confirmMsg := tgbotapi.NewMessage(chatID, text)
	confirmMsg.ParseMode = "Markdown"

	repluMenu := GetActivityMenuKeyboard()
	repluMenu.ResizeKeyboard = true
	confirmMsg.ReplyMarkup = repluMenu
	if _, err := bot.Send(confirmMsg); err != nil {
		log.Error().Err(err).Msg("err showing add_activity")
	}

	ActivityCollections[userID] = append(ActivityCollections[userID], Activity{
		NameActivity: input,
		TimeEntry:    []TimeEntry{},
	})

	followupMsg := tgbotapi.NewMessage(chatID, "➕ Теперь вы можете добавить таймер для трекинга.")
	bot.Send(followupMsg)

}
func SelectionActivityPromt(bot interfaces.BotAPI, chatID int64, userID int64) {
	text := `
📂 *Выбрать активность*

📂 Текущие активности: *🦫Go*
📂 Архив активносте: *12*

*Выберите активность для трека:*
`
	activities := ActivityCollections[userID]

	if len(activities) == 0 {
		msg := tgbotapi.NewMessage(chatID, "нет активностей")
		if _, err := bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("ошибка при отправке сообщения")
			return
		}
	}

	// activities := []string{"🦫Go", "📘English", "🏋️‍♀️Workout"}

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

	replyMenu.ResizeKeyboard = true

	msgReply := tgbotapi.NewMessage(chatID, "🗂")
	msgReply.ReplyMarkup = replyMenu
	if _, err := bot.Send(msgReply); err != nil {
		log.Error().Err(err).Msg("error showing calendar reply")
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = inlineMenu
	if _, err := bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("error showing calendar reply")
	}
}
