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

func main() {

	cfg := config.Load()

	tgBot, err := bot.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("error initialization bot : %v\n", err)
	}

	fmt.Println("START APP")

	dispatcher.New(tgBot)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	<-signalChan
	log.Info().Msg("telegram bot gracefully shutdown")

	waitForShutdown(tgBot, nil)

	log.Info().Msg("telegram bot is shutdown")
}

func waitForShutdown(tgBot interfaces.BotAPI, db *sql.DB) {

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
