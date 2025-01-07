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

type ParseAccessJwtTokenUseCase struct {
	Users     repository.UserRepository
	JwtSecret string
}

func (uc *ParseAccessJwtTokenUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	accessToken string,
) (
	user *models.UserModel,
	err error,
) {
	log = log.With(slog.String("use-case", "ParseAccessJwtTokenUseCase"))

	tokenParams, err := jwtlib.Parse[JwtAccessTokenParams](accessToken, uc.JwtSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Error(
				"token is expired",
				"accessToken", accessToken,
				sl.Err(err),
			)
			return nil, models.ErrJwtExpired
		}
		log.Error(
			"error on parse access jwt token",
			"accessToken", accessToken,
			sl.Err(err),
		)
		return nil, models.ErrInvalidCredentials
	}

	user, err = uc.Users.GetByEmail(ctx, tokenParams.Email)
	if err != nil {
		log.Error(
			"error on get user info by email",
			"email", tokenParams.Email,
			sl.Err(err),
		)
		return nil, err
	}
	if user.ID == 0 {
		log.Info(
			"user not found",
			"email", tokenParams.Email,
		)
		return nil, models.ErrUserNotFound
	}

	return user, nil

}
