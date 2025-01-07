package auth

import (
	"context"
	"errors"
	"log/slog"

	"sso/internal/lib/logger"
	models "sso/internal/services/auth/model"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginRequestValidate struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	AppId    int    `validate:"required"`
}

func (s *serverApi) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	log := logger.Logger()
	log = log.With(
		slog.String("requestid", uuid.New().String()),
	)

	dto := LoginRequestValidate{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppId:    int(req.GetAppId()),
	}

	if err := validator.New().Struct(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	authorizationCode, err := s.auth.Login(ctx, log, dto.Email, dto.Password, dto.AppId)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return nil, status.Error(codes.Internal, models.ErrInvalidCredentials.Error())
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// TODO: изменить формат сообщений protobuf. Не должно быть "Token"
	return &ssov1.LoginResponse{Token: authorizationCode}, nil
}
