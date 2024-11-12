package auth

import (
	"log/slog"

	authService "sso/internal/services/auth"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"github.com/go-playground/validator"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
	log       *slog.Logger
	validator validator.Validate
	auth      authService.Auth
}

func Register(log *slog.Logger, gRPC *grpc.Server, auth authService.Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{log: log, validator: *validator.New(), auth: auth})
}
