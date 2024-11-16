package logger

import (
	"fmt"
	"log/slog"
	"os"

	"sso/internal/lib/logger/handlers/slogpretty"
)

var env string

func SetEnv(environment string) {
	env = environment
}

func Logger() *slog.Logger {
	switch env {
	case "local", "test":
		// return setupPrettySlog()
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "stage":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic(fmt.Sprintf("Unknown environment: %s", env))
	}
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
