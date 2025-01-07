package auth_service

import (
	"context"
	"log/slog"
	"sso/internal/lib/logger/sl"

	models "sso/internal/services/auth/model"
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
) (user *models.UserModel, err error) {
	log = log.With(slog.String("use-case", "AuthenticateUserUseCase"))

	user, err = uc.Users.GetByEmail(ctx, email)
	if err != nil {
		log.Error(
			"failed to get user info from repository",
			"email", email,
			sl.Err(err),
		)
		return nil, err
	}
	if user.ID == 0 {
		log.Info(
			"user not found",
			"email", email,
		)
		return nil, models.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info(
			"failed on compare hash from password",
			"email", email,
			"password", password,
			sl.Err(err),
		)
		return nil, models.ErrInvalidCredentials
	}

	return user, nil
}
