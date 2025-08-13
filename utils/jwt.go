package utils

import (
	"time"
	"blogapi/models"
	"github.com/golang-jwt/jwt/v5"
)


const jwtKey = "my-secret-key"

func GenerateJwt(username string, userId int64) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{

		UserId: userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "",err
	}

	return tokenString,nil
}
