package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId   int64
	Username string
	jwt.RegisteredClaims
}
