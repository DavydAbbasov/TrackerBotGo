package interfaces

import (
	"context"

	"github.com/DavydAbbasov/trecker_bot/internal/model"
)

type UserRepository interface {
	CreateUserByTelegramID(ctx context.Context, u *model.User) error
	GetUserByTelegramID(ctx context.Context, tgID int64) (*model.User, error)
}
