package postgresql

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

	"github.com/DavydAbbasov/trecker_bot/config"
	"github.com/DavydAbbasov/trecker_bot/internal/model"
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
	log.Info().Msg("database connection is success")

	storage := &UserRepo{db: db}

	return storage, nil
}

func (s *UserRepo) Close() error {
	return s.db.Close()
}

func (r *UserRepo) CreateUserByTelegramID(ctx context.Context, u *model.User) error {
	//	VALUES ($1,$2,$3,$4,$5,$6)- шесть плейсхолдеров. Значения для них
	// ты передашь из Go (в ExecContext(...)) в том же порядке:
	q := `
	INSERT INTO users (tg_user_id, username, phone_number, email, language, timezone)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (tg_user_id) DO UPDATESET
		username     = COALESCE(EXCLUDED.username, users.username),
		phone_number = COALESCE(EXCLUDED.phone_number, users.phone_number),
		email        = COALESCE(EXCLUDED.email, users.email),
		language     = COALESCE(EXCLUDED.language, users.language),
		timezone     = COALESCE(EXCLUDED.timezone, users.timezone)
	;`

	_, err := r.db.ExecContext(ctx, q,
		u.TgUserID, u.UserName, u.PhoneNumber, u.Email, u.Language, u.TimeZone)
	if err != nil {
		return fmt.Errorf("user upsert:%w", err)
	}
	return nil
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
	err := r.db.QueryRowContext(ctx, q, tgID). //“обратись к базе и верни одну строку”
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

func (r *UserRepo) UpdateUserByTelegramID(ctx context.Context, u *model.User) error {
	// not implemented
	return nil
}

func (r *UserRepo) DeleteUserByTelegramID(ctx context.Context, telegramID int64) error {
	// not implemented
	return nil
}
