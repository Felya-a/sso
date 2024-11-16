package auth_service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"sso/internal/config"
	"sso/internal/lib/logger/sl"
	authModels "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	usecase "sso/internal/services/auth/use-case"

	"github.com/jmoiron/sqlx"
)

type Auth interface {
	Login(
		ctx context.Context,
		log *slog.Logger,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		log *slog.Logger,
		email string,
		password string,
	) (userID int64, err error)
	GetUserInfo(
		ctx context.Context,
		log *slog.Logger,
		token string,
	) (user *authModels.UserModel, err error)
}

type AuthService struct {
	authenticateUser usecase.AuthenticateUserUseCase
	generateToken    usecase.GenerateTokenUseCase
	parseToken       usecase.ParseTokenUseCase
	registrationUser usecase.RegistrationUserUseCase
	mu               sync.Mutex
}

func New(
	db *sqlx.DB,
) *AuthService {
	userRepository := repository.NewPostgresUserRepository(db)
	registrationUser := usecase.RegistrationUserUseCase{Users: userRepository}
	authenticateUser := usecase.AuthenticateUserUseCase{Users: userRepository}
	generateToken := usecase.GenerateTokenUseCase{TokenTtl: config.Get().TokenTtl}
	parseToken := usecase.ParseTokenUseCase{Users: userRepository}

	return &AuthService{
		authenticateUser: authenticateUser,
		generateToken:    generateToken,
		parseToken:       parseToken,
		registrationUser: registrationUser,
	}
}

func (a *AuthService) Login(
	ctx context.Context,
	log *slog.Logger,
	email string,
	password string,
	appID int,
) (token string, err error) {
	const op = "authService.Login"
	log = log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	user, err := a.authenticateUser.Execute(ctx, log, email, password)
	if err != nil {
		log.Error("failed on get user info", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = a.generateToken.Execute(ctx, log, user, config.Get().JWTSecret)
	if err != nil {
		log.Error("failed on generate jwt token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthService) RegisterNewUser(
	ctx context.Context,
	log *slog.Logger,
	email string,
	password string,
) (userID int64, err error) {
	const op = "authService.RegisterNewUser"
	log = log.With(
		slog.String("op", op),
	)

	a.mu.Lock()
	defer a.mu.Unlock()

	log.Info("registration new user", "email", email)

	user, err := a.registrationUser.Execute(ctx, log, email, password)
	if err != nil {
		log.Error("failed on registration new user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user success registered", "email", email)

	return user.ID, nil
}

func (a *AuthService) GetUserInfo(
	ctx context.Context,
	log *slog.Logger,
	token string,
) (user *authModels.UserModel, err error) {
	const op = "authService.GetUserInfo"
	log = log.With(
		slog.String("op", op),
	)

	log.Info("try parse jwt token", "token", token)

	user, err = a.parseToken.Execute(ctx, log, token, config.Get().JWTSecret)
	if err != nil {
		log.Error("failed on parse jwt token", sl.Err(err))
		return &authModels.UserModel{}, err
	}

	log.Info("jwt success parsed", "userId", user.ID)

	return user, err
}
