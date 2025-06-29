package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DavydAbbasov/trecker_bot/config"
	"github.com/DavydAbbasov/trecker_bot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

// var wg sync.WaitGroup

func main() {
	cfg := config.MustLoad()

	_, err := bot.New(cfg)
	if err != nil {
		log.Fatal().Msgf("error initialization bot : %v\n", err)
	}

	fmt.Println("START APP")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	<-signalChan
	log.Info().Msg("telegram bot gracefully shutdown")

	// Функции или код, который завершает работу

	log.Info().Msg("telegram bot is shutdown")
}

func SetupGracefulShutdown(bot *tgbotapi.BotAPI, db *sql.DB) {

	fmt.Println("The completion signal has been received. Completing the work.")

	// if db != nil {
	// 	err := db.Close()
	// 	if err != nil {
	// 		fmt.Println("Error when closing the database:", err)
	// 	} else {
	// 		fmt.Println("The database connection is closed.")
	// 	}
	// }
	// if bot != nil {
	// 	bot.StopReceivingUpdates()
	// 	fmt.Println("Bot Stoped")
	// }
}
