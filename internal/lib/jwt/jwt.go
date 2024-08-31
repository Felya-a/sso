package jwt

import (
	"sso/internal/config"
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NetToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(config.Get().JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
