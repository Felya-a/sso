package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort string,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{GRPCServer: grpcApp}
}
