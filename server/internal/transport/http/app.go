package http_app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"sso/internal/config"
	"sso/internal/lib/logger/sl"
	authService "sso/internal/services/auth"
	"sso/internal/transport/http/middleware"
	"sso/internal/transport/http/router"

	"github.com/gin-gonic/gin"
)

type HttpTransport struct {
	log        *slog.Logger
	httpServer *http.Server
	port       string
}

func New(
	log *slog.Logger,
	port string,
	authService authService.Auth,
) *HttpTransport {
	gin.SetMode(getGinMode())
	handler := gin.Default()
	handler.Use(middleware.CORSMiddleware())

	router.SetupRoutes(handler, authService)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	return &HttpTransport{
		log,
		httpServer,
		port,
	}
}

func (a *HttpTransport) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *HttpTransport) run() error {
	const op = "http_app.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	log.Info("http server is running", slog.String("addr", a.httpServer.Addr))
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("error on start http server", sl.Err(err))
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func getGinMode() string {
	var mode string

	switch config.Get().Env {
	case "local", "test":
		mode = "debug"
	case "stage", "prod":
		mode = "release"
	default:
		mode = "release"
	}

	return mode
}

func (a *HttpTransport) Stop() {
	a.log.Info("stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.log.Error("Server forced to shutdown: ", sl.Err(err))
	}
}
