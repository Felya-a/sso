package auth_service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"sso/internal/config"
	"sso/internal/lib/logger/sl"
	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	usecase "sso/internal/services/auth/use-case"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const AUTHORIZATION_CODE_TTL = 10 * time.Minute

type Auth interface {
	Register(
		ctx context.Context,
		log *slog.Logger,
		email string,
		password string,
	) (user *models.UserModel, err error)
	Login(
		ctx context.Context,
		log *slog.Logger,
		email string,
		password string,
		appID int,
	) (authorizationCode string, err error)
	Token(
		ctx context.Context,
		log *slog.Logger,
		authorizationCode string,
	) (tokens *models.JwtTokens, err error)
	Refresh(
		ctx context.Context,
		log *slog.Logger,
		refreshToken string,
	) (tokens *models.JwtTokens, err error)
	UserInfo(
		ctx context.Context,
		log *slog.Logger,
		token string,
	) (user *models.UserModel, err error)
}

type AuthService struct {
	registrationUser          usecase.RegistrationUserUseCase
	authenticateUser          usecase.AuthenticateUserUseCase
	generateAuthorizationCode usecase.GenerateAuthorizationCodeUseCase
	parseAuthorizationCode    usecase.ParseAuthorizationCodeUseCase
	generateJwtToken          usecase.GenerateJwtTokensUseCase
	parseAccessJwtToken       usecase.ParseAccessJwtTokenUseCase
	parseRefreshJwtToken      usecase.ParseRefreshJwtTokenUseCase
}

func New(
	db *sqlx.DB,
	rdb *redis.Client,
) *AuthService {
	config := config.Get()
	usersRepository := repository.NewPostgresUserRepository(db)
	authorizationCodeRepository := repository.NewRedisAuthorizationCodeRepository(rdb, AUTHORIZATION_CODE_TTL)

	registrationUser := usecase.RegistrationUserUseCase{
		Users: usersRepository,
	}
	authenticateUser := usecase.AuthenticateUserUseCase{
		Users: usersRepository,
	}
	generateJwtToken := usecase.GenerateJwtTokensUseCase{
		JwtSecret:  config.Jwt.Secret,
		AccessTtl:  config.Jwt.AccessTtl,
		RefreshTtl: config.Jwt.RefreshTtl,
	}
	generateAuthorizationCode := usecase.GenerateAuthorizationCodeUseCase{
		AuthorizationCodes: authorizationCodeRepository,
		JwtSecret:          config.Jwt.Secret,
	}
	parseAuthorizationCode := usecase.ParseAuthorizationCodeUseCase{
		AuthorizationCodes: authorizationCodeRepository,
		Users:              usersRepository,
		JwtSecret:          config.Jwt.Secret,
	}
	parseAccessJwtToken := usecase.ParseAccessJwtTokenUseCase{
		Users:     usersRepository,
		JwtSecret: config.Jwt.Secret,
	}
	parseRefreshJwtToken := usecase.ParseRefreshJwtTokenUseCase{
		Users:     usersRepository,
		JwtSecret: config.Jwt.Secret,
	}

	return &AuthService{
		registrationUser:          registrationUser,
		authenticateUser:          authenticateUser,
		generateJwtToken:          generateJwtToken,
		generateAuthorizationCode: generateAuthorizationCode,
		parseAuthorizationCode:    parseAuthorizationCode,
		parseAccessJwtToken:       parseAccessJwtToken,
		parseRefreshJwtToken:      parseRefreshJwtToken,
	}
}

func (a *AuthService) Register(
	ctx context.Context,
	log *slog.Logger,
	email string,
	password string,
) (*models.UserModel, error) {
	user, err := a.registrationUser.Execute(ctx, log, email, password)
	if err != nil {
		if models.IsDefinedError(err) {
			return &models.UserModel{}, err
		}
		log.Error("failed on registration new user", sl.Err(err))
		return &models.UserModel{}, fmt.Errorf("%s: %w", "AuthService", err)
	}

	log.Info("user success registered", "email", email)

	return user, nil
}

func (a *AuthService) Login(
	ctx context.Context,
	log *slog.Logger,
	email string,
	password string,
	appID int,
) (authorizationCode string, err error) {
	user, err := a.authenticateUser.Execute(ctx, log, email, password)
	if err != nil {
		if models.IsDefinedError(err) {
			return "", err
		}
		return "", fmt.Errorf("%s: %w", "AuthService", err)
	}

	authorizationCode, err = a.generateAuthorizationCode.Execute(ctx, log, user, AUTHORIZATION_CODE_TTL)
	if err != nil {
		if models.IsDefinedError(err) {
			return "", err
		}
		return "", fmt.Errorf("%s: %w", "AuthService", err)
	}

	return authorizationCode, nil
}

func (a *AuthService) Token(
	ctx context.Context,
	log *slog.Logger,
	authorizationCode string,
) (*models.JwtTokens, error) {
	user, err := a.parseAuthorizationCode.Execute(ctx, log, authorizationCode)
	if err != nil {
		if models.IsDefinedError(err) {
			return &models.JwtTokens{}, err
		}
		return &models.JwtTokens{}, fmt.Errorf("%s: %w", "AuthService", err)
	}

	tokens, err := a.generateJwtToken.Execute(ctx, log, user)
	if err != nil {
		if models.IsDefinedError(err) {
			return &models.JwtTokens{}, err
		}
		return &models.JwtTokens{}, fmt.Errorf("%s: %w", "AuthService", err)
	}

	return tokens, nil
}

func (a *AuthService) Refresh(
	ctx context.Context,
	log *slog.Logger,
	refreshToken string,
) (*models.JwtTokens, error) {
	user, err := a.parseRefreshJwtToken.Execute(ctx, log, refreshToken)
	if err != nil {
		if models.IsDefinedError(err) {
			return &models.JwtTokens{}, err
		}
		return &models.JwtTokens{}, fmt.Errorf("%s: %w", "AuthService", err)
	}

	tokens, err := a.generateJwtToken.Execute(ctx, log, user)
	if err != nil {
		if models.IsDefinedError(err) {
			return &models.JwtTokens{}, err
		}
		return &models.JwtTokens{}, fmt.Errorf("%s: %w", "AuthService", err)
	}

	return tokens, nil
}

func (a *AuthService) UserInfo(
	ctx context.Context,
	log *slog.Logger,
	token string,
) (*models.UserModel, error) {
	user, err := a.parseAccessJwtToken.Execute(ctx, log, token)
	if err != nil {
		if models.IsDefinedError(err) {
			return &models.UserModel{}, err
		}
		return &models.UserModel{}, fmt.Errorf("%s: %w", "AuthService", err)
	}

	return user, nil
}
