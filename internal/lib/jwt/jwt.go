package jwt

import (
	"errors"
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

// ParseToken извлекает и проверяет данные из токена
func ParseToken(tokenString string, secret string) (JwtBodyParams, error) {
	// Функция для проверки подписи
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что алгоритм подписи соответствует
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}

	// Парсим токен
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return JwtBodyParams{}, jwt.ErrTokenExpired
		}
		return JwtBodyParams{}, err
	}

	// Извлекаем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем и извлекаем данные из claims
		id, idOk := claims["id"].(float64) // jwt.MapClaims возвращает float64 для чисел
		email, emailOk := claims["email"].(string)

		if !idOk || !emailOk {
			return JwtBodyParams{}, errors.New("invalid token claims")
		}

		// Возвращаем распарсенные данные
		return JwtBodyParams{
			ID:    int64(id),
			Email: email,
		}, nil
	}

	return JwtBodyParams{}, errors.New("invalid token")
}
