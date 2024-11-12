package handlers

import (
	authModels "sso/internal/services/auth/model"
)

type UserInfoDto struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func NewUserInfoDto(user *authModels.UserModel) UserInfoDto {
	return UserInfoDto{
		ID:    user.ID,
		Email: user.Email,
	}
}
