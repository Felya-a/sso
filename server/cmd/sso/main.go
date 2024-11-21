package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger"
	authService "sso/internal/services/auth"
	"sso/internal/utils"

	_ "database/sql"

	_ "github.com/lib/pq"
)

func main() {
	config := config.MustLoad()
	logger.SetEnv(config.Env)
	log := logger.Logger()

	db := utils.MustConnectPostgres(config)
	utils.Migrate(db)

	authService := authService.New(db)

	application := app.New(db, log, config.Grpc.Port, config.Http.Port, authService)

	go application.GrpcServer.MustRun()
	go application.HttpServer.MustRun()

	log.Info("Starting application", slog.Any("env", config.Env))

	// Graceful shutdown
	sgnl := gracefulShutdown()
	log.Info("Stopping application", slog.String("signal", sgnl.String()))

	application.GrpcServer.Stop()
	application.HttpServer.Stop()
	db.Close()

	log.Info("Application stopped")
}

func gracefulShutdown() os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sgnl := <-stop
	return sgnl
}