package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "sso/internal/grpc/auth"
	authService "sso/internal/services/auth"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(
	db *sqlx.DB,
	log *slog.Logger,
	port string,
) *App {
	gRPCServer := grpc.NewServer()
	authService := authService.New(db, log)
	authgrpc.Register(log, gRPCServer, authService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))
	if err = a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
