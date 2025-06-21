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

type UserInfoRequestValidate struct {
	AccessToken string `validate:"required"`
}

func (s *serverApi) UserInfo(
	ctx context.Context,
	req *ssov1.UserInfoRequest,
) (*ssov1.UserInfoResponse, error) {
	log := logger.Logger()
	log = log.With(
		slog.String("requestid", uuid.New().String()),
	)

	dto := UserInfoRequestValidate{
		AccessToken: req.GetAccessToken(),
	}

	if err := validator.New().Struct(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.auth.UserInfo(ctx, log, dto.AccessToken)
	if err != nil {
		if models.IsDefinedError(err) {
			if errors.Is(err, models.ErrJwtExpired) {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.UserInfoResponse{
		UserId: user.ID,
		Email:  user.Email,
	}, nil
}
