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
	r.POST("/registration", handlers.GetRegistrationHandler(authService))
	r.POST("/login", handlers.GetLoginHandler(authService))
	r.POST("/token", handlers.GetTokenHandler(authService))
	r.GET("/refresh", handlers.GetRefreshHandler(authService))
	r.GET("/userinfo", handlers.GetUserInfoHandler(authService))

	/* FOR DEBUG ONLY */
	r.GET("/redirect", func(ctx *gin.Context) { ctx.Redirect(301, "https://google.com") })
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
