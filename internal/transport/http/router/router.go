package router

import (
	authService "sso/internal/services/auth"
	v1 "sso/internal/transport/http/router/v1"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authService authService.Auth,
) {
	api := r.Group("/api/v1")
	v1.SetupUserRoutes(api, authService)
}
