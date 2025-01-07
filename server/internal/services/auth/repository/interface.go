package repository

import (
	"context"
	auth "sso/internal/services/auth/model"
)

type UserRepository interface {
	GetById(
		ctx context.Context,
		id int64,
	) (*auth.UserModel, error)
	GetByEmail(
		ctx context.Context,
		email string,
	) (*auth.UserModel, error)
	Save(
		ctx context.Context,
		email string,
		passHash []byte,
	) (err error)
}

type AuthorizationCodeRepository interface {
	CheckEndDelete(
		ctx context.Context,
		code string,
	) (bool, error)
	Save(
		ctx context.Context,
		code string,
	) error
}
