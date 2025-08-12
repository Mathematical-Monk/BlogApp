package utils

import (
	"time"
	"blogapi/models"
	"github.com/golang-jwt/jwt/v5"
)


const jwtKey = "my-secret-key"

func GenerateJwt(username string) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{

		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "",err
	}

	return tokenString,nil
}
