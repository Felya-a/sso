package auth_service

import (
	"context"
	"log/slog"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	"strconv"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type GenerateAuthorizationCodeUseCase struct {
	AuthorizationCodes repository.AuthorizationCodeRepository
	JwtSecret          string
}

type AuthorizationCodeParams struct {
	Id    string `validate:"required"`
	Email string `validate:"required"`
}

func (uc *GenerateAuthorizationCodeUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	user *models.UserModel,
	ttlCode time.Duration,
) (authorizationCode string, err error) {
	log = log.With(slog.String("use-case", "GenerateAuthorizationCodeUseCase"))

	codeParams := &AuthorizationCodeParams{
		Id:    strconv.FormatInt(user.ID, 10),
		Email: string(user.Email),
	}

	claims := jwtlib.MapClaims{
		"id":    codeParams.Id,
		"email": codeParams.Email,
		"exp":   time.Now().Add(ttlCode).Unix(),
	}

	authorizationCode, err = jwt.New(claims, uc.JwtSecret)
	if err != nil {
		log.Error(
			"error during generation authorization token",
			"id", user.ID,
			"email", user.Email,
			sl.Err(err),
		)
		return "", err
	}

	if err = uc.AuthorizationCodes.Save(ctx, authorizationCode); err != nil {
		log.Error(
			"error on save authorization code",
			"authorizationCode", authorizationCode,
			sl.Err(err),
		)
		return "", err
	}

	return authorizationCode, nil
}
