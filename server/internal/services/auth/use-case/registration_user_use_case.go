package auth_service

import (
	"context"
	"log/slog"
	"sso/internal/lib/logger/sl"
	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationUserUseCase struct {
	mu    sync.Mutex
	Users repository.UserRepository
}

func (uc *RegistrationUserUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	email string,
	password string,
) (user *models.UserModel, err error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	log = log.With(slog.String("use-case", "RegistrationUserUseCase"))

	// Проверка наличия пользователя
	existingUser, err := uc.Users.GetByEmail(ctx, email)
	if err != nil {
		log.Error(
			"failed to check exists user",
			"email", email,
			sl.Err(err),
		)
		return nil, err
	}
	if existingUser.ID != 0 {
		log.Info(
			"user already exist",
			"email", email,
		)
		return nil, models.ErrUserAlreadyExists
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(
			"failed to generate password hash",
			"password", password,
			sl.Err(err),
		)
		return nil, err
	}

	// Сохранение пользователя
	if err := uc.Users.Save(ctx, email, hashedPassword); err != nil {
		log.Error(
			"failed to save user",
			"email", email,
			sl.Err(err),
		)
		return nil, err
	}

	// Получение id созданного пользователя
	user, err = uc.Users.GetByEmail(ctx, email)
	if err != nil {
		log.Error(
			"failed to get new user id",
			"email", email,
			sl.Err(err),
		)
		return nil, err
	}
	if user.ID == 0 {
		log.Error(
			"user not was saved",
			"email", email,
		)
		return nil, models.ErrUserNotFound
	}

	return user, nil
}
