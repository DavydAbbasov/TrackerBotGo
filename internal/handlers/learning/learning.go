package learning

import (
	"fmt"
	"strings"

	"github.com/DavydAbbasov/trecker_bot/interfaces"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher/context"
	"github.com/DavydAbbasov/trecker_bot/internal/handlers/entry"
	"github.com/DavydAbbasov/trecker_bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

// type UserState struct {
// 	State        string
// 	CurrentColl  string
// 	PendingWorld string
// }
// type Collections struct {
// 	TextInput1 string
// 	TextInput2 string
// }
// type Collection struct {
// 	NameCollection string
// 	Collections    []Collections
// }

// var UserStates = map[int64]*UserState{}

// var userCollections = map[int64][]Collection{}
type LearningModule struct {
	bot             interfaces.BotAPI
	fsm             interfaces.FSMManager
	entry           *entry.EntryModule
	learningStorage storage.LearningStorage
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager, entry *entry.EntryModule, learningStorage storage.LearningStorage) *LearningModule {
	return &LearningModule{
		bot:             bot,
		fsm:             fsm,
		entry:           entry,
		learningStorage: learningStorage,
	}
}
func (l *LearningModule) ShowLearningMenu(ctx *context.MsgContext) {
	text := `
🧠 *Learning*

🌐 Язык: *Английский 🇬🇧* 
📊 Добавлено слов: *463*  
📘 На сегодня: *10*  
✅ Выучено: *296*  
🕐 Следующее слово: *через 25 мин*

`
	msg := tgbotapi.NewMessage(ctx.ChatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = buuildLerningKeyboard()

	_, err := l.bot.Send(msg)
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

func (l *LearningModule) AddCollection(ctx *context.CallbackContext) {

	l.fsm.Set(ctx.UserID, "waiting_for_collection_name")
	// UserStates[chatID] = &UserState{
	// 	State: "waiting_for_collection_name",
	// }

	replyMenu := tgbotapi.NewReplyKeyboard(

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
			tgbotapi.NewKeyboardButton("↩ Назад Home"),
		),
	)

	replyMenu.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(ctx.ChatID, "📝")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = replyMenu
	if _, err := l.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}
	msg1 := tgbotapi.NewMessage(ctx.ChatID, "✏️ Введите имя новой подборки:")

	if _, err := l.bot.Send(msg1); err != nil {
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
func (l *LearningModule) ProcessCollectionCreation(ctx *context.MsgContext) {

	input := strings.TrimSpace(ctx.Text)

	if input == "ℹ️ Помощь" {
		l.bot.Send(tgbotapi.NewMessage(ctx.ChatID, "помощи нет"))
		return
	}

	if input == "↩ Назад Home" {
		l.fsm.Reset(ctx.UserID)
		// delete(UserStates, ctx.UserID)
		l.entry.ShowMainMenu(ctx)
		return
	}

	if input == "" || len(input) < 2 {
		msg := tgbotapi.NewMessage(ctx.ChatID, "⚠️ Имя подборки не может быть пустым. Пожалуйста, введите название.")
		l.bot.Send(msg)
		return
	}
	// t.fsm.Set(ctx.UserID, "activity_created")
	// t.fsm.SetData(ctx.UserID, "activity_name", input)

	// l.fsm.Set(ctx.UserID, "collection_created")
	// l.fsm.SetData(ctx.UserID, "createdivity_name", input)
	l.fsm.Reset(ctx.UserID)

	// state := UserStates[ctx.UserID]
	// state.CurrentColl = input
	// state.State = "collection_created"

	confirmMsg := tgbotapi.NewMessage(ctx.ChatID, fmt.Sprintf("📚 Подборка *%s* сохранена!", input))
	confirmMsg.ParseMode = "Markdown"
	confirmMsg.ReplyMarkup = GetLearningMenuKeyboard()
	if _, err := l.bot.Send(confirmMsg); err != nil {
		log.Error().Err(err).Msg("err showing learning") 
	}

	l.learningStorage.AddCollection(ctx.UserID, storage.Collection{
		NameCollection: input,
		Collection:     []storage.WordPair{},
	})

	// userCollections[ctx.UserID] = append(userCollections[ctx.UserID], Collection{
	// 	NameCollection: input,
	// 	Collections:    []Collections{},
	// })

	followupMsg := tgbotapi.NewMessage(ctx.ChatID, "➕ Теперь вы можете добавить слова для изучения.")
	l.bot.Send(followupMsg)
}

func (l *LearningModule) SowUserCollections(ctx *context.CallbackContext) {

	collections := l.learningStorage.ListCollections(ctx.UserID)

	if len(collections) == 0 {
		msg := tgbotapi.NewMessage(ctx.ChatID, "❌ У вас пока нет подборок.")
		if _, err := l.bot.Send(msg); err != nil {
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

	msg := tgbotapi.NewMessage(ctx.ChatID, "📂 Ваши подборки:")
	msg.ReplyMarkup = keyboard

	if _, err := l.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("ошибка при отправке сообщения подборки")
		return
	}
}
