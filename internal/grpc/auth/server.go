package auth

import (
	"context"
	"fmt"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"github.com/go-playground/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
	validator validator.Validate
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{validator: *validator.New()})
}

func (s *serverApi) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	validationReq := LoginRequestValidate{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppId:    int(req.GetAppId()),
	}

	if err := s.validator.Struct(validationReq); err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.LoginResponse{Token: "example"}, nil
}

func (s *serverApi) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverApi) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
