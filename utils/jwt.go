package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	Username string
	jwt.RegisteredClaims
}

const jwtKey = "my-secret-key"

func GenerateJwt(payloadString string) (string, error) {

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &claims{

		Username: payloadString,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "",err
	}

	return tokenString,nil
}
