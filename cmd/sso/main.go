package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger/handlers/slogpretty"
	"strconv"
	"syscall"

	_ "database/sql" // ?

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const (
	envLocal = "local"
	envStage = "stage"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	log := setupLogger(config.Env)

	db := mustConnectPostgres(config)

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

func mustConnectPostgres(config config.Config) *sqlx.DB {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			config.Postgres.Host,
			strconv.Itoa(config.Postgres.Port),
			config.Postgres.User,
			config.Postgres.Database,
			config.Postgres.Password,
			"disable",
		),
	)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
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

func gracefulShutdown() os.Signal {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sgnl := <-stop
	return sgnl
}
