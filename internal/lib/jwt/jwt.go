package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtBodyParams struct {
	ID    int64
	Email string
}

func NewToken(body JwtBodyParams, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = body.ID
	claims["email"] = body.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
