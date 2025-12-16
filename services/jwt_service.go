package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(tokenString string) (string, error)
}

type JWTServiceImpl struct {
	jwtSecret string
}

func NewJWTService() JWTService {
	return &JWTServiceImpl{
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}

func (s *JWTServiceImpl) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
		"iss": "money-api",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("erro ao fazer parse do token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("erro ao extrair claims do token")
	}

	if iss, ok := claims["iss"].(string); !ok || iss != "money-api" {
		return "", fmt.Errorf("issuer inválido")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("claim 'sub' não encontrada ou inválida")
	}

	return userID, nil
}
