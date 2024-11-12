package handlers

import (
	authModels "sso/internal/services/auth/model"
)

type UserInfoResponseDto struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func GetUserInfoResponseDto(user *authModels.UserModel) UserInfoResponseDto {
	return UserInfoResponseDto{
		ID:    user.ID,
		Email: user.Email,
	}
}

type LoginRequestDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDto struct {
	Token string `json:"token"`
}

func GetLoginResponseDto(token string) LoginResponseDto {
	return LoginResponseDto{Token: token}
}
