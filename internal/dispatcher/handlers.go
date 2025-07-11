package dispatcher

import (
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

/*
handlers - исполнители (обработчики команд)
Реализуют, что делать по команде (/start, /help, /track)
Отвечают за бизнес-логику команды
(например, "Привет, я бот" или "добавить активность")
*/
func HandleStart(bot interfaces.StoppableBot, msg *tgbotapi.Message) {
	//We get the chat ID to know who to send the reply to.
	chatID := msg.Chat.ID
	//create inline buuton for ckoose
	buttonRu := tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru")
	buttonEN := tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang_en")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonRu, buttonEN),
	)
	//Sending a message with a language selection
	text := "Choose your language"
	message := tgbotapi.NewMessage(chatID, text)
	message.ReplyMarkup = keyboard
	//creating button text
	// buttonTrack := tgbotapi.NewKeyboardButton("/track") //a visual button in the keyboard.
	// buttonHelp := tgbotapi.NewKeyboardButton("/help")
	//merge two button in row for display ux
	// row := tgbotapi.NewKeyboardButtonRow(buttonTrack, buttonHelp)
	//creating object keyboard
	// keyboard := tgbotapi.NewReplyKeyboard(row)
	// creating message text
	// text := "Choose a action"                    //gets the user
	// message := tgbotapi.NewMessage(chatID, text) //chatID-to whom to send,text-What to say
	//attach the keyboard
	// message.ReplyMarkup = keyboard //"Show this keyboard along with the message."
	//send
	_, err := bot.Send(message) //Sends any text, emoji, command, button
	if err != nil {
		log.Error().Err(err).Msg("Error when sending the /start")
	}

}

//	func HandleTrack(bot interfaces.StoppableBot, msg *tgbotapi.Message) {
//		chatID := msg.Chat.ID
//		buttonTimerStandart := tgbotapi.NewKeyboardButton("Standart")
//		buttonTimerByYourself := tgbotapi.NewKeyboardButton("By yourself")
//		row := tgbotapi.NewKeyboardButtonRow(buttonTimerStandart, buttonTimerByYourself)
//		keyboard := tgbotapi.NewReplyKeyboard(row)
//		keyboard.ResizeKeyboard = true
//		text := "Choose"
//		message := tgbotapi.NewMessage(chatID, text)
//		message.ReplyMarkup = keyboard
//		_, err := bot.Send(message)
//		if err != nil {
//			log.Error().Err(err).Msg("eror track")
//		}
//	}
func ShowMainMenu(bot interfaces.StoppableBot, chatID int64) {

	buttonMyAccount := tgbotapi.NewKeyboardButton("👤My account")
	buttonTrack := tgbotapi.NewKeyboardButton("📈Track")
	buttonSupport := tgbotapi.NewKeyboardButton("🧠Learning")
	buttonSubscriptions := tgbotapi.NewKeyboardButton("💳Subscription")

	row1 := tgbotapi.NewKeyboardButtonRow(buttonMyAccount, buttonTrack)
	row2 := tgbotapi.NewKeyboardButtonRow(buttonSupport, buttonSubscriptions)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)

	//Уменьшает клавиатуру под контент (не будет занимать весь экран).
	keyboard.ResizeKeyboard = true
	//Означает, что клавиатура не исчезнет после одного нажатия.
	keyboard.OneTimeKeyboard = false

	msg := tgbotapi.NewMessage(chatID, "🏠Home")
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing menu")
	}
}
func ShowProfileMock(bot interfaces.StoppableBot, chatID int64) {
	text := `
👤 *My account*

— 👤 Имя: *Давид*
— 🧠 Активность: *Трекается*
— 🔥 Streak: *5 дней*
— 🌐 Язык: *Русский*
— 📍 Часовой пояс: *Europe/Berlin*
— 🗃 Подписка : *12 month*
— 📧 Контакт: @alaamov

Настроить профиль можно ниже:
`
	msg := tgbotapi.NewMessage(chatID, text)
	//Эта строка включает режим разметки Markdown
	msg.ParseMode = "Markdown"
	//Эта строка прикрепляет клавиатуру к сообщению,
	msg.ReplyMarkup = buildProfileKeyboard()
	//это вызов API Telegram,
	// который отправляет сообщение msg пользователю.
	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("err showing profil")

	}
}
func buildProfileKeyboard() tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🌐 Язык", "edit_language"),
		tgbotapi.NewInlineKeyboardButtonData("📍 Часовой пояс", "edit_timezone"),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📧 Контакт", "edit_contact"),
		tgbotapi.NewInlineKeyboardButtonData("🔁 Обновить", "refresh_profile"),
	)
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2)
}
func ShowTrackingMenu(bot interfaces.StoppableBot, chatID int64) {
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
func ShowLearningMenu(bot interfaces.StoppableBot, chatID int64) {
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
func buuildLerningKeyboard() tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Добавить слова на день", "add_wordsday"),
		tgbotapi.NewInlineKeyboardButtonData("🎲 Случайная подборка", "random_words"),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔁 Сменить язык", "switch_language"),
		tgbotapi.NewInlineKeyboardButtonData("📈 Статистика", "summary"),
	)
	row3 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🗂 База слов", "base_words"))
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)
}
