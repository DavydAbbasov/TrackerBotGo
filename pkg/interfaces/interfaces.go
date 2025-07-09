package interfaces
/*
Изоляция зависимостей
Гибкость замены реализации
Чистая архитектура и SOLID
«Модули не должны зависеть от деталей реализации,
а только от абстракций» — принцип Dependency Inversion (из SOLID)
Это метод из библиотеки tgbotapi, который:
Останавливает получение обновлений (сообщений) от Telegram.
*/
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StoppableBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	StopReceivingUpdates()
}
