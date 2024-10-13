package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

func GenerateAccessToken(username string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(username string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
