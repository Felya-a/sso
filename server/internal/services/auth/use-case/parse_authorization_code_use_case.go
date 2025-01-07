package auth_service

import (
	"context"
	"log/slog"
	jwtlib "sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
)

type ParseAuthorizationCodeUseCase struct {
	Users              repository.UserRepository
	AuthorizationCodes repository.AuthorizationCodeRepository
	JwtSecret          string
}

func (uc *ParseAuthorizationCodeUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	authorizationCode string,
) (user *models.UserModel, err error) {
	log = log.With(slog.String("use-case", "ParseAuthorizationCodeUseCase"))

	isExist, err := uc.AuthorizationCodes.CheckEndDelete(ctx, authorizationCode)
	if err != nil {
		log.Error(
			"error on check exists authorization code",
			"authorizationCode", authorizationCode,
			sl.Err(err),
		)
		return nil, err
	}

	if !isExist {
		log.Info(
			"authorization code is not exists",
			"authorizationCode", authorizationCode,
		)
		return nil, models.ErrInvalidCredentials
	}

	authorizationCodeParams, err := jwtlib.Parse[AuthorizationCodeParams](authorizationCode, uc.JwtSecret)
	if err != nil {
		log.Error(
			"error on parse authorization code",
			"authorizationCode", authorizationCode,
			sl.Err(err),
		)
		return nil, err
	}

	user, err = uc.Users.GetByEmail(ctx, authorizationCodeParams.Email)
	if err != nil {
		log.Error(
			"error on get user by email",
			"email", authorizationCodeParams.Email,
			sl.Err(err),
		)
		return nil, err
	}
	if user.ID == 0 {
		log.Error(
			"user not found",
			"email", authorizationCodeParams.Email,
		)
		return nil, models.ErrUserNotFound
	}

	return user, nil
}
