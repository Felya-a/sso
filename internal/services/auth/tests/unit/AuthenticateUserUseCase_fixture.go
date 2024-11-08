package auth_test

import (
	model "sso/internal/services/auth/model"

	"golang.org/x/crypto/bcrypt"
)

type UserWithPassword struct {
	model.UserModel
	password string
}

var ValidUser UserWithPassword = UserWithPassword{
	UserModel: model.UserModel{
		ID:       1,
		Email:    "fake@local.com",
		PassHash: getPassHash("password123"),
	},
	password: "password123",
}

func getPassHash(password string) []byte {
	passHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return passHash
}
