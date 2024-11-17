package auth_service

import (
	"context"
	"errors"
	"fmt"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"log/slog"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	authModels "sso/internal/services/auth/model"
	models "sso/internal/services/auth/model/errors"
	"sso/internal/services/auth/repository"
)

type ParseTokenUseCase struct {
	Users repository.UserRepository
}

func (uc *ParseTokenUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	token string,
	JWTSecret string,
) (
	user *authModels.UserModel,
	err error,
) {
	const op = "authService.ParseTokenUseCase"
	log = log.With(
		slog.String("op", op),
	)

	jwtInfo, err := jwt.ParseToken(token, JWTSecret)
	if err != nil {
		if errors.Is(err, jwtLib.ErrTokenExpired) {
			log.Error("jwt token expired", sl.Err(err))
			return &authModels.UserModel{}, jwtLib.ErrTokenExpired
		}
		log.Error("failed on parse jwt token", sl.Err(err))
		return &authModels.UserModel{}, models.ErrInvalidJwt
	}

	user, err = uc.Users.GetByEmail(ctx, jwtInfo.Email)
	if err != nil {
		log.Error("failed on get user by email", sl.Err(err))
		return &authModels.UserModel{}, fmt.Errorf("%s: %w", "ParseTokenUseCase", err)
	}
	if user.ID == 0 {
		return &authModels.UserModel{}, models.ErrUserNotFound
	}

	return user, nil
}
