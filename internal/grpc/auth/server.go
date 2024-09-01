package auth

import (
	"log/slog"
	"sso/internal/services/auth"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"github.com/go-playground/validator"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
	log       *slog.Logger
	validator validator.Validate
	auth      auth.Auth
}

func Register(log *slog.Logger, gRPC *grpc.Server, auth auth.Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{log: log, validator: *validator.New(), auth: auth})
}
