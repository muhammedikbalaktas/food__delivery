package controllers

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const password = "SystemTokenPassword"

func parseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(password), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse claims")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found in claims")
	}

	return role, nil
}
