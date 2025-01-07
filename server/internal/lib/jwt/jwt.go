package jwt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
)

func New(claims jwt.MapClaims, secret string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("secret cannot be empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// Parse извлекает и проверяет данные из токена
func Parse[T any](tokenString string, secret string) (T, error) {
	var result T

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
			return result, jwt.ErrTokenExpired
		}
		return result, err
	}
	if !token.Valid {
		return result, fmt.Errorf("token not valid")
	}

	// Извлекаем claims
	claims := token.Claims.(jwt.MapClaims)

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return result, fmt.Errorf("failed to marshal claims: %w", err)
	}

	// Разбираем JSON строку в структуру типа T
	if err = json.Unmarshal(claimsJSON, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal claims into result: %w", err)
	}

	// Валидируем claims
	if err := validator.New().Struct(result); err != nil {
		return result, fmt.Errorf("failed to validation claims into result: %w", err)
	}

	return result, nil
}
