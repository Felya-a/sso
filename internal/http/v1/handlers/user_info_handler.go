package handlers

import (
	"errors"
	"net/http"
	"strings"

	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model/errors"

	"github.com/gin-gonic/gin"
)

func GetUserInfoHandler(authService authService.Auth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.Request.Header.Get("Authorization")

		if authorizationHeader == "" {
			response := ErrorResponse{
				Status:  "error",
				Message: "failed get user info",
				Error:   "authorization header is missing",
			}
			ctx.JSON(401, response.FormatResponse())
			return
		}

		accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

		userInfo, err := authService.GetUserInfo(ctx, accessToken)
		if err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "failed parse jwt token",
				Error:   models.ErrInternal.Error(),
			}

			if errors.Is(err, models.ErrUserNotFound) ||
				errors.Is(err, models.ErrInvalidJwt) {
				response.Error = err.Error()
				ctx.JSON(400, response.FormatResponse())
				return
			}

			ctx.JSON(http.StatusInternalServerError, response.FormatResponse())
			return
		}

		response := SuccessResponse{
			Status:  "ok",
			Message: "success parse jwt token",
			Data:    NewUserInfoDto(userInfo),
		}

		ctx.JSON(http.StatusOK, response.FormatResponse())
	}
}
