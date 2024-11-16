package auth

import (
	"log/slog"

	authService "sso/internal/services/auth"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth authService.Auth
}

func Register(log *slog.Logger, gRPC *grpc.Server, auth authService.Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}
