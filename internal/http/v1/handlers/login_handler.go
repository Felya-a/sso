package handlers

import (
	"errors"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func GetLoginHandler(authService authService.Auth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto LoginRequestDto

		if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "parse body error",
				Error:   err.Error(),
			}
			ctx.JSON(400, response)
		}

		if err := validator.New().Struct(dto); err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "validation error",
				Error:   err.Error(),
			}
			ctx.JSON(400, response)
		}

		token, err := authService.Login(ctx, dto.Email, dto.Password, 1)
		if err != nil {
			response := ErrorResponse{
				Status:  "error",
				Message: "error on login",
				Error:   "internal error",
			}
			if errors.Is(err, models.ErrInternal) ||
				errors.Is(err, models.ErrInvalidCredentials) {
				response.Error = err.Error()
				ctx.JSON(400, response)
				return
			}
			ctx.JSON(500, response)
		}

		response := SuccessResponse{
			Status:  "ok",
			Message: "success login",
			Data:    GetLoginResponseDto(token),
		}
		ctx.JSON(200, response)
	}
}
