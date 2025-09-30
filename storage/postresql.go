package storage

//Storage (твой файл) — это инициализация соединения с БД
// и управление пулом *sql.DB (Ping, настройки, Close).

//Адаптер — это слой, который умеет «разговаривать»
// с внешними системами (в нашем случае — с базой данных).

//«Адаптер хранилища» как бы переводит бизнес-команды (сохрани пользователя,
// получи пользователя) в SQL-запросы.«Адаптер хранилища» как бы переводит
// бизнес-команды (сохрани пользователя, получи пользователя) в SQL-запросы.
import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/DavydAbbasov/trecker_bot/config"

	helper "github.com/DavydAbbasov/trecker_bot/internal/lib/postgresql"
	"github.com/DavydAbbasov/trecker_bot/internal/model"
	_ "github.com/lib/pq"
	log "github.com/rs/zerolog/log"
)

type UserRepo struct {
	db *sql.DB
}

func New(config *config.Config) (*UserRepo, error) {
	//DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s ", config.HostDB, config.PortDB, config.UserDB,
		config.PasswordDB, config.NameDB, config.SSLMode)

	//Важно: sql.Open не устанавливает сетевое соединение сразу. Он:
	//проверяет базовую валидность параметров,
	//создаёт объект *sql.DB — это пул соединений (менеджер коннекшенов).
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error get database driver %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connection to database %w", err)
	}
	//Ping() принудительно открывает реальное соединение
	// и проверяет, что база доступна и DSN валиден.
	log.Info().Msg("database UserRepo connection is success")

	storage := &UserRepo{
		db: db,
	}

	return storage, nil
}

func (s *UserRepo) Close() error {
	return s.db.Close()
}

// EnsureIDByTelegram: найдёт или создаст пользователя и вернёт users.id
func (r *UserRepo) EnsureIDByTelegram(ctx context.Context, tgID int64, username string) (int64, error) {
	username = strings.TrimSpace(username)

	ns := sql.NullString{
		String: username,
		Valid:  username != "",
	}

	q := `
	INSERT INTO users (tg_user_id, username)
	VALUES ($1, $2)
	ON CONFLICT (tg_user_id) DO UPDATE
	SET username = COALESCE(EXCLUDED.username, users.username)
	RETURNING id
	`

	var id int64
	if err := r.db.QueryRowContext(ctx, q, tgID, ns).
		Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepo) GetUserByTelegramID(ctx context.Context, tgID int64) (*model.User, error) {

	q := `
	SELECT
		id,tg_user_id,
		username,
		phone_number,
		email,
		language,
		timezone,
		created_at
	FROM users
	WHERE tg_user_id = $1;`

	var u model.User
	//QueryRowContext возвращает «обёртку» для одной строки (*Row
	err := r.db.QueryRowContext(ctx, q, tgID).
		Scan(
			&u.ID,
			&u.TgUserID,
			&u.UserName,
			&u.PhoneNumber,
			&u.Email,
			&u.Language,
			&u.TimeZone,
			&u.CreatedAt,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // нет пользователя — вернём nil
		}
		return nil, fmt.Errorf("user get by tg: %w", err)
	}
	return &u, nil
}

// InsertUser Adding a user to the database.
func (r *UserRepo) InsertUser(ctx context.Context, tgID int64, username *string) error {
	// Request to add data.
	q :=
		`INSERT INTO users (tg_user_id,username)
	VALUES ($1,$2)
	ON CONFLICT (tg_user_id)DO NOTHING;`
	// Executing a request to add data.
	if _, err := r.db.ExecContext(ctx, q, tgID, helper.ToNullable(username)); err != nil {
		return err
	}
	return nil
}
func (r *UserRepo) UpdateUsername(ctx context.Context, tgID int64, username string) (int64, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return 0, nil
	}

	q :=
		`UPDATE users
	SET username = $1
	WHERE tg_user_id = $2
	AND (username IS DISTINCT FROM $1);`

	result, err := r.db.ExecContext(ctx, q, username, tgID)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err
}
func (r *UserRepo) UpdateLanguage(ctx context.Context, tgID int64, lang string) (int64, error) {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return 0, nil
	}

	q :=
		`UPDATE users
	SET language = $2
	WHERE tg_user_id = $1
	AND (language IS DISTINCT FROM $2 );`

	result, err := r.db.ExecContext(ctx, q, tgID, lang)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err
}
func (r *UserRepo) UpdateUserByTelegramID(ctx context.Context, u *model.User) error {
	// not implemented
	return nil
}

func (r *UserRepo) DeleteUserByTelegramID(ctx context.Context, telegramID int64) error {
	// not implemented
	return nil
}

// func (r *UserRepo) CheckIfUserExist(ctx context.Context, tgID int64) (bool, error) {
// 	q := `SELECT COUNT(id) AS countusers FROM users WHERE tg_user_id;`
// }

// func (r *UserRepo) UpdateEmail(ctx context.Context, email string, id int64) (int64, error) {
// }
