package auth

import (
	"context"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
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
	dto := LoginRequestValidate{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := s.validator.Struct(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.auth.RegisterNewUser(ctx, dto.Email, dto.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.RegisterResponse{UserId: int64(userId)}, nil
}
