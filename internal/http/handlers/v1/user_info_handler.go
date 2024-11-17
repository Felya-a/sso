package http_handlers_v1

import (
	"errors"
	"log/slog"
	"strings"

	. "sso/internal/http/handlers"
	"sso/internal/lib/logger"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model/errors"

	"github.com/gin-gonic/gin"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetUserInfoHandler(
	authService authService.Auth,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.Logger()
		log = log.With(
			slog.String("requestid", uuid.New().String()),
		)

		authorizationHeader := ctx.Request.Header.Get("Authorization")

		if authorizationHeader == "" {
			response := ErrorResponse{
				Status:  "error",
				Message: "failed",
				Error:   "authorization header is missing",
			}
			ctx.JSON(401, response)
			return
		}

		accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

		userInfo, err := authService.GetUserInfo(ctx, log, accessToken)
		if err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "failed",
				Error:   models.ErrInternal.Error(),
			}

			if errors.Is(err, models.ErrUserNotFound) ||
				errors.Is(err, models.ErrInvalidJwt) ||
				errors.Is(err, jwtLib.ErrTokenExpired) {
				response.Error = err.Error()
				ctx.JSON(400, response)
				return
			}

			ctx.JSON(500, response)
			return
		}

		response := SuccessResponse{
			Status:  "ok",
			Message: "success",
			Data:    GetUserInfoResponseDto(userInfo),
		}

		ctx.JSON(200, response)
	}
}
