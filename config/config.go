package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/rs/zerolog/log"
)

type Config struct {
	TelegramToken    string `env:"telegram_token"`     //нужен, чтобы подключить бота.ключ авторизации Telegram-бота.
	TelegramBotDebug bool   `env:"telegram_bot_debug"` //включает/выключает отладку.Например, ты увидишь, какие апдейты приходят, как отправляются сообщения.
	DBHost           string `env:"db_host"`            //Где (на каком сервере) живёт база
	DBPORT           string `env:"db_port"`            //Через какой порт к ней подключаться
}

func MustLoadConfig(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Msg("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal().Msg("cannot read config: " + err.Error())
	}

	return &cfg
}
