package app

import (
	"log/slog"
	grpcapp "sso/internal/transport/grpc"
	httpapp "sso/internal/transport/http"

	authService "sso/internal/services/auth"

	"github.com/jmoiron/sqlx"
)

type App struct {
	GrpcServer *grpcapp.App
	HttpServer *httpapp.HttpTransport
}

func New(
	db *sqlx.DB,
	log *slog.Logger,
	grpcPort string,
	httpPort string,
	authService authService.Auth,
) *App {
	grpcApp := grpcapp.New(db, log, grpcPort, authService)
	httpApp := httpapp.New(log, httpPort, authService)

	return &App{GrpcServer: grpcApp, HttpServer: httpApp}
}
