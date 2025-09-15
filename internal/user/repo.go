package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DavydAbbasov/trecker_bot/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
func (r *UserRepo) UpsertByTG(ctx context.Context, u *model.User) error {
	//	VALUES ($1,$2,$3,$4,$5,$6)- шесть плейсхолдеров. Значения для них
	// ты передашь из Go (в ExecContext(...)) в том же порядке:
	const q = `
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
func (r *UserRepo) GetByTG(ctx context.Context, tgID int64) (*model.User, error) {

	const q = `
	SELECT
id,tg_user_id, username,phone_number,email,language,timezone,created_at
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
