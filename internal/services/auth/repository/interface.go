package repository

import (
	"context"
	auth "sso/internal/services/auth/model"
)

type UserRepository interface {
	Save(
		ctx context.Context,
		email string,
		passHash []byte,
	) (err error)

	GetByEmail(
		ctx context.Context,
		email string,
	) (*auth.UserModel, error)
}
