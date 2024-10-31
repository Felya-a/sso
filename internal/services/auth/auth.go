package auth

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jmoiron/sqlx"

	"sso/internal/config"
	"sso/internal/lib/logger/sl"
	"sso/internal/services/auth/repository"
	usecase "sso/internal/services/auth/use-case"
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
}

type AuthService struct {
	log              *slog.Logger
	authenticateUser usecase.AuthenticateUserUseCase
	generateToken    usecase.GenerateTokenUseCase
	registrationUser usecase.RegistrationUserUseCase
	mu               sync.Mutex
}

func New(
	db *sqlx.DB,
	log *slog.Logger,
) *AuthService {
	userRepository := repository.NewPostgresUserRepository(db, log)
	registrationUser := usecase.RegistrationUserUseCase{Users: userRepository}
	authenticateUser := usecase.AuthenticateUserUseCase{Users: userRepository}
	generateToken := usecase.GenerateTokenUseCase{TokenTtl: config.Get().TokenTtl}

	return &AuthService{
		log:              log,
		authenticateUser: authenticateUser,
		generateToken:    generateToken,
		registrationUser: registrationUser,
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

	user, err := a.authenticateUser.Execute(ctx, log, email, password)
	if err != nil {
		log.Error("failed on get user info", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = a.generateToken.Execute(ctx, log, user, config.Get().TokenTtl, config.Get().JWTSecret)
	if err != nil {
		log.Error("failed on generate jwt token", sl.Err(err))
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

	a.mu.Lock()
	defer a.mu.Unlock()

	user, err := a.registrationUser.Execute(ctx, log, email, password)
	if err != nil {
		log.Error("failed on registration new user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, nil
}
