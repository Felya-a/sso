package handlers

import (
	"errors"
	"log/slog"
	"sso/internal/lib/logger"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func GetLoginHandler(
	authService authService.Auth,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto LoginRequestDto

		log := logger.Logger()
		log = log.With(
			slog.String("requestid", uuid.New().String()),
		)

		if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
			log.Info("parse body error", slog.String("error", err.Error()))
			response := ErrorResponse{
				Status:  "error",
				Message: "parse body error",
				Error:   err.Error(),
			}
			ctx.JSON(400, response)
			return
		}

		if err := validator.New().Struct(dto); err != nil {
			log.Info("validation error", slog.String("error", err.Error()))
			response := ErrorResponse{
				Status:  "error",
				Message: "validation error",
				Error:   err.Error(),
			}
			ctx.JSON(400, response)
			return
		}

		token, err := authService.Login(ctx, log, dto.Email, dto.Password, 1)
		if err != nil {
			log.Info("error on login", slog.String("error", err.Error()))
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
			return
		}

		log.Info("success login", slog.String("email", dto.Email))
		response := SuccessResponse{
			Status:  "ok",
			Message: "success login",
			Data:    GetLoginResponseDto(token),
		}
		ctx.JSON(200, response)
	}
}
