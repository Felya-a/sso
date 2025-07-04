package http_handlers_v1

import (
	"log/slog"
	"strings"

	. "sso/internal/http/handlers"
	"sso/internal/lib/logger"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserInfoHandler(
	authService authService.Auth,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		log := logger.Logger()
		log = log.With(
			slog.String("requestid", uuid.New().String()),
		)

		accessTokenFromCookie, _ := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		accessTokenFromHeaders := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if accessTokenFromHeaders != "" {
			accessToken = accessTokenFromHeaders
		}
		if accessTokenFromCookie != "" {
			accessToken = accessTokenFromCookie
		}

		if accessToken == "" {
			response := ErrorResponse{
				Status:  "error",
				Message: "failed",
				Error:   "access token is missing",
			}
			ctx.JSON(401, response)
			return
		}

		userInfo, err := authService.UserInfo(ctx, log, accessToken)
		if err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "failed",
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

		response := SuccessResponse{
			Status:  "ok",
			Message: "success",
			Data:    UserInfoResponseDto{ID: userInfo.ID, Email: userInfo.Email},
		}

		ctx.JSON(200, response)
	}
}
