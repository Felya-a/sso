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

type RefreshRequestValidate struct {
	RefreshToken string `validate:"required"`
}

func (s *serverApi) Refresh(
	ctx context.Context,
	req *ssov1.RefreshRequest,
) (*ssov1.RefreshResponse, error) {
	log := logger.Logger()
	requestID := uuid.New().String()
	log = log.With(
		slog.String("requestid", requestID),
	)

	log.Info("Refresh request started", slog.Any("request", req))

	dto := RefreshRequestValidate{
		RefreshToken: req.GetRefreshToken(),
	}

	if err := validator.New().Struct(dto); err != nil {
		log.Error("Validation failed", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, models.ErrInvalidCredentials.Error())
	}

	log.Info("Validation passed", slog.Any("dto", dto))

	tokens, err := s.auth.Refresh(ctx, log, dto.RefreshToken)
	if err != nil {
		if models.IsDefinedError(err) {
			log.Error("Refresh failed with defined error", slog.Any("error", err))
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		log.Error("Refresh failed with internal error", slog.Any("error", err))
		return nil, status.Error(codes.Internal, "Internal error")
	}

	log.Info("Refresh request completed successfully", slog.Any("tokens", tokens))

	return &ssov1.RefreshResponse{
		AccessToken:  tokens.AccessJwtToken,
		RefreshToken: tokens.RefreshJwtToken,
	}, nil
}
