package adapters

import (
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	SecretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{SecretKey: secretKey}
}

func (j *JWTManager) GenerateToken(userId int, email string, name string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"name":    name,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (bool, map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return false, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, claims, nil
	}

	return false, nil, fmt.Errorf("invalid token")
}