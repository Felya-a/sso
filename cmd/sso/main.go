package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger/handlers/slogpretty"
	"syscall"
)

const (
	envLocal = "local"
	envStage = "stage"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	log := setupLogger(config.Env)

	application := app.New(log, config.GrpcConfig.Port, "", config.TokenTtl)

	go application.GRPCServer.MustRun()

	log.Info("Starting application", slog.Any("env", config.Env))

	// Graceful shutdown
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sgnl := <-stop

	log.Info("Stopping application", slog.String("signal", sgnl.String()))

	application.GRPCServer.Stop()

	log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
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
