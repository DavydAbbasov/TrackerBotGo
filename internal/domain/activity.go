package domain

import (
	"context"
	"time"
)

type Activity struct {
	ID         int64
	UserID     int64
	Name       string
	Emoji      string
	IsArchived bool
	CreatedAt  time.Time
}
type ActivityRepo interface {
	Create(ctx context.Context, userID int64, name string, emoji string) (Activity, error)
	ListActive(ctx context.Context, userID int64) ([]Activity, error)
	SelectedListActive(ctx context.Context, userID int64) ([]int64, error)
	ToggleSelectedActive(ctx context.Context, userID int64, activityID int64) error
}
