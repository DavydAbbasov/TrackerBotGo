package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var TelegramToken string

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	TelegramToken = os.Getenv("TELEGRAM_TOKEN")

	if TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is not set in the environment variables")
	}
	
}
