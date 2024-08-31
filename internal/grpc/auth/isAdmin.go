package auth

import (
	"context"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IsAdminRequestValidate struct {
	UserId int `validate:"required"`
}

func (s *serverApi) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	dto := IsAdminRequestValidate{
		UserId: int(req.GetUserId()),
	}

	if err := s.validator.Struct(dto); err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, int64(dto.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
