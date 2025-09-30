package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/DavydAbbasov/trecker_bot/internal/domain"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ActivityRepo struct {
	db *sql.DB
}

func NewActivityRepo(db *sql.DB) (*ActivityRepo, error) {
	return &ActivityRepo{
		db: db,
	}, nil
}
func (r *ActivityRepo) Create(ctx context.Context, userID int64, name string, emoji string) (domain.Activity, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Activity{}, fmt.Errorf("empty activity name")
	}
	emoji = strings.TrimSpace(emoji)
	q := `
	INSERT INTO activities (user_id,name, emoji)
	VALUES ($1,$2,$3)
	RETURNING id,user_id,name,emoji,is_archived,created_at;`

	var a domain.Activity
	err := r.db.QueryRowContext(ctx, q, userID, name, emoji).Scan(
		&a.ID,
		&a.UserID,
		&a.Name,
		&a.Emoji,
		&a.IsArchived, //зачем нам надо сканировать всю структуру и в цеелом зпачем сканировать
		&a.CreatedAt,
	)
	if err != nil {
		// maping errors postgres -> domain
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation (uq_activities_user_lower_name)//?
				return domain.Activity{}, domain.ErrActivityExists
			case "23514": // check_violation (если добавишь CHECK-и)
				return domain.Activity{}, fmt.Errorf("invalid activity: %s", pgErr.Message)
			}
		}
		return domain.Activity{}, err
	}

	return a, nil
}
func (r *ActivityRepo) ListActive(ctx context.Context, userID int64) ([]domain.Activity, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid userID: must be > 0")
	}
	q := `
	SELECT id, user_id, name, emoji,is_archived, created_at
	FROM activities
	WHERE user_id = $1 AND is_archived = false
	ORDER BY lower(name), id
	;`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list active activities: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Activity, 0, 16)
	for rows.Next() {
		var a domain.Activity
		if err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Name,
			&a.Emoji,
			&a.IsArchived,
			&a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan active activity: %w", err)
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration (active): %w", err)
	}

	return out, nil
}
func (r *ActivityRepo) SelectedListActive(ctx context.Context, userID int64) ([]int64, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid userID: must be > 0")
	}

	q := `
	SELECT activity_id
	FROM user_selected_activities
	WHERE user_id = $1
	ORDER BY activity_id;`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list selected ids query: %w", err)
	}
	defer rows.Close()

	ids := make([]int64, 0, 16)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan selected id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration (selected ids): %w", err)
	}
	return ids, nil
}
func (r *ActivityRepo) ToggleSelectedActive(ctx context.Context, userID int64, activityID int64) error {
	if userID <= 0 || activityID <= 0 {
		return fmt.Errorf("invalid ids")
	}

	// 1) проверка владения
	ownQ := `
	SELECT EXISTS(
    SELECT 1
    FROM activities
    WHERE id = $1 AND user_id = $2);`

	var owned bool
	if err := r.db.QueryRowContext(ctx, ownQ, activityID, userID).Scan(&owned); err != nil {
		return fmt.Errorf("ownership check: %w", err)
	}
	if !owned {
		return domain.ErrNoExistActivity
	}

	// 2) попытка снять выбор (delete-first)
	delQ := `
    DELETE
	FROM user_selected_activities
    WHERE user_id = $1 AND activity_id = $2;`

	res, err := r.db.ExecContext(ctx, delQ, userID, activityID)
	if err != nil {
		return fmt.Errorf("toggle delete: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 1 {
		// была выбрана → сняли и выходим
		return nil
	}

	// 3) не было выбрано → добавляем
	insQ := `
    INSERT INTO user_selected_activities(user_id, activity_id)
    VALUES ($1, $2)
    ON CONFLICT DO NOTHING;` // на случай гонок

	if _, err := r.db.ExecContext(ctx, insQ, userID, activityID); err != nil {
		return fmt.Errorf("toggle insert: %w", err)
	}
	return nil
}
