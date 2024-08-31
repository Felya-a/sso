package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"

	"github.com/jmoiron/sqlx"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(
	db *sqlx.DB,
	log *slog.Logger,
	grpcPort string,
) *App {
	grpcApp := grpcapp.New(db, log, grpcPort)

	return &App{GrpcServer: grpcApp}
}
