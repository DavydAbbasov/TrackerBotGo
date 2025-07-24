package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	TelegramToken    string `env:"telegram_token"`   
	TelegramBotDebug bool   `env:"telegram_bot_debug"`  
	DBHost           string `env:"db_host"`            
	DBPORT           string `env:"db_port"`            

}


func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		TelegramToken:    os.Getenv("TELEGRAM_TOKEN"),
		TelegramBotDebug: os.Getenv("TELEGRAM_BOT_DEBUG") == "true",
		DBHost:           os.Getenv("DB_HOST"),
		DBPORT:           os.Getenv("DB_PORT"),
	}

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is not set in the environment variables")
	}
	return cfg
}

