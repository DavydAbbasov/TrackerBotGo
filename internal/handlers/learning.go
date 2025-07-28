package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)



type UserState struct {
	State        string
	CurrentColl  string
	PendingWorld string
}
type Collections struct {
	TextInput1 string
	TextInput2 string
}
type Collection struct {
	NameCollection string
	Collections    []Collections
}

var UserStates = map[int64]*UserState{}

var userCollections = map[int64][]Collection{}

func (d *Dispatcher) ShowLearningMenu(chatID int64) {
	text := `
🧠 *Learning*

🌐 Язык: *Английский 🇬🇧* 
📊 Добавлено слов: *463*  
📘 На сегодня: *10*  
✅ Выучено: *296*  
🕐 Следующее слово: *через 25 мин*

`
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = buuildLerningKeyboard()

	_, err := d.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}
}
func buuildLerningKeyboard() tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Создать подборку", "add_collection"),
		tgbotapi.NewInlineKeyboardButtonData("🎲 Случайная подборка", "random_words"),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔁 Архив подборок", "switch_learning_actv"),
		tgbotapi.NewInlineKeyboardButtonData("📈 Статистика", "summary_learning"),
	)
	row3 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🗂 База слов", "base_words"))
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)
}

func (d *Dispatcher) AddCollection(chatID int64) {

	UserStates[chatID] = &UserState{
		State: "waiting_for_collection_name",
	}

	replyMenu := tgbotapi.NewReplyKeyboard(

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)

	replyMenu.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, "📝")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu
	if _, err := d.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}
	msg1 := tgbotapi.NewMessage(chatID, "✏️ Введите имя новой подборки:")

	if _, err := d.bot.Send(msg1); err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}

}
func GetLearningMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Добавить слово"),
			tgbotapi.NewKeyboardButton("✅ Завершить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)
}
func (d *Dispatcher) ProcessCollectionCreation(ctx *MsgContext) {

	input := strings.TrimSpace(ctx.Text)

	if input == "ℹ️ Помощь" {
		d.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "помощи нет"))
		return
	}

	if input == "↩ Назад Home" {
		delete(UserStates, ctx.UserID)
		d.ShowMainMenu(ctx.ChatID)
		return
	}

	if input == "" || len(input) < 2 {
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Имя подборки не может быть пустым. Пожалуйста, введите название.")
		d.bot.Send(msg)
		return
	}

	state := UserStates[ctx.UserID]
	state.CurrentColl = input
	state.State = "collection_created"

	confirmMsg := tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("📚 Подборка *%s* сохранена!", input))
	confirmMsg.ParseMode = "Markdown"
	confirmMsg.ReplyMarkup = GetLearningMenuKeyboard()
	if _, err := d.bot.Send(confirmMsg); err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}

	userCollections[ctx.UserID] = append(userCollections[ctx.UserID], Collection{
		NameCollection: input,
		Collections:    []Collections{},
	})

	followupMsg := tgbotapi.NewMessage(ctx.ChatID, "➕ Теперь вы можете добавить слова для изучения.")
	d.bot.Send(followupMsg)
}

func (d *Dispatcher) SowUserCollections(chatID int64, userID int64) {

	collections := userCollections[userID]

	if len(collections) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ У вас пока нет подборок.")
		if _, err := d.bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("ошибка при отправке сообщения")
			return
		}
	}
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, coll := range collections {
		if coll.NameCollection == "" {
			log.Warn().Msg("Обнаружена подборка без названия, пропускаем")
			continue
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(coll.NameCollection, "view_collection_"+coll.NameCollection)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, "📂 Ваши подборки:")
	msg.ReplyMarkup = keyboard

	if _, err := d.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("ошибка при отправке сообщения подборки")
		return
	}
}
