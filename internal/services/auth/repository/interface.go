package repository

import (
	"context"
	"sso/internal/domain/models"
)

type UserRepository interface {
	Save(
		ctx context.Context,
		email string,
		passHash []byte,
	) (err error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppRepository interface {
	App(ctx context.Context, appID int) (models.App, error)
}
