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
)

func main() {
	config.LoadConfig()

	tgBot, err := bot.InitBot(config.TelegramToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialization bot : %v\n", err)
		os.Exit(1)
	}

	SetupGracefulShutdown(tgBot, nil)
}

func SetupGracefulShutdown(bot *tgbotapi.BotAPI, db *sql.DB) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("The completion signal has been received. Completing the work.")

	if db != nil {
		err := db.Close()
		if err != nil {
			fmt.Println("Error when closing the database:", err)
		} else {
			fmt.Println("The database connection is closed.")
		}
	}
	if bot != nil {
		bot.StopReceivingUpdates()
		fmt.Println("Bot Stoped")
	}
}
