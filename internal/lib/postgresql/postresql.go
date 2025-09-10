package postgresql

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

type Storage struct{ db *sql.DB }

func New(config *config.Config) (*Storage, error) {
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

	storage := &Storage{db: db}

	return storage, nil
}
func (s *Storage) Close() error {
	return s.db.Close()
}

// ctx context.Context → нужен, чтобы можно было отменить
// или ограничить по времени запрос.
func (s *Storage) SaveUser(ctx context.Context, u *model.User) error {
	q := `
	INSERT INTO users (id, tg_user_id,username,phone_number,email,language,timezone,created_at)

	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	//Подготовка выражения
	//Prepare создаёт подготовленное выражение (prepared statement).

	//stmt (Statement) — это подготовленный запрос: база «поняла рецепт»,
	// проверила синтаксис и готова быстро готовить по нему много раз.

	//«зафиксировать шаблон» и подготовить его
	// к многократному использованию.
	stmt, err := s.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("failed prepare query %w", err)
	}
	defer stmt.Close()
	//Выполнение с параметрами

	//ExecContext — это выполнение подготовленного
	// запроса с конкретными ингредиентами (значениями).

	//«подставить конкретные значения и выполнить».
	_, err = stmt.ExecContext(ctx,
		u.ID,
		u.TgUserID,
		u.UserName,
		u.PhoneNumber,
		u.Email,
		u.Language,
		u.TimeZone,
		u.CreatedAt,
	)
	if err != nil { //Количество параметров = количеству плейсхолдеров.
		return fmt.Errorf("failed exec query to database %w", err)
	}
	return nil
}
func (s *Storage) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	q := `
	SELECT
	u.id,
	u.tg_user_id,
	u.username,
	u.phone_number,
	u.email,
	u.language,
	u.timezone,
	u.created_at
	FROM users AS u
	WHERE u.id = $1
	`
	stmt, err := s.db.Prepare(q) //Prepare лишь фиксирует форму запроса.
	if err != nil {
		return nil, fmt.Errorf("failed prepare query %w", err)
	}

	defer stmt.Close()

	//QueryContext — выполняет SELECT и возвращает курсор
	// (итератор по строкам результата).
	rows, err := stmt.QueryContext(ctx, id) //типы/ограничения реально проверяются на Exec/Query (во время выполнения).
	//rows — это не сами данные, а поток (чтение построчно).
	//Думаешь о нём как о файле: читаешь строку → следующую → и т.д.
	if err != nil {
		return nil, fmt.Errorf("failed exec query to database %w", err)
	}

	defer rows.Close()

	user := &model.User{}

	for rows.Next() { //rows.Next() — «есть ли ещё строка в результате?» Если да — переходим к ней.
		//«есть ли ещё следующая строка в результате этого запроса?»
		//если да → делает её «текущей», и ты можешь прочитать её поля через Scan(...);
		//если нет → цикл заканчивается.
		err := rows.Scan( //rows.Scan(&...) — перекладывает значения колонок текущей строки в поля твоей структуры.
			&user.ID, //Порядок должен соответствовать SELECT:
			&user.TgUserID,
			&user.UserName,
			&user.PhoneNumber,
			&user.Email,
			&user.Language,
			&user.TimeZone,
			&user.CreatedAt,
		) //добавляешь каждого user в слайс,
		//В твоём случае ты ищешь по ID (один пользователь) → цикл избыточен.
		// Лучше QueryRowContext
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows %w", err)
		}
	}

	return user, nil
}
