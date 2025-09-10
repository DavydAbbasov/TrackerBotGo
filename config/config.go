package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/rs/zerolog/log"
)

type Config struct {
	TelegramToken    string `env:"TELEGRAM_TOKEN"`
	TelegramBotDebug bool   `env:"TELEGRAM_BOT_DEBUG"`

	HostDB     string `env:"HOST_DB"`
	PortDB     int64  `env:"PORT_DB"`
	UserDB     string `env:"USER_DB"`
	PasswordDB string `env:"PASSWORD_DB"`
	NameDB     string `env:"NAME_DB"`
	SSLMode    string `env:"SSL_MODE"`
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
