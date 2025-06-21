package http_handlers_v1

import (
	"log/slog"
	. "sso/internal/http/handlers"
	"sso/internal/lib/logger"
	"sso/internal/lib/logger/sl"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func GetTokenHandler(
	authService authService.Auth,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto TokenRequestDto

		log := logger.Logger()
		log = log.With(
			slog.String("requestid", uuid.New().String()),
		)

		if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
			log.Info("parse body error", sl.Err(err))
			response := ErrorResponse{
				Status:  "error",
				Message: "parse body error",
				Error:   err.Error(),
			}
			ctx.JSON(400, response)
			return
		}

		if err := validator.New().Struct(dto); err != nil {
			log.Info("validation error", sl.Err(err))
			response := ErrorResponse{
				Status:  "error",
				Message: "validation error",
				Error:   err.Error(),
			}
			ctx.JSON(400, response)
			return
		}

		tokens, err := authService.Tokens(ctx, log, dto.AuthorizationCode)
		if err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "error on generate jwt tokens",
				Error:   "internal error",
			}
			if models.IsDefinedError(err) {
				response.Error = err.Error()
				ctx.JSON(400, response)
				return
			}
			ctx.JSON(500, response)
			return
		}

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
