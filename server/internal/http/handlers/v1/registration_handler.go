package http_handlers_v1

import (
	"errors"
	"log/slog"
	. "sso/internal/http/handlers"
	"sso/internal/lib/logger"
	"sso/internal/lib/logger/sl"
	authService "sso/internal/services/auth"
	models "sso/internal/services/auth/model/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func GetRegistrationHandler(
	authService authService.Auth,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto RegistrationRequestDto

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

		userid, err := authService.RegisterNewUser(ctx, log, dto.Email, dto.Password)
		if err != nil {
			log.Info("error on registration", sl.Err(err))
			response := ErrorResponse{
				Status:  "error",
				Message: "error on registration",
				Error:   "internal error",
			}
			if errors.Is(err, models.ErrInternal) ||
				errors.Is(err, models.ErrInvalidCredentials) ||
				errors.Is(err, models.ErrUserAlreadyExists) {
				response.Error = err.Error()
				ctx.JSON(400, response)
				return
			}
			ctx.JSON(500, response)
			return
		}

		log.Info("success register new user", slog.Int64("userid", userid))
		response := SuccessResponse{
			Status:  "ok",
			Message: "success registration",
			Data:    GetRegistrationResponseDto(userid),
		}
		ctx.JSON(200, response)
	}
}
