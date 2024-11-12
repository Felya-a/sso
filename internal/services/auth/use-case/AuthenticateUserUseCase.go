package auth_service

import (
	"context"
	"log/slog"
	"sso/internal/lib/logger/sl"
	models "sso/internal/services/auth/model/errors"

	auth "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticateUserUseCase struct {
	Users repository.UserRepository
}

func (uc *AuthenticateUserUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	email string,
	password string,
) (user *auth.UserModel, err error) {
	user, err = uc.Users.GetByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get user info from repository", sl.Err(err))
		return &auth.UserModel{}, models.ErrInternal
	}
	if user.ID == 0 {
		log.Info("user not found")
		return &auth.UserModel{}, models.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("failed on compare hash from password", sl.Err(err))
		return &auth.UserModel{}, models.ErrInvalidCredentials
	}

	return user, nil
}
