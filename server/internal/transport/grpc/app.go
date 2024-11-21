package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	authgrpc "sso/internal/grpc/auth"
	authService "sso/internal/services/auth"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func LoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	log *slog.Logger,
) (resp interface{}, err error) {
	log.Info("Received request: %v", req)
	log.Info("Received request: %v", info.FullMethod)

	// Вызываем основной обработчик
	resp, err = handler(ctx, req)

	// Логируем ответ и ошибку, если она есть
	if err != nil {
		log.Error("Method: %s, Error: %v", info.FullMethod, status.Convert(err).Message())
	} else {
		log.Info("Response: %v", resp)
	}

	return resp, err
}

func LoggerInterceptorWithLog(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		return LoggerInterceptor(ctx, req, info, handler, log)
	}
}

func New(
	db *sqlx.DB,
	log *slog.Logger,
	port string,
	authService authService.Auth,
) *App {
	// gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(LoggerInterceptor))
	// gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(LoggerInterceptorWithLog(log)))
	gRPCServer := grpc.NewServer()
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
