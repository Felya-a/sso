package auth

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/lib/logger/sl"
	"sso/internal/models"
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
		return &auth.UserModel{}, fmt.Errorf("%s: %w", "AuthenticateUserUseCase", err)
	}
	if user.ID == 0 {
		log.Warn("user not found")
		return &auth.UserModel{}, models.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("failed to get user info from repository", sl.Err(err))
		return &auth.UserModel{}, fmt.Errorf("%s: %w", "AuthenticateUserUseCase", models.ErrInvalidCredentials)
	}

	return user, nil
}
