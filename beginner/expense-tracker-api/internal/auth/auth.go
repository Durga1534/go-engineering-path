package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("your_secret_key")

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return SecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			return fmt.Sprintf("%d", int(userID)), nil
		}
	}
	return "", jwt.ErrInvalidKey
}
