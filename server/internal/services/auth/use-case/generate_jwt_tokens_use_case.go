package auth_service

import (
	"context"
	"log/slog"
	"time"

	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	models "sso/internal/services/auth/model"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type GenerateJwtTokensUseCase struct {
	JwtSecret  string
	AccessTtl  time.Duration
	RefreshTtl time.Duration
}

type JwtAccessTokenParams struct {
	Id    int64  `json:"id" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type JwtRefreshTokenParams struct {
	Id int64 `json:"id" validate:"required"`
}

func (uc *GenerateJwtTokensUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	user *models.UserModel,
) (tokens *models.JwtTokens, err error) {
	log = log.With(slog.String("use-case", "GenerateJwtTokenUseCase"))

	accessToken, err := uc.generateAccessToken(
		log,
		uc.JwtSecret,
		&JwtAccessTokenParams{
			Id:    user.ID,
			Email: string(user.Email),
		},
	)
	if err != nil {
		log.Error(
			"error on generate access token",
			"id", user.ID,
			"email", user.Email,
			sl.Err(err),
		)
		return nil, err
	}

	refreshToken, err := uc.generateRefreshToken(
		log,
		uc.JwtSecret,
		&JwtRefreshTokenParams{
			Id: user.ID,
		},
	)
	if err != nil {
		log.Error(
			"error on generate refresh token",
			"id", user.ID,
			sl.Err(err),
		)
		return nil, err
	}

	return &models.JwtTokens{AccessJwtToken: accessToken, RefreshJwtToken: refreshToken}, nil
}

func (uc *GenerateJwtTokensUseCase) generateAccessToken(
	log *slog.Logger,
	secret string,
	params *JwtAccessTokenParams,
) (string, error) {
	claims := jwtlib.MapClaims{
		"id":    params.Id,
		"email": params.Email,
		"exp":   time.Now().Add(uc.AccessTtl).Unix(),
	}

	token, err := jwt.New(claims, secret)
	if err != nil {
		log.Error(
			"error during access token generation",
			"id", params.Id,
			"email", params.Email,
			sl.Err(err),
		)
		return "", err
	}

	return token, nil
}

func (uc *GenerateJwtTokensUseCase) generateRefreshToken(
	log *slog.Logger,
	secret string,
	params *JwtRefreshTokenParams,
) (string, error) {
	claims := jwtlib.MapClaims{
		"id":  params.Id,
		"exp": time.Now().Add(uc.RefreshTtl).Unix(),
	}

	token, err := jwt.New(claims, secret)
	if err != nil {
		log.Error(
			"error during refresh token generation",
			"id", params.Id,
			sl.Err(err),
		)
		return "", err
	}

	return token, err
}
