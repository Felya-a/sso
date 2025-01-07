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

type RegisterRequestValidate struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func (s *serverApi) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	log := logger.Logger()
	log = log.With(
		slog.String("requestid", uuid.New().String()),
	)

	dto := RegisterRequestValidate{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := validator.New().Struct(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, models.ErrInvalidCredentials.Error())
	}

	user, err := s.auth.Register(ctx, log, dto.Email, dto.Password)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return nil, status.Error(codes.Internal, models.ErrUserAlreadyExists.Error())
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.RegisterResponse{UserId: user.ID}, nil
}
