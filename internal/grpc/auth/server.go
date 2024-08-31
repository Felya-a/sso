package auth

import (
	"sso/internal/services/auth"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"github.com/go-playground/validator"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
	validator validator.Validate
	auth      auth.Auth
}

func Register(gRPC *grpc.Server, auth auth.Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{validator: *validator.New(), auth: auth})
}
