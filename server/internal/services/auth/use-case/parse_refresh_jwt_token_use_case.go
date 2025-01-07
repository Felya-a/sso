package auth_service

import (
	"context"
	"errors"
	"log/slog"
	jwtlib "sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"

	"github.com/golang-jwt/jwt/v5"
)

type ParseRefreshJwtTokenUseCase struct {
	Users     repository.UserRepository
	JwtSecret string
}

func (uc *ParseRefreshJwtTokenUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	refreshToken string,
) (
	user *models.UserModel,
	err error,
) {
	log = log.With(slog.String("use-case", "ParseRefreshJwtTokenUseCase"))

	tokenParams, err := jwtlib.Parse[JwtRefreshTokenParams](refreshToken, uc.JwtSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Error(
				"token is expired",
				"refreshToken", refreshToken,
				sl.Err(err),
			)
			return nil, models.ErrJwtExpired
		}
		log.Error(
			"error on parse refresh jwt token",
			"refreshToken", refreshToken,
			sl.Err(err),
		)
		return nil, models.ErrInvalidCredentials
	}

	user, err = uc.Users.GetById(ctx, tokenParams.Id)
	if err != nil {
		log.Error(
			"error on get user info by id",
			"user.id", tokenParams.Id,
			sl.Err(err),
		)
		return nil, err
	}
	if user.ID == 0 {
		log.Error(
			"user not found",
			"user.id", tokenParams.Id,
		)
		return nil, models.ErrUserNotFound
	}

	return user, nil

}
