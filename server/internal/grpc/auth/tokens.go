package auth

import (
	"context"
	"log/slog"

	"sso/internal/lib/logger"
	models "sso/internal/services/auth/model"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokensRequestValidate struct {
	AuthorizationCode string `validate:"required"`
}

func (s *serverApi) Tokens(
	ctx context.Context,
	req *ssov1.TokensRequest,
) (*ssov1.TokensResponse, error) {
	log := logger.Logger()
	log = log.With(
		slog.String("requestid", uuid.New().String()),
	)

	dto := TokensRequestValidate{
		AuthorizationCode: req.GetAuthorizationCode(),
	}

	if err := validator.New().Struct(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokens, err := s.auth.Tokens(ctx, log, dto.AuthorizationCode)
	if err != nil {
		if models.IsDefinedError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.TokensResponse{
		AccessToken:  tokens.AccessJwtToken,
		RefreshToken: tokens.RefreshJwtToken,
	}, nil
}
