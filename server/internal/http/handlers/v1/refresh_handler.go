package http_handlers_v1

import (
	"fmt"
	"log/slog"
	. "sso/internal/http/handlers"
	"sso/internal/lib/logger"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetRefreshHandler(
	authService authService.Auth,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.Logger()
		log = log.With(
			slog.String("requestid", uuid.New().String()),
		)

		refreshToken, _ := ctx.Cookie("refresh_token")
		if refreshToken == "" {
			response := ErrorResponse{
				Status:  "error",
				Message: "failed",
				Error:   "refresh token is empty",
			}
			ctx.JSON(400, response)
			return
		}

		tokens, err := authService.Refresh(ctx, log, refreshToken)
		if err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "error on registration",
				Error:   "internal error",
			}
			if models.IsDefinedError(err) {
				fmt.Println(err)
				response.Error = err.Error()
				ctx.JSON(400, response)
				return
			}
			ctx.JSON(500, response)
			return
		}

		// TODO: согласовать maxAge с временем жизни токенов
		ctx.SetCookie("access_token", tokens.AccessJwtToken, 30*24*60*60, "/", "", true, true)
		ctx.SetCookie("refresh_token", tokens.RefreshJwtToken, 30*24*60*60, "/", "", true, true)

		response := SuccessResponse{
			Status:  "ok",
			Message: "success",
			Data: TokenResponseDto{
				AccessToken:  tokens.AccessJwtToken,
				RefreshToken: tokens.RefreshJwtToken,
			},
		}
		ctx.JSON(200, response)
	}
}
