package dispatcher

import (
	"fmt"
	"strings"

	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

func ShowLearningMenu(bot interfaces.BotAPI, chatID int64) {
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

	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}
}

// "контекст взаимодействия" с пользователем.
// моделируешь этапы взаимодействия с пользователем и для каждого этапа сохраняешь нужную информацию.
type UserState struct {
	State string /* <--- ключ для управления поведением
		ключевое поле, которое хранит текущее "состояние" пользователя.
		Узнать, что сейчас делает пользователь

	| Значение `State`                | Что это означает                           |
	| ------------------------------- | ------------------------------------------ |
	| `"waiting_for_collection_name"` | пользователь должен ввести имя подборки    |
	| `"collection_created"`          | подборка создана, ждём действий            |
	| `"waiting_for_word"`            | пользователь должен ввести слово           |
	| `"waiting_for_translation"`     | пользователь должен ввести перевод к слову |
	*/
	CurrentColl  string // CurrentТекущийCollсбор - имя подборки, в которую добавляем слова. Если пользователь создаёт подборку, нужно знать её имя
	PendingWorld string // PendingОжидающее слово -  временно запоминаем слово, ждём перевода. Когда пользователь вводит слово, а потом перевод — надо запомнить слово
}
type Collections struct {
	TextInput1 string //первый ввод
	TextInput2 string //второй ввод
}
type Collection struct { //[]Collection → список всех коллекций пользователя
	NameCollection string
	Collections    []Collections //[]Collections → стркуктура колекции
}

var UserStates = map[int64]*UserState{} // используется при создании подборки, когда мы отслеживаем состояние

var userCollections = map[int64][]Collection{} //Это глобальная переменная — мапа (map), в которой:
// int64 — ключ — это userID, уникальный ID каждого пользователя Telegram.
// []Collection — значение — список подборок (collections), принадлежащих этому пользователю.

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

func AddCollection(bot interfaces.BotAPI, chatID int64) {

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
	if _, err := bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}
	msg1 := tgbotapi.NewMessage(chatID, "✏️ Введите имя новой подборки:")

	if _, err := bot.Send(msg1); err != nil {
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
func ProcessCollectionCreation(bot interfaces.BotAPI, msg *tgbotapi.Message) {
	//начало обработки кнопки создать подборку
	userID := msg.From.ID
	chatID := msg.Chat.ID
	input := strings.TrimSpace(msg.Text)

	// фильтруем кнопки
	if input == "ℹ️ Помощь" {
		bot.Send(tgbotapi.NewMessage(chatID, "помощи нет"))
		return
	}

	if input == "↩ Назад Home" {
		delete(UserStates, userID) //после возврата в главное меню бот забыл текущее состояние
		// и не воспринимал следующую фразу как имя подборки.
		ShowMainMenu(bot, chatID)
		return
	}

	if input == "" || len(input) < 2 {
		msg := tgbotapi.NewMessage(chatID, "⚠️ Имя подборки не может быть пустым. Пожалуйста, введите название.")
		bot.Send(msg)
		return
	}

	// 2. Сохраняем введённое имя //Обновляем userState
	state := UserStates[userID]
	// 3. Обновляем состояние
	state.CurrentColl = input
	state.State = "collection_created"

	// 4. Формируем сообщение с подтверждением
	confirmMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📚 Подборка *%s* сохранена!", input))
	confirmMsg.ParseMode = "Markdown"
	confirmMsg.ReplyMarkup = GetLearningMenuKeyboard()
	if _, err := bot.Send(confirmMsg); err != nil {
		log.Error().Err(err).Msg("err showing learning")
	}

	// Здесь добавляем подборку в список пользователя
	userCollections[userID] = append(userCollections[userID], Collection{
		NameCollection: input,
		Collections:    []Collections{},
	})

	followupMsg := tgbotapi.NewMessage(chatID, "➕ Теперь вы можете добавить слова для изучения.")
	bot.Send(followupMsg)
}

func SowUserCollections(bot interfaces.BotAPI, chatID int64, userID int64) {

	collections := userCollections[userID]

	if len(collections) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ У вас пока нет подборок.")
		if _, err := bot.Send(msg); err != nil {
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

	if _, err := bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("ошибка при отправке сообщения подборки")
		return
	}
}
