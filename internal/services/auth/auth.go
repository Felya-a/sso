package auth

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/config"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"

	repo "sso/internal/services/auth/repository"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(
		ctx context.Context,
		userID int64,
	) (bool, error)
}

type AuthService struct {
	log            *slog.Logger
	appRepository  repo.AppRepository
	userRepository repo.UserRepository
}

// New returns a new instance of the Auth service.
func New(
	db *sqlx.DB,
	log *slog.Logger,
) *AuthService {
	userRepository := repo.NewPostgresUserRepository(db, log)
	appRepository := repo.NewPostgresAppRepository(db, log)

	return &AuthService{
		log:            log,
		appRepository:  appRepository,
		userRepository: userRepository,
	}
}

func (a *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (token string, err error) {
	const op = "authService.Login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("login user")

	user, err := a.userRepository.GetByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get user info from repository", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if user.ID == 0 {
		log.Warn("user not found")
		return "", models.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("failed to get user info from repository", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, models.ErrInvalidCredentials)
	}

	app, err := a.appRepository.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err = jwt.NetToken(user, app, config.Get().TokenTtl)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthService) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const op = "authService.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registering user")

	// Проверка наличия пользователя
	existingUser, err := a.userRepository.GetByEmail(ctx, email)
	if err != nil {
		log.Error("failed to check exists user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if existingUser.ID != 0 {
		return 0, models.ErrUserAlreadyExists
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Сохранение пользователя
	userId, err := a.userRepository.Save(ctx, email, hashedPassword)
	a.log.Debug("", userId)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (a *AuthService) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	// TODO
	return true, nil
}
