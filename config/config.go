package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var TelegramToken string

type Config struct {
	TelegramToken    string `env:"telegram_token"`
	TelegramBotDebug bool   `env:"telegram_bot_debug"`

	DBHost string `env:"db_host"` // data base
	DBPORT string `env:"db_port"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	TelegramToken = os.Getenv("TELEGRAM_TOKEN")

	if TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is not set in the environment variables")
	}

	return &Config{}
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	TelegramToken = os.Getenv("TELEGRAM_TOKEN")

	if TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is not set in the environment variables")
	}

	return &Config{}
}
