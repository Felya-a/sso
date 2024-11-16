package http_handlers_v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	. "sso/internal/http/handlers"
	"sso/internal/lib/logger"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model/errors"

	"github.com/gin-gonic/gin"
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
				Message: "failed get user info",
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
				Message: "failed parse jwt token",
				Error:   models.ErrInternal.Error(),
			}

			if errors.Is(err, models.ErrUserNotFound) ||
				errors.Is(err, models.ErrInvalidJwt) {
				response.Error = err.Error()
				ctx.JSON(400, response)
				return
			}

			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		response := SuccessResponse{
			Status:  "ok",
			Message: "success parse jwt token",
			Data:    GetUserInfoResponseDto(userInfo),
		}

		ctx.JSON(http.StatusOK, response)
	}
}
