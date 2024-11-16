package http_routes_v1

import (
	"sso/internal/http/v1/handlers"
	authService "sso/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(
	r *gin.RouterGroup,
	authService authService.Auth,
) {
	r.GET("/userinfo", handlers.GetUserInfoHandler(authService))
	r.POST("/login", handlers.GetLoginHandler(authService))

	/* FOR DEBUG ONLY */
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
