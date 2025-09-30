package storage

import (
	"database/sql"
	"fmt"

	"github.com/DavydAbbasov/trecker_bot/config"
	storage "github.com/DavydAbbasov/trecker_bot/storage/postgres"

	"github.com/DavydAbbasov/trecker_bot/internal/domain"
	_ "github.com/lib/pq"
	log "github.com/rs/zerolog/log"
)

type Repos struct {
	db         *sql.DB
	Activities domain.ActivityRepo
}

func NewRepos(cfg *config.Config, dsn string) (*Repos, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error get database driver %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connection to database %w", err)
	}

	log.Info().Msg("database Repos connection is success")

	r := &Repos{
		db: db,
	}
	r.Activities, err = storage.NewActivityRepo(db)
	if err != nil {
		log.Error().Err(err).Msg("shuher")
	}
	return r, nil
}

func (r *Repos) Close() error {
	return r.db.Close()
}
