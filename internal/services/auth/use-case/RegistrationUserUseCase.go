package auth_service

import (
	"context"
	"log/slog"
	"sso/internal/lib/logger/sl"
	auth "sso/internal/services/auth/model"
	models "sso/internal/services/auth/model/errors"
	"sso/internal/services/auth/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationUserUseCase struct {
	Users repository.UserRepository
}

func (uc *RegistrationUserUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	email string,
	password string,
) (user *auth.UserModel, err error) {
	// Проверка наличия пользователя
	existingUser, err := uc.Users.GetByEmail(ctx, email)
	if err != nil {
		log.Error("failed to check exists user", sl.Err(err))
		return &auth.UserModel{}, models.ErrInternal
	}
	if existingUser.ID != 0 {
		return &auth.UserModel{}, models.ErrUserAlreadyExists
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return &auth.UserModel{}, models.ErrInternal
	}

	// Сохранение пользователя
	if err := uc.Users.Save(ctx, email, hashedPassword); err != nil {
		log.Error("failed to save user", sl.Err(err))
		return &auth.UserModel{}, models.ErrInternal
	}

	// Получение id созданного пользователя
	user, err = uc.Users.GetByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get new user id", sl.Err(err))
		return &auth.UserModel{}, models.ErrInternal
	}
	if user.ID == 0 {
		log.Error("failed to save new user")
		return &auth.UserModel{}, models.ErrUserNotSaved
	}

	return user, nil
}
