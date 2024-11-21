package router

import (
	handle404 "sso/internal/http/handlers/404"
	authService "sso/internal/services/auth"
	v1 "sso/internal/transport/http/router/v1"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authService authService.Auth,
) {
	v1Group := r.Group("/api/v1")
	v1.SetupV1Routes(v1Group, authService)

	r.NoRoute(handle404.Handle404)
}
