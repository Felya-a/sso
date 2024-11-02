package auth

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	authModels "sso/internal/services/auth/model"
	"time"
)

type GenerateTokenUseCase struct {
	TokenTtl time.Duration
}

func (uc *GenerateTokenUseCase) Execute(
	ctx context.Context,
	log *slog.Logger,
	user *authModels.UserModel,
	JWTSecret string,
) (token string, err error) {
	token, err = jwt.NewToken(jwt.JwtBodyParams{ID: user.ID, Email: user.Email}, uc.TokenTtl, JWTSecret)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", "GenerateTokenUseCase", err)
	}
	return token, nil
}
