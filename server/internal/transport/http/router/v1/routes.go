package http_routes_v1

import (
	handlers "sso/internal/http/handlers/v1"
	authService "sso/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func SetupV1Routes(
	r *gin.RouterGroup,
	authService authService.Auth,
) {
	r.GET("/userinfo", handlers.GetUserInfoHandler(authService))
	r.POST("/login", handlers.GetLoginHandler(authService))
	r.POST("/registration", handlers.GetRegistrationHandler(authService))

	/* FOR DEBUG ONLY */
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
