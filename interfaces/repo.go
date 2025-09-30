package interfaces

import (
	"context"

	"github.com/DavydAbbasov/trecker_bot/internal/model"
)

type UserRepository interface {
	EnsureIDByTelegram(ctx context.Context, tgID int64, username string) (int64, error)
	InsertUser(ctx context.Context, tgID int64, username *string) error
	// UpdateUser(ctx context.Context, u *model.User) error
	GetUserByTelegramID(ctx context.Context, tgID int64) (*model.User, error)
	UpdateUsername(ctx context.Context, tgID int64, uname string) (int64, error)
	UpdateLanguage(ctx context.Context, tgID int64, lang string) (int64, error)
}
