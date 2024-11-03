package auth_test_common

import (
	authModel "sso/internal/services/auth/model"

	"golang.org/x/crypto/bcrypt"
)

type UserWithPassword struct {
	authModel.UserModel
	Password string
}

var ValidUser UserWithPassword = UserWithPassword{
	UserModel: authModel.UserModel{
		ID:       1,
		Email:    "fake@local.com",
		PassHash: getPassHash("password123"),
	},
	Password: "password123",
}

func getPassHash(password string) []byte {
	passHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return passHash
}
