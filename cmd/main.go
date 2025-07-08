package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DavydAbbasov/trecker_bot/config"
	"github.com/DavydAbbasov/trecker_bot/internal/bot"
	"github.com/DavydAbbasov/trecker_bot/internal/dispatcher"
	"github.com/DavydAbbasov/trecker_bot/pkg/interfaces"

	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

// var wg sync.WaitGroup
// Здесь управляется жизненный цикл приложения.
func main() {
	// 1. Загружаем конфигурацию из .env
	cfg := config.Load()
	/*Загрузи .env файл
	2.Прочитай из него переменные
	3.Помести всё в структуру
	4.Config Сохрани в переменную cfg
	*/

	// 2. Инициализируем Telegram-бота
	tgBot, err := bot.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("error initialization bot : %v\n", err)
	}

	fmt.Println("START APP")

	dispatcher.Start(tgBot)

	// 3. Запускаем обработку команд (логика команд — в dispatcher)
	// go dispatcher.Start(tgBot)
	// 4. Ожидаем сигнал завершения (Ctrl+C, kill и т.д.)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	<-signalChan
	log.Info().Msg("telegram bot gracefully shutdown")

	// 5. Завершаем соединения (бот, БД)
	waitForShutdown(tgBot, nil)

	log.Info().Msg("telegram bot is shutdown")
}

// Завершение соединений (бот и база данных)
func waitForShutdown(tgBot interfaces.StoppableBot, db *sql.DB) {

	fmt.Println("The completion signal has been received. Completing the work.")

	if db != nil {
		err := db.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing DB connection")
		} else {
			log.Info().Msg("db connection is closed.")
		}
	}
	if tgBot != nil {
		tgBot.StopReceivingUpdates()
		log.Info().Msg("Telegram bot stopped")
	}
}
