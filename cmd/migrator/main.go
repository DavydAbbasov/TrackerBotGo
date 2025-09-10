package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/DavydAbbasov/trecker_bot/config"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //миграции можно применять к PostgreSQL.
	_ "github.com/golang-migrate/migrate/v4/source/file"       //миграции можно читать из файловой системы
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	cfg := config.MustLoadConfig(".env")

	user := url.QueryEscape(cfg.UserDB)
	pass := url.QueryEscape(cfg.PasswordDB)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		pass,
		cfg.HostDB,
		cfg.PortDB,
		cfg.NameDB,
		cfg.SSLMode,
	)

	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		panic(err)
	}

	defer m.Close()

	if err := m.Up(); err != nil {

		if errors.Is(err, migrate.ErrNoChange) { //Это ошибка, которую возвращает
			// библиотека, если нет новых миграций.
			//миграции в папке есть, но все они уже применены;
			//новых файлов .up.sql не появилось.
			//Здесь мы говорим: «Если ошибка — это именно ErrNoChange
			// (прямо или внутри другой ошибки) → не паникуем, просто логируем».

			log.Info().Msg("no migrations to apply")

			return
		}
		panic(err)
	}
	log.Info().Msg("migrations applied")
}
