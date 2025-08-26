package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/rs/zerolog/log"
)

type Config struct {
	TelegramToken    string `env:"TELEGRAM_TOKEN"`
	TelegramBotDebug bool   `env:"telegram_bot_debug"`
	DBHost           string `env:"db_host"`
	DBPORT           string `env:"db_port"`
}

func MustLoadConfig(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) { //
		log.Fatal().Msg("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil { //
		log.Fatal().Msg("cannot read config: " + err.Error())
	}

	return &cfg
}
