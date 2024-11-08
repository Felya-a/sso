package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger/handlers/slogpretty"
	"sso/internal/utils"
	"syscall"

	_ "database/sql"

	_ "github.com/lib/pq"
)

const (
	envTest  = "test"
	envLocal = "local"
	envStage = "stage"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	log := setupLogger(config.Env)

	db := utils.MustConnectPostgres(config)

	application := app.New(db, log, config.Grpc.Port)

	go application.GrpcServer.MustRun()

	log.Info("Starting application", slog.Any("env", config.Env))

	// Graceful shutdown
	sgnl := gracefulShutdown()
	log.Info("Stopping application", slog.String("signal", sgnl.String()))

	application.GrpcServer.Stop()
	db.Close()

	log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal, envTest:
		log = setupPrettySlog()
		// log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envStage:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func gracefulShutdown() os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sgnl := <-stop
	return sgnl
}
