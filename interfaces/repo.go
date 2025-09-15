package interfaces

import (
	"context"

	"github.com/DavydAbbasov/trecker_bot/internal/model"
)

type Repo interface {
	UpsertByTG(ctx context.Context, u *model.User) error
	GetByTG(ctx context.Context, tgID int64) (*model.User, error)
}
